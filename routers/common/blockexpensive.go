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

package common

import (
	"net/http"
	"strings"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/reqctx"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/web/middleware"

	"github.com/go-chi/chi/v5"
)

func BlockExpensive() func(next http.Handler) http.Handler {
	if !setting.Service.BlockAnonymousAccessExpensive {
		return nil
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ret := determineRequestPriority(reqctx.FromContext(req.Context()))
			if !ret.SignedIn {
				if ret.Expensive || ret.LongPolling {
					http.Redirect(w, req, setting.AppSubURL+"/user/login", http.StatusSeeOther)
					return
				}
			}
			next.ServeHTTP(w, req)
		})
	}
}

func isRoutePathExpensive(routePattern string) bool {
	if strings.HasPrefix(routePattern, "/user/") || strings.HasPrefix(routePattern, "/login/") {
		return false
	}

	expensivePaths := []string{
		// code related
		"/{username}/{reponame}/archive/",
		"/{username}/{reponame}/blame/",
		"/{username}/{reponame}/commit/",
		"/{username}/{reponame}/commits/",
		"/{username}/{reponame}/compare/",
		"/{username}/{reponame}/graph",
		"/{username}/{reponame}/media/",
		"/{username}/{reponame}/raw/",
		"/{username}/{reponame}/rss/branch/",
		"/{username}/{reponame}/src/",

		// issue & PR related (no trailing slash)
		"/{username}/{reponame}/issues",
		"/{username}/{reponame}/{type:issues}",
		"/{username}/{reponame}/pulls",
		"/{username}/{reponame}/{type:pulls}",

		// wiki
		"/{username}/{reponame}/wiki/",

		// activity
		"/{username}/{reponame}/activity/",
	}
	for _, path := range expensivePaths {
		if strings.HasPrefix(routePattern, path) {
			return true
		}
	}
	return false
}

func isRoutePathForLongPolling(routePattern string) bool {
	return routePattern == "/user/events"
}

func determineRequestPriority(reqCtx reqctx.RequestContext) (ret struct {
	SignedIn    bool
	Expensive   bool
	LongPolling bool
},
) {
	chiRoutePath := chi.RouteContext(reqCtx).RoutePattern()
	if _, ok := reqCtx.GetData()[middleware.ContextDataKeySignedUser].(*user_model.User); ok {
		ret.SignedIn = true
	} else {
		ret.Expensive = isRoutePathExpensive(chiRoutePath)
		ret.LongPolling = isRoutePathForLongPolling(chiRoutePath)
	}
	return ret
}
