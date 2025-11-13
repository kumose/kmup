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

package user

import (
	"net/http"
	"slices"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	code_indexer "github.com/kumose/kmup/modules/indexer/code"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/routers/common"
	shared_user "github.com/kumose/kmup/routers/web/shared/user"
	"github.com/kumose/kmup/services/context"
)

const (
	tplUserCode templates.TplName = "user/code"
)

// CodeSearch render user/organization code search page
func CodeSearch(ctx *context.Context) {
	if !setting.Indexer.RepoIndexerEnabled {
		ctx.Redirect(ctx.ContextUser.HomeLink())
		return
	}
	if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	ctx.Data["IsPackageEnabled"] = setting.Packages.Enabled
	ctx.Data["Title"] = ctx.Tr("explore.code")
	ctx.Data["IsCodePage"] = true

	prepareSearch := common.PrepareCodeSearch(ctx)
	if prepareSearch.Keyword == "" {
		ctx.HTML(http.StatusOK, tplUserCode)
		return
	}

	var (
		repoIDs []int64
		err     error
	)

	page := ctx.FormInt("page")
	if page <= 0 {
		page = 1
	}

	repoIDs, err = repo_model.FindUserCodeAccessibleOwnerRepoIDs(ctx, ctx.ContextUser.ID, ctx.Doer)
	if err != nil {
		ctx.ServerError("FindUserCodeAccessibleOwnerRepoIDs", err)
		return
	}

	var (
		total                 int
		searchResults         []*code_indexer.Result
		searchResultLanguages []*code_indexer.SearchResultLanguages
	)

	if len(repoIDs) > 0 {
		total, searchResults, searchResultLanguages, err = code_indexer.PerformSearch(ctx, &code_indexer.SearchOptions{
			RepoIDs:    repoIDs,
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

		loadRepoIDs := make([]int64, 0, len(searchResults))
		for _, result := range searchResults {
			if !slices.Contains(loadRepoIDs, result.RepoID) {
				loadRepoIDs = append(loadRepoIDs, result.RepoID)
			}
		}

		repoMaps, err := repo_model.GetRepositoriesMapByIDs(ctx, loadRepoIDs)
		if err != nil {
			ctx.ServerError("GetRepositoriesMapByIDs", err)
			return
		}

		ctx.Data["RepoMaps"] = repoMaps
	}
	ctx.Data["SearchResults"] = searchResults
	ctx.Data["SearchResultLanguages"] = searchResultLanguages

	pager := context.NewPagination(total, setting.UI.RepoSearchPagingNum, page, 5)
	pager.AddParamFromRequest(ctx.Req)
	ctx.Data["Page"] = pager

	ctx.HTML(http.StatusOK, tplUserCode)
}
