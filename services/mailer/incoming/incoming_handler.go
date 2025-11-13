// Copyright (C) Kumo inc. and its affiliates.
// Author: Jeff.li lijippy@163.com
// All rights reserved.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package incoming

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	issues_model "github.com/kumose/kmup/models/issues"
	access_model "github.com/kumose/kmup/models/perm/access"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
	attachment_service "github.com/kumose/kmup/services/attachment"
	"github.com/kumose/kmup/services/context/upload"
	issue_service "github.com/kumose/kmup/services/issue"
	incoming_payload "github.com/kumose/kmup/services/mailer/incoming/payload"
	"github.com/kumose/kmup/services/mailer/token"
	pull_service "github.com/kumose/kmup/services/pull"
)

type MailHandler interface {
	Handle(ctx context.Context, content *MailContent, doer *user_model.User, payload []byte) error
}

var handlers = map[token.HandlerType]MailHandler{
	token.ReplyHandlerType:       &ReplyHandler{},
	token.UnsubscribeHandlerType: &UnsubscribeHandler{},
}

// ReplyHandler handles incoming emails to create a reply from them
type ReplyHandler struct{}

func (h *ReplyHandler) Handle(ctx context.Context, content *MailContent, doer *user_model.User, payload []byte) error {
	if doer == nil {
		return util.NewInvalidArgumentErrorf("doer can't be nil")
	}

	ref, err := incoming_payload.GetReferenceFromPayload(ctx, payload)
	if err != nil {
		return err
	}

	var issue *issues_model.Issue

	switch r := ref.(type) {
	case *issues_model.Issue:
		issue = r
	case *issues_model.Comment:
		comment := r

		if err := comment.LoadIssue(ctx); err != nil {
			return err
		}

		issue = comment.Issue
	default:
		return util.NewInvalidArgumentErrorf("unsupported reply reference: %v", ref)
	}

	if err := issue.LoadRepo(ctx); err != nil {
		return err
	}

	perm, err := access_model.GetUserRepoPermission(ctx, issue.Repo, doer)
	if err != nil {
		return err
	}

	// Locked issues require write permissions
	if issue.IsLocked && !perm.CanWriteIssuesOrPulls(issue.IsPull) && !doer.IsAdmin {
		log.Debug("can't write issue or pull")
		return nil
	}

	if !perm.CanReadIssuesOrPulls(issue.IsPull) {
		log.Debug("can't read issue or pull")
		return nil
	}

	attachmentIDs := make([]string, 0, len(content.Attachments))
	if setting.Attachment.Enabled {
		for _, attachment := range content.Attachments {
			attachmentBuf := bytes.NewReader(attachment.Content)
			uploaderFile := attachment_service.NewLimitedUploaderKnownSize(attachmentBuf, attachmentBuf.Size())
			a, err := attachment_service.UploadAttachmentGeneralSizeLimit(ctx, uploaderFile, setting.Attachment.AllowedTypes, &repo_model.Attachment{
				Name:       attachment.Name,
				UploaderID: doer.ID,
				RepoID:     issue.Repo.ID,
			})
			if err != nil {
				if upload.IsErrFileTypeForbidden(err) {
					log.Info("Skipping disallowed attachment type: %s", attachment.Name)
					continue
				}
				if errors.Is(err, util.ErrContentTooLarge) {
					log.Info("Skipping attachment exceeding size limit: %s", attachment.Name)
					continue
				}

				return err
			}
			attachmentIDs = append(attachmentIDs, a.UUID)
		}
	}

	if content.Content == "" && len(attachmentIDs) == 0 {
		return nil
	}

	switch r := ref.(type) {
	case *issues_model.Issue:
		_, err := issue_service.CreateIssueComment(ctx, doer, issue.Repo, issue, content.Content, attachmentIDs)
		if err != nil {
			return fmt.Errorf("CreateIssueComment failed: %w", err)
		}
	case *issues_model.Comment:
		comment := r

		switch comment.Type {
		case issues_model.CommentTypeCode:
			_, err := pull_service.CreateCodeComment(
				ctx,
				doer,
				nil,
				issue,
				comment.Line,
				content.Content,
				comment.TreePath,
				false, // not pending review but a single review
				comment.ReviewID,
				"",
				attachmentIDs,
			)
			if err != nil {
				return fmt.Errorf("CreateCodeComment failed: %w", err)
			}
		default:
			_, err := issue_service.CreateIssueComment(ctx, doer, issue.Repo, issue, content.Content, attachmentIDs)
			if err != nil {
				return fmt.Errorf("CreateIssueComment failed: %w", err)
			}
		}
	}
	return nil
}

// UnsubscribeHandler handles unwatching issues/pulls
type UnsubscribeHandler struct{}

func (h *UnsubscribeHandler) Handle(ctx context.Context, _ *MailContent, doer *user_model.User, payload []byte) error {
	if doer == nil {
		return util.NewInvalidArgumentErrorf("doer can't be nil")
	}

	ref, err := incoming_payload.GetReferenceFromPayload(ctx, payload)
	if err != nil {
		return err
	}

	switch r := ref.(type) {
	case *issues_model.Issue:
		issue := r

		if err := issue.LoadRepo(ctx); err != nil {
			return err
		}

		perm, err := access_model.GetUserRepoPermission(ctx, issue.Repo, doer)
		if err != nil {
			return err
		}

		if !perm.CanReadIssuesOrPulls(issue.IsPull) {
			log.Debug("can't read issue or pull")
			return nil
		}

		return issues_model.CreateOrUpdateIssueWatch(ctx, doer.ID, issue.ID, false)
	}

	return fmt.Errorf("unsupported unsubscribe reference: %v", ref)
}
