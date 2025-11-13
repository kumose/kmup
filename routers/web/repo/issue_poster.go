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
	"slices"
	"strings"

	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	shared_user "github.com/kumose/kmup/routers/web/shared/user"
	"github.com/kumose/kmup/services/context"
)

type userSearchInfo struct {
	UserID     int64  `json:"user_id"`
	UserName   string `json:"username"`
	AvatarLink string `json:"avatar_link"`
	FullName   string `json:"full_name"`
}

type userSearchResponse struct {
	Results []*userSearchInfo `json:"results"`
}

func IssuePullPosters(ctx *context.Context) {
	isPullList := ctx.PathParam("type") == "pulls"
	issuePosters(ctx, isPullList)
}

func issuePosters(ctx *context.Context, isPullList bool) {
	repo := ctx.Repo.Repository
	search := strings.TrimSpace(ctx.FormString("q"))
	posters, err := repo_model.GetIssuePostersWithSearch(ctx, repo, isPullList, search, setting.UI.DefaultShowFullName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if search == "" && ctx.Doer != nil {
		// the returned posters slice only contains limited number of users,
		// to make the current user (doer) can quickly filter their own issues, always add doer to the posters slice
		if !slices.ContainsFunc(posters, func(user *user_model.User) bool { return user.ID == ctx.Doer.ID }) {
			posters = append(posters, ctx.Doer)
		}
	}

	posters = shared_user.MakeSelfOnTop(ctx.Doer, posters)

	resp := &userSearchResponse{}
	resp.Results = make([]*userSearchInfo, len(posters))
	for i, user := range posters {
		resp.Results[i] = &userSearchInfo{UserID: user.ID, UserName: user.Name, AvatarLink: user.AvatarLink(ctx)}
		if setting.UI.DefaultShowFullName {
			resp.Results[i].FullName = user.FullName
		}
	}
	ctx.JSON(http.StatusOK, resp)
}
