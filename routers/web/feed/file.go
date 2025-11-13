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
	"strings"
	"time"

	"github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"

	"github.com/gorilla/feeds"
)

// ShowFileFeed shows tags and/or releases on the repo as RSS / Atom feed
func ShowFileFeed(ctx *context.Context, repo *repo.Repository, formatType string) {
	fileName := ctx.Repo.TreePath
	if len(fileName) == 0 {
		return
	}
	commits, err := ctx.Repo.GitRepo.CommitsByFileAndRange(
		git.CommitsByFileAndRangeOptions{
			Revision: ctx.Repo.RefFullName.ShortName(), // FIXME: legacy code used ShortName
			File:     fileName,
			Page:     1,
		})
	if err != nil {
		ctx.ServerError("ShowBranchFeed", err)
		return
	}

	title := "Latest commits for file " + ctx.Repo.TreePath

	link := &feeds.Link{Href: repo.HTMLURL() + "/" + ctx.Repo.RefTypeNameSubURL() + "/" + util.PathEscapeSegments(ctx.Repo.TreePath)}

	feed := &feeds.Feed{
		Title:       title,
		Link:        link,
		Description: repo.Description,
		Created:     time.Now(),
	}

	for _, commit := range commits {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:    commit.ID.String(),
			Title: strings.TrimSpace(strings.Split(commit.Message(), "\n")[0]),
			Link:  &feeds.Link{Href: repo.HTMLURL() + "/commit/" + commit.ID.String()},
			Author: &feeds.Author{
				Name:  commit.Author.Name,
				Email: commit.Author.Email,
			},
			Description: commit.Message(),
			Content:     commit.Message(),
			Created:     commit.Committer.When,
		})
	}

	writeFeed(ctx, feed, formatType)
}
