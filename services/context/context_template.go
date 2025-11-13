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

package context

import (
	"context"
	"net/http"
	"time"

	"github.com/kumose/kmup/services/webtheme"
)

var _ context.Context = TemplateContext(nil)

func NewTemplateContext(ctx context.Context, req *http.Request) TemplateContext {
	return TemplateContext{"_ctx": ctx, "_req": req}
}

func (c TemplateContext) parentContext() context.Context {
	return c["_ctx"].(context.Context)
}

func (c TemplateContext) Deadline() (deadline time.Time, ok bool) {
	return c.parentContext().Deadline()
}

func (c TemplateContext) Done() <-chan struct{} {
	return c.parentContext().Done()
}

func (c TemplateContext) Err() error {
	return c.parentContext().Err()
}

func (c TemplateContext) Value(key any) any {
	return c.parentContext().Value(key)
}

func (c TemplateContext) CurrentWebTheme() *webtheme.ThemeMetaInfo {
	req := c["_req"].(*http.Request)
	var themeName string
	if webCtx := GetWebContext(c); webCtx != nil {
		if webCtx.Doer != nil {
			themeName = webCtx.Doer.Theme
		}
	}
	if themeName == "" {
		if cookieTheme, _ := req.Cookie("kmup_theme"); cookieTheme != nil {
			themeName = cookieTheme.Value
		}
	}
	return webtheme.GuaranteeGetThemeMetaInfo(themeName)
}
