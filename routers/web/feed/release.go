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
	"time"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/services/context"

	"github.com/gorilla/feeds"
)

// shows tags and/or releases on the repo as RSS / Atom feed
func ShowReleaseFeed(ctx *context.Context, repo *repo_model.Repository, isReleasesOnly bool, formatType string) {
	releases, err := db.Find[repo_model.Release](ctx, repo_model.FindReleasesOptions{
		IncludeTags: !isReleasesOnly,
		RepoID:      ctx.Repo.Repository.ID,
	})
	if err != nil {
		ctx.ServerError("GetReleasesByRepoID", err)
		return
	}

	var title string
	var link *feeds.Link

	if isReleasesOnly {
		title = ctx.Locale.TrString("repo.release.releases_for", repo.FullName())
		link = &feeds.Link{Href: repo.HTMLURL() + "/release"}
	} else {
		title = ctx.Locale.TrString("repo.release.tags_for", repo.FullName())
		link = &feeds.Link{Href: repo.HTMLURL() + "/tags"}
	}

	feed := &feeds.Feed{
		Title:       title,
		Link:        link,
		Description: repo.Description,
		Created:     time.Now(),
	}

	feed.Items, err = releasesToFeedItems(ctx, releases)
	if err != nil {
		ctx.ServerError("releasesToFeedItems", err)
		return
	}

	writeFeed(ctx, feed, formatType)
}
