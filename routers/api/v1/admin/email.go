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

package admin

import (
	"net/http"

	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// GetAllEmails
func GetAllEmails(ctx *context.APIContext) {
	// swagger:operation GET /admin/emails admin adminGetAllEmails
	// ---
	// summary: List all emails
	// produces:
	// - application/json
	// parameters:
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/EmailList"
	//   "403":
	//     "$ref": "#/responses/forbidden"

	listOptions := utils.GetListOptions(ctx)

	emails, maxResults, err := user_model.SearchEmails(ctx, &user_model.SearchEmailOptions{
		Keyword:     ctx.PathParam("email"),
		ListOptions: listOptions,
	})
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	results := make([]*api.Email, len(emails))
	for i := range emails {
		results[i] = convert.ToEmailSearch(emails[i])
	}

	ctx.SetLinkHeader(int(maxResults), listOptions.PageSize)
	ctx.SetTotalCountHeader(maxResults)
	ctx.JSON(http.StatusOK, &results)
}

// SearchEmail
func SearchEmail(ctx *context.APIContext) {
	// swagger:operation GET /admin/emails/search admin adminSearchEmails
	// ---
	// summary: Search all emails
	// produces:
	// - application/json
	// parameters:
	// - name: q
	//   in: query
	//   description: keyword
	//   type: string
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/EmailList"
	//   "403":
	//     "$ref": "#/responses/forbidden"

	ctx.SetPathParam("email", ctx.FormTrim("q"))
	GetAllEmails(ctx)
}
