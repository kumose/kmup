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

package feed

import (
	"github.com/kumose/kmup/services/context"
)

// RenderBranchFeed render format for branch or file
func RenderBranchFeed(ctx *context.Context, feedType string) {
	if ctx.Repo.TreePath == "" {
		ShowBranchFeed(ctx, ctx.Repo.Repository, feedType)
	} else {
		ShowFileFeed(ctx, ctx.Repo.Repository, feedType)
	}
}

func RenderBranchFeedRSS(ctx *context.Context) {
	RenderBranchFeed(ctx, "rss")
}

func RenderBranchFeedAtom(ctx *context.Context) {
	RenderBranchFeed(ctx, "atom")
}
