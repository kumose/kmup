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

package org

import (
	"net/http"
	"path"
	"strings"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/organization"
	"github.com/kumose/kmup/models/renderhelper"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/markup/markdown"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/util"
	shared_user "github.com/kumose/kmup/routers/web/shared/user"
	"github.com/kumose/kmup/services/context"
)

const tplOrgHome templates.TplName = "org/home"

// Home show organization home page
func Home(ctx *context.Context) {
	uname := ctx.PathParam("username")

	if strings.HasSuffix(uname, ".keys") || strings.HasSuffix(uname, ".gpg") {
		ctx.NotFound(nil)
		return
	}

	ctx.SetPathParam("org", uname)
	context.OrgAssignment(context.OrgAssignmentOptions{})(ctx)
	if ctx.Written() {
		return
	}

	home(ctx, false)
}

func Repositories(ctx *context.Context) {
	home(ctx, true)
}

func home(ctx *context.Context, viewRepositories bool) {
	org := ctx.Org.Organization

	ctx.Data["PageIsUserProfile"] = true
	ctx.Data["Title"] = org.DisplayName()

	var orderBy db.SearchOrderBy
	sortOrder := ctx.FormString("sort")
	if _, ok := repo_model.OrderByFlatMap[sortOrder]; !ok {
		sortOrder = setting.UI.ExploreDefaultSort // TODO: add new default sort order for org home?
	}
	ctx.Data["SortType"] = sortOrder
	orderBy = repo_model.OrderByFlatMap[sortOrder]

	keyword := ctx.FormTrim("q")
	ctx.Data["Keyword"] = keyword

	language := ctx.FormTrim("language")
	ctx.Data["Language"] = language

	page := ctx.FormInt("page")
	if page <= 0 {
		page = 1
	}

	archived := ctx.FormOptionalBool("archived")
	ctx.Data["IsArchived"] = archived

	fork := ctx.FormOptionalBool("fork")
	ctx.Data["IsFork"] = fork

	mirror := ctx.FormOptionalBool("mirror")
	ctx.Data["IsMirror"] = mirror

	template := ctx.FormOptionalBool("template")
	ctx.Data["IsTemplate"] = template

	private := ctx.FormOptionalBool("private")
	ctx.Data["IsPrivate"] = private

	opts := &organization.FindOrgMembersOpts{
		Doer:         ctx.Doer,
		OrgID:        org.ID,
		IsDoerMember: ctx.Org.IsMember,
		ListOptions:  db.ListOptions{Page: 1, PageSize: 25},
	}

	members, _, err := organization.FindOrgMembers(ctx, opts)
	if err != nil {
		ctx.ServerError("FindOrgMembers", err)
		return
	}
	ctx.Data["Members"] = members
	ctx.Data["Teams"] = ctx.Org.Teams
	ctx.Data["DisableNewPullMirrors"] = setting.Mirror.DisableNewPull
	ctx.Data["ShowMemberAndTeamTab"] = ctx.Org.IsMember || len(members) > 0

	prepareResult, err := shared_user.RenderUserOrgHeader(ctx)
	if err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	// if no profile readme, it still means "view repositories"
	isViewOverview := !viewRepositories && prepareOrgProfileReadme(ctx, prepareResult)
	ctx.Data["PageIsViewRepositories"] = !isViewOverview
	ctx.Data["PageIsViewOverview"] = isViewOverview
	ctx.Data["ShowOrgProfileReadmeSelector"] = isViewOverview && prepareResult.ProfilePublicReadmeBlob != nil && prepareResult.ProfilePrivateReadmeBlob != nil

	repos, count, err := repo_model.SearchRepository(ctx, repo_model.SearchRepoOptions{
		ListOptions: db.ListOptions{
			PageSize: setting.UI.User.RepoPagingNum,
			Page:     page,
		},
		Keyword:            keyword,
		OwnerID:            org.ID,
		OrderBy:            orderBy,
		Private:            ctx.IsSigned,
		Actor:              ctx.Doer,
		Language:           language,
		IncludeDescription: setting.UI.SearchRepoDescription,
		Archived:           archived,
		Fork:               fork,
		Mirror:             mirror,
		Template:           template,
		IsPrivate:          private,
	})
	if err != nil {
		ctx.ServerError("SearchRepository", err)
		return
	}

	ctx.Data["Repos"] = repos
	ctx.Data["Total"] = count

	pager := context.NewPagination(int(count), setting.UI.User.RepoPagingNum, page, 5)
	pager.AddParamFromRequest(ctx.Req)
	ctx.Data["Page"] = pager

	ctx.HTML(http.StatusOK, tplOrgHome)
}

func prepareOrgProfileReadme(ctx *context.Context, prepareResult *shared_user.PrepareOwnerHeaderResult) bool {
	viewAs := ctx.FormString("view_as", util.Iif(ctx.Org.IsMember, "member", "public"))
	viewAsMember := viewAs == "member"

	var profileRepo *repo_model.Repository
	var readmeBlob *git.Blob
	if viewAsMember {
		if prepareResult.ProfilePrivateReadmeBlob != nil {
			profileRepo, readmeBlob = prepareResult.ProfilePrivateRepo, prepareResult.ProfilePrivateReadmeBlob
		} else {
			profileRepo, readmeBlob = prepareResult.ProfilePublicRepo, prepareResult.ProfilePublicReadmeBlob
			viewAsMember = false
		}
	} else {
		if prepareResult.ProfilePublicReadmeBlob != nil {
			profileRepo, readmeBlob = prepareResult.ProfilePublicRepo, prepareResult.ProfilePublicReadmeBlob
		} else {
			profileRepo, readmeBlob = prepareResult.ProfilePrivateRepo, prepareResult.ProfilePrivateReadmeBlob
			viewAsMember = true
		}
	}
	if readmeBlob == nil {
		return false
	}

	readmeBytes, err := readmeBlob.GetBlobContent(setting.UI.MaxDisplayFileSize)
	if err != nil {
		log.Error("failed to GetBlobContent for profile %q (view as %q) readme: %v", profileRepo.FullName(), viewAs, err)
		return false
	}

	rctx := renderhelper.NewRenderContextRepoFile(ctx, profileRepo, renderhelper.RepoFileOptions{
		CurrentRefPath: path.Join("branch", util.PathEscapeSegments(profileRepo.DefaultBranch)),
	})
	ctx.Data["ProfileReadmeContent"], err = markdown.RenderString(rctx, readmeBytes)
	if err != nil {
		log.Error("failed to GetBlobContent for profile %q (view as %q) readme: %v", profileRepo.FullName(), viewAs, err)
		return false
	}
	ctx.Data["IsViewingOrgAsMember"] = viewAsMember
	return true
}
