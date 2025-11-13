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

	repo_module "github.com/kumose/kmup/modules/repository"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// Shows a list of all Label templates
func ListLabelTemplates(ctx *context.APIContext) {
	// swagger:operation GET /label/templates miscellaneous listLabelTemplates
	// ---
	// summary: Returns a list of all label templates
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/LabelTemplateList"
	result := make([]string, len(repo_module.LabelTemplateFiles))
	for i := range repo_module.LabelTemplateFiles {
		result[i] = repo_module.LabelTemplateFiles[i].DisplayName
	}

	ctx.JSON(http.StatusOK, result)
}

// Shows all labels in a template
func GetLabelTemplate(ctx *context.APIContext) {
	// swagger:operation GET /label/templates/{name} miscellaneous getLabelTemplateInfo
	// ---
	// summary: Returns all labels in a template
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
	//     "$ref": "#/responses/LabelTemplateInfo"
	//   "404":
	//     "$ref": "#/responses/notFound"
	name := util.PathJoinRelX(ctx.PathParam("name"))

	labels, err := repo_module.LoadTemplateLabelsByDisplayName(name)
	if err != nil {
		ctx.APIErrorNotFound()
		return
	}

	ctx.JSON(http.StatusOK, convert.ToLabelTemplateList(labels))
}
