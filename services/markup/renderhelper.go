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

	"github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/markup"
	kmup_context "github.com/kumose/kmup/services/context"
)

func FormalRenderHelperFuncs() *markup.RenderHelperFuncs {
	return &markup.RenderHelperFuncs{
		RenderRepoFileCodePreview: renderRepoFileCodePreview,
		RenderRepoIssueIconTitle:  renderRepoIssueIconTitle,
		IsUsernameMentionable: func(ctx context.Context, username string) bool {
			mentionedUser, err := user.GetUserByName(ctx, username)
			if err != nil {
				return false
			}

			kmupCtx := kmup_context.GetWebContext(ctx)
			if kmupCtx == nil {
				// when using general context, use user's visibility to check
				return mentionedUser.Visibility.IsPublic()
			}

			// when using kmup context (web context), use user's visibility and user's permission to check
			return user.IsUserVisibleToViewer(kmupCtx, mentionedUser, kmupCtx.Doer)
		},
	}
}
