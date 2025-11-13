// Copyright 2014 The Gogs Authors. All rights reserved.
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

package web

import (
	"net/http"
	"strconv"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/sitemap"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/web/middleware"
	"github.com/kumose/kmup/routers/web/auth"
	"github.com/kumose/kmup/routers/web/user"
	"github.com/kumose/kmup/services/context"
)

const (
	// tplHome home page template
	tplHome templates.TplName = "home"
)

// Home render home page
func Home(ctx *context.Context) {
	if ctx.IsSigned {
		if !ctx.Doer.IsActive && setting.Service.RegisterEmailConfirm {
			ctx.Data["Title"] = ctx.Tr("auth.active_your_account")
			ctx.HTML(http.StatusOK, auth.TplActivate)
		} else if !ctx.Doer.IsActive || ctx.Doer.ProhibitLogin {
			log.Info("Failed authentication attempt for %s from %s", ctx.Doer.Name, ctx.RemoteAddr())
			ctx.Data["Title"] = ctx.Tr("auth.prohibit_login")
			ctx.HTML(http.StatusOK, "user/auth/prohibit_login")
		} else if ctx.Doer.MustChangePassword {
			ctx.Data["Title"] = ctx.Tr("auth.must_change_password")
			ctx.Data["ChangePasscodeLink"] = setting.AppSubURL + "/user/change_password"
			middleware.SetRedirectToCookie(ctx.Resp, setting.AppSubURL+ctx.Req.URL.RequestURI())
			ctx.Redirect(setting.AppSubURL + "/user/settings/change_password")
		} else {
			user.Dashboard(ctx)
		}
		return
		// Check non-logged users landing page.
	} else if setting.LandingPageURL != setting.LandingPageHome {
		ctx.Redirect(setting.AppSubURL + string(setting.LandingPageURL))
		return
	}

	// Check auto-login.
	if ctx.GetSiteCookie(setting.CookieRememberName) != "" {
		ctx.Redirect(setting.AppSubURL + "/user/login")
		return
	}

	ctx.Data["PageIsHome"] = true
	ctx.Data["IsRepoIndexerEnabled"] = setting.Indexer.RepoIndexerEnabled
	ctx.HTML(http.StatusOK, tplHome)
}

// HomeSitemap renders the main sitemap
func HomeSitemap(ctx *context.Context) {
	m := sitemap.NewSitemapIndex()
	if !setting.Service.Explore.DisableUsersPage {
		_, cnt, err := user_model.SearchUsers(ctx, user_model.SearchUserOptions{
			Types:       []user_model.UserType{user_model.UserTypeIndividual},
			ListOptions: db.ListOptions{PageSize: 1},
			IsActive:    optional.Some(true),
			Visible:     []structs.VisibleType{structs.VisibleTypePublic},
		})
		if err != nil {
			ctx.ServerError("SearchUsers", err)
			return
		}
		count := int(cnt)
		idx := 1
		for i := 0; i < count; i += setting.UI.SitemapPagingNum {
			m.Add(sitemap.URL{URL: setting.AppURL + "explore/users/sitemap-" + strconv.Itoa(idx) + ".xml"})
			idx++
		}
	}

	_, cnt, err := repo_model.SearchRepository(ctx, repo_model.SearchRepoOptions{
		ListOptions: db.ListOptions{
			PageSize: 1,
		},
		Actor:     ctx.Doer,
		AllPublic: true,
	})
	if err != nil {
		ctx.ServerError("SearchRepository", err)
		return
	}
	count := int(cnt)
	idx := 1
	for i := 0; i < count; i += setting.UI.SitemapPagingNum {
		m.Add(sitemap.URL{URL: setting.AppURL + "explore/repos/sitemap-" + strconv.Itoa(idx) + ".xml"})
		idx++
	}

	ctx.Resp.Header().Set("Content-Type", "text/xml")
	if _, err := m.WriteTo(ctx.Resp); err != nil {
		log.Error("Failed writing sitemap: %v", err)
	}
}

// NotFound render 404 page
func NotFound(ctx *context.Context) {
	ctx.Data["Title"] = "Page Not Found"
	ctx.NotFound(nil)
}
