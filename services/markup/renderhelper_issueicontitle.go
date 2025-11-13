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

package markup

import (
	"context"
	"errors"
	"fmt"
	"html/template"

	"github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/perm/access"
	"github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/htmlutil"
	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/util"
	kmup_context "github.com/kumose/kmup/services/context"
)

func renderRepoIssueIconTitle(ctx context.Context, opts markup.RenderIssueIconTitleOptions) (_ template.HTML, err error) {
	webCtx := kmup_context.GetWebContext(ctx)
	if webCtx == nil {
		return "", errors.New("context is not a web context")
	}

	textIssueIndex := fmt.Sprintf("(#%d)", opts.IssueIndex)
	dbRepo := webCtx.Repo.Repository
	if opts.OwnerName != "" {
		dbRepo, err = repo.GetRepositoryByOwnerAndName(ctx, opts.OwnerName, opts.RepoName)
		if err != nil {
			return "", err
		}
		textIssueIndex = fmt.Sprintf("(%s/%s#%d)", dbRepo.OwnerName, dbRepo.Name, opts.IssueIndex)
	}
	if dbRepo == nil {
		return "", nil
	}

	issue, err := issues.GetIssueByIndex(ctx, dbRepo.ID, opts.IssueIndex)
	if err != nil {
		return "", err
	}

	if webCtx.Repo.Repository == nil || dbRepo.ID != webCtx.Repo.Repository.ID {
		perms, err := access.GetUserRepoPermission(ctx, dbRepo, webCtx.Doer)
		if err != nil {
			return "", err
		}
		if !perms.CanReadIssuesOrPulls(issue.IsPull) {
			return "", util.ErrPermissionDenied
		}
	}

	if issue.IsPull {
		if err = issue.LoadPullRequest(ctx); err != nil {
			return "", err
		}
	}

	htmlIcon, err := webCtx.RenderToHTML("shared/issueicon", issue)
	if err != nil {
		return "", err
	}

	return htmlutil.HTMLFormat(`<a href="%s">%s %s %s</a>`, opts.LinkHref, htmlIcon, issue.Title, textIssueIndex), nil
}
