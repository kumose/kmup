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

package repo

import (
	"net/http"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/context"
)

// GetLicenses returns licenses
func GetLicenses(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/licenses repository repoGetLicenses
	// ---
	// summary: Get repo licenses
	// produces:
	//   - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// responses:
	//   "404":
	//     "$ref": "#/responses/notFound"
	//   "200":
	//     "$ref": "#/responses/LicensesList"

	licenses, err := repo_model.GetRepoLicenses(ctx, ctx.Repo.Repository)
	if err != nil {
		log.Error("GetRepoLicenses failed: %v", err)
		ctx.APIErrorInternal(err)
		return
	}

	resp := make([]string, len(licenses))
	for i := range licenses {
		resp[i] = licenses[i].License
	}

	ctx.JSON(http.StatusOK, resp)
}
