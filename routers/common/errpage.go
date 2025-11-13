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
	"fmt"
	"net/http"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/httpcache"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/web/middleware"
	"github.com/kumose/kmup/modules/web/routing"
	"github.com/kumose/kmup/services/context"
)

const tplStatus500 templates.TplName = "status/500"

// RenderPanicErrorPage renders a 500 page, and it never panics
func RenderPanicErrorPage(w http.ResponseWriter, req *http.Request, err any) {
	combinedErr := fmt.Sprintf("%v\n%s", err, log.Stack(2))
	log.Error("PANIC: %s", combinedErr)

	defer func() {
		if err := recover(); err != nil {
			log.Error("Panic occurs again when rendering error page: %v. Stack:\n%s", err, log.Stack(2))
		}
	}()

	routing.UpdatePanicError(req.Context(), err)

	httpcache.SetCacheControlInHeader(w.Header(), &httpcache.CacheControlOptions{NoTransform: true})
	w.Header().Set(`X-Frame-Options`, setting.CORSConfig.XFrameOptions)

	tmplCtx := context.NewTemplateContext(req.Context(), req)
	tmplCtx["Locale"] = middleware.Locale(w, req)
	ctxData := middleware.GetContextData(req.Context())

	// This recovery handler could be called without Kmup's web context, so we shouldn't touch that context too much.
	// Otherwise, the 500-page may cause new panics, eg: cache.GetContextWithData, it makes the developer&users couldn't find the original panic.
	user, _ := ctxData[middleware.ContextDataKeySignedUser].(*user_model.User)
	if !setting.IsProd || (user != nil && user.IsAdmin) {
		ctxData["ErrorMsg"] = "PANIC: " + combinedErr
	}

	err = templates.HTMLRenderer().HTML(w, http.StatusInternalServerError, tplStatus500, ctxData, tmplCtx)
	if err != nil {
		log.Error("Error occurs again when rendering error page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error, please collect error logs and report to Kmup issue tracker"))
	}
}
