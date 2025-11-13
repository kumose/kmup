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
	"net/http"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/git"
	code_indexer "github.com/kumose/kmup/modules/indexer/code"
	"github.com/kumose/kmup/modules/indexer/code/gitgrep"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/routers/common"
	"github.com/kumose/kmup/services/context"
)

const tplSearch templates.TplName = "repo/search"

// Search render repository search page
func Search(ctx *context.Context) {
	ctx.Data["PageIsViewCode"] = true
	prepareSearch := common.PrepareCodeSearch(ctx)
	if prepareSearch.Keyword == "" {
		ctx.HTML(http.StatusOK, tplSearch)
		return
	}

	page := ctx.FormInt("page")
	if page <= 0 {
		page = 1
	}

	var total int
	var searchResults []*code_indexer.Result
	var searchResultLanguages []*code_indexer.SearchResultLanguages
	if setting.Indexer.RepoIndexerEnabled {
		var err error
		total, searchResults, searchResultLanguages, err = code_indexer.PerformSearch(ctx, &code_indexer.SearchOptions{
			RepoIDs:    []int64{ctx.Repo.Repository.ID},
			Keyword:    prepareSearch.Keyword,
			SearchMode: prepareSearch.SearchMode,
			Language:   prepareSearch.Language,
			Paginator: &db.ListOptions{
				Page:     page,
				PageSize: setting.UI.RepoSearchPagingNum,
			},
		})
		if err != nil {
			if code_indexer.IsAvailable(ctx) {
				ctx.ServerError("SearchResults", err)
				return
			}
			ctx.Data["CodeIndexerUnavailable"] = true
		} else {
			ctx.Data["CodeIndexerUnavailable"] = !code_indexer.IsAvailable(ctx)
		}
	} else {
		var err error
		// ref should be default branch or the first existing branch
		searchRef := git.RefNameFromBranch(ctx.Repo.Repository.DefaultBranch)
		searchResults, total, err = gitgrep.PerformSearch(ctx, page, ctx.Repo.Repository.ID, ctx.Repo.GitRepo, searchRef, prepareSearch.Keyword, prepareSearch.SearchMode)
		if err != nil {
			ctx.ServerError("gitgrep.PerformSearch", err)
			return
		}
	}

	ctx.Data["Repo"] = ctx.Repo.Repository
	ctx.Data["SearchResults"] = searchResults
	ctx.Data["SearchResultLanguages"] = searchResultLanguages

	pager := context.NewPagination(total, setting.UI.RepoSearchPagingNum, page, 5)
	pager.AddParamFromRequest(ctx.Req)
	ctx.Data["Page"] = pager

	ctx.HTML(http.StatusOK, tplSearch)
}
