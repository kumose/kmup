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

package install

import (
	"fmt"
	"html"
	"net/http"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/public"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/routers/common"
	"github.com/kumose/kmup/routers/web/healthcheck"
	"github.com/kumose/kmup/routers/web/misc"
	"github.com/kumose/kmup/services/forms"
)

// Routes registers the installation routes
func Routes() *web.Router {
	base := web.NewRouter()
	base.Use(common.ProtocolMiddlewares()...)
	base.Methods("GET, HEAD", "/assets/*", public.FileHandlerFunc())

	r := web.NewRouter()
	if sessionMid, err := common.Sessioner(); err == nil && sessionMid != nil {
		r.Use(sessionMid, Contexter())
	} else {
		log.Fatal("common.Sessioner failed: %v", err)
	}
	r.Get("/", Install) // it must be on the root, because the "install.js" use the window.location to replace the "localhost" AppURL
	r.Post("/", web.Bind(forms.InstallForm{}), SubmitInstall)
	r.Get("/post-install", InstallDone)

	r.Get("/-/web-theme/list", misc.WebThemeList)
	r.Post("/-/web-theme/apply", misc.WebThemeApply)
	r.Get("/api/healthz", healthcheck.Check)

	r.NotFound(installNotFound)

	base.Mount("", r)
	return base
}

func installNotFound(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Refresh", "1; url="+setting.AppSubURL+"/")
	// do not use 30x status, because the "post-install" page needs to use 404/200 to detect if Kmup has been installed.
	// the fetch API could follow 30x requests to the page with 200 status.
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprintf(w, `Not Found. <a href="%s">Go to default page</a>.`, html.EscapeString(setting.AppSubURL+"/"))
}
