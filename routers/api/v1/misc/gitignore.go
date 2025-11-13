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

package misc

import (
	"net/http"

	"github.com/kumose/kmup/modules/options"
	repo_module "github.com/kumose/kmup/modules/repository"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
)

// Shows a list of all Gitignore templates
func ListGitignoresTemplates(ctx *context.APIContext) {
	// swagger:operation GET /gitignore/templates miscellaneous listGitignoresTemplates
	// ---
	// summary: Returns a list of all gitignore templates
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/GitignoreTemplateList"
	ctx.JSON(http.StatusOK, repo_module.Gitignores)
}

// SHows information about a gitignore template
func GetGitignoreTemplateInfo(ctx *context.APIContext) {
	// swagger:operation GET /gitignore/templates/{name} miscellaneous getGitignoreTemplateInfo
	// ---
	// summary: Returns information about a gitignore template
	// produces:
	// - application/json
	// parameters:
	// - name: name
	//   in: path
	//   description: name of the template
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/GitignoreTemplateInfo"
	//   "404":
	//     "$ref": "#/responses/notFound"
	name := util.PathJoinRelX(ctx.PathParam("name"))

	text, err := options.Gitignore(name)
	if err != nil {
		ctx.APIErrorNotFound()
		return
	}

	ctx.JSON(http.StatusOK, &structs.GitignoreTemplateInfo{Name: name, Source: string(text)})
}
