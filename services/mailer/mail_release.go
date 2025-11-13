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

package mailer

import (
	"bytes"
	"context"
	"fmt"

	"github.com/kumose/kmup/models/renderhelper"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/markup/markdown"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/translation"
	sender_service "github.com/kumose/kmup/services/mailer/sender"
)

const tplNewReleaseMail templates.TplName = "repo/release"

func generateMessageIDForRelease(release *repo_model.Release) string {
	return fmt.Sprintf("<%s/releases/%d@%s>", release.Repo.FullName(), release.ID, setting.Domain)
}

// MailNewRelease send new release notify to all repo watchers.
func MailNewRelease(ctx context.Context, rel *repo_model.Release) {
	if setting.MailService == nil {
		// No mail service configured
		return
	}

	watcherIDList, err := repo_model.GetRepoWatchersIDs(ctx, rel.RepoID)
	if err != nil {
		log.Error("GetRepoWatchersIDs(%d): %v", rel.RepoID, err)
		return
	}

	recipients, err := user_model.GetMailableUsersByIDs(ctx, watcherIDList, false)
	if err != nil {
		log.Error("user_model.GetMailableUsersByIDs: %v", err)
		return
	}

	langMap := make(map[string][]*user_model.User)
	for _, user := range recipients {
		if user.ID != rel.PublisherID {
			langMap[user.Language] = append(langMap[user.Language], user)
		}
	}

	for lang, tos := range langMap {
		mailNewRelease(ctx, lang, tos, rel)
	}
}

func mailNewRelease(ctx context.Context, lang string, tos []*user_model.User, rel *repo_model.Release) {
	locale := translation.NewLocale(lang)

	var err error
	rctx := renderhelper.NewRenderContextRepoComment(ctx, rel.Repo).WithUseAbsoluteLink(true)
	rel.RenderedNote, err = markdown.RenderString(rctx,
		rel.Note)
	if err != nil {
		log.Error("markdown.RenderString(%d): %v", rel.RepoID, err)
		return
	}

	subject := locale.TrString("mail.release.new.subject", rel.TagName, rel.Repo.FullName())
	mailMeta := map[string]any{
		"locale":   locale,
		"Release":  rel,
		"Subject":  subject,
		"Language": locale.Language(),
		"Link":     rel.HTMLURL(),
	}

	var mailBody bytes.Buffer

	if err := LoadedTemplates().BodyTemplates.ExecuteTemplate(&mailBody, string(tplNewReleaseMail), mailMeta); err != nil {
		log.Error("ExecuteTemplate [%s]: %v", string(tplNewReleaseMail)+"/body", err)
		return
	}

	msgs := make([]*sender_service.Message, 0, len(tos))
	publisherName := fromDisplayName(rel.Publisher)
	msgID := generateMessageIDForRelease(rel)
	for _, to := range tos {
		msg := sender_service.NewMessageFrom(to.EmailTo(), publisherName, setting.MailService.FromEmail, subject, mailBody.String())
		msg.Info = subject
		msg.SetHeader("Message-ID", msgID)
		msgs = append(msgs, msg)
	}

	SendAsync(msgs...)
}
