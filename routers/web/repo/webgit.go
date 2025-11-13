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

package repo

import (
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
	files_service "github.com/kumose/kmup/services/repository/files"
)

func WebGitOperationCommonData(ctx *context.Context) {
	// TODO: more places like "wiki page" and "merging a pull request or creating an auto merge merging task"
	emails, err := user_model.GetActivatedEmailAddresses(ctx, ctx.Doer.ID)
	if err != nil {
		log.Error("WebGitOperationCommonData: GetActivatedEmailAddresses: %v", err)
	}
	if ctx.Doer.KeepEmailPrivate {
		emails = append([]string{ctx.Doer.GetPlaceholderEmail()}, emails...)
	}
	ctx.Data["CommitCandidateEmails"] = emails
	ctx.Data["CommitDefaultEmail"] = ctx.Doer.GetEmail()
}

func WebGitOperationGetCommitChosenEmailIdentity(ctx *context.Context, email string) (_ *files_service.IdentityOptions, valid bool) {
	if ctx.Data["CommitCandidateEmails"] == nil {
		setting.PanicInDevOrTesting("no CommitCandidateEmails in context data")
	}
	emails, _ := ctx.Data["CommitCandidateEmails"].([]string)
	if email == "" {
		return nil, true
	}
	if util.SliceContainsString(emails, email, true) {
		return &files_service.IdentityOptions{GitUserEmail: email}, true
	}
	return nil, false
}
