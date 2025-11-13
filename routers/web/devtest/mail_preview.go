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

package devtest

import (
	"net/http"
	"strings"

	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/mailer"

	"gopkg.in/yaml.v3"
)

func MailPreviewRender(ctx *context.Context) {
	tmplName := ctx.PathParam("*")
	mockDataContent, err := templates.AssetFS().ReadFile("mail/" + tmplName + ".devtest.yml")
	mockData := map[string]any{}
	if err == nil {
		err = yaml.Unmarshal(mockDataContent, &mockData)
		if err != nil {
			http.Error(ctx.Resp, "Failed to parse mock data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	mockData["locale"] = ctx.Locale
	err = mailer.LoadedTemplates().BodyTemplates.ExecuteTemplate(ctx.Resp, tmplName, mockData)
	if err != nil {
		_, _ = ctx.Resp.Write([]byte(err.Error()))
	}
}

func prepareMailPreviewRender(ctx *context.Context, tmplName string) {
	tmplSubject := mailer.LoadedTemplates().SubjectTemplates.Lookup(tmplName)
	if tmplSubject == nil {
		ctx.Data["RenderMailSubject"] = "default subject"
	} else {
		var buf strings.Builder
		err := tmplSubject.Execute(&buf, nil)
		if err != nil {
			ctx.Data["RenderMailSubject"] = err.Error()
		} else {
			ctx.Data["RenderMailSubject"] = buf.String()
		}
	}
	ctx.Data["RenderMailTemplateName"] = tmplName
}

func MailPreview(ctx *context.Context) {
	ctx.Data["MailTemplateNames"] = mailer.LoadedTemplates().TemplateNames
	tmplName := ctx.FormString("tmpl")
	if tmplName != "" {
		prepareMailPreviewRender(ctx, tmplName)
	}
	ctx.HTML(http.StatusOK, "devtest/mail-preview")
}
