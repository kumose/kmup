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
	"fmt"
	"net/http"
	"net/url"

	"github.com/kumose/kmup/modules/options"
	repo_module "github.com/kumose/kmup/modules/repository"
	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
)

// Returns a list of all License templates
func ListLicenseTemplates(ctx *context.APIContext) {
	// swagger:operation GET /licenses miscellaneous listLicenseTemplates
	// ---
	// summary: Returns a list of all license templates
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/LicenseTemplateList"
	response := make([]api.LicensesTemplateListEntry, len(repo_module.Licenses))
	for i, license := range repo_module.Licenses {
		response[i] = api.LicensesTemplateListEntry{
			Key:  license,
			Name: license,
			URL:  fmt.Sprintf("%sapi/v1/licenses/%s", setting.AppURL, url.PathEscape(license)),
		}
	}
	ctx.JSON(http.StatusOK, response)
}

func GetLicenseTemplateInfo(ctx *context.APIContext) {
	// swagger:operation GET /licenses/{name} miscellaneous getLicenseTemplateInfo
	// ---
	// summary: Returns information about a license template
	// produces:
	// - application/json
	// parameters:
	// - name: name
	//   in: path
	//   description: name of the license
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/LicenseTemplateInfo"
	//   "404":
	//     "$ref": "#/responses/notFound"
	name := util.PathJoinRelX(ctx.PathParam("name"))

	text, err := options.License(name)
	if err != nil {
		ctx.APIErrorNotFound()
		return
	}

	response := api.LicenseTemplateInfo{
		Key:  name,
		Name: name,
		URL:  fmt.Sprintf("%sapi/v1/licenses/%s", setting.AppURL, url.PathEscape(name)),
		Body: string(text),
		// This is for combatibilty with the GitHub API. This Text is for some reason added to each License response.
		Implementation: "Create a text file (typically named LICENSE or LICENSE.txt) in the root of your source code and copy the text of the license into the file",
	}

	ctx.JSON(http.StatusOK, response)
}
