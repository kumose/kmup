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

package org

import (
	"github.com/kumose/kmup/routers/api/v1/shared"
	"github.com/kumose/kmup/services/context"
)

func ListBlocks(ctx *context.APIContext) {
	// swagger:operation GET /orgs/{org}/blocks organization organizationListBlocks
	// ---
	// summary: List users blocked by the organization
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/UserList"

	shared.ListBlocks(ctx, ctx.Org.Organization.AsUser())
}

func CheckUserBlock(ctx *context.APIContext) {
	// swagger:operation GET /orgs/{org}/blocks/{username} organization organizationCheckUserBlock
	// ---
	// summary: Check if a user is blocked by the organization
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: username
	//   in: path
	//   description: username of the user to check
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"

	shared.CheckUserBlock(ctx, ctx.Org.Organization.AsUser())
}

func BlockUser(ctx *context.APIContext) {
	// swagger:operation PUT /orgs/{org}/blocks/{username} organization organizationBlockUser
	// ---
	// summary: Block a user
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: username
	//   in: path
	//   description: username of the user to block
	//   type: string
	//   required: true
	// - name: note
	//   in: query
	//   description: optional note for the block
	//   type: string
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"
	//   "422":
	//     "$ref": "#/responses/validationError"

	shared.BlockUser(ctx, ctx.Org.Organization.AsUser())
}

func UnblockUser(ctx *context.APIContext) {
	// swagger:operation DELETE /orgs/{org}/blocks/{username} organization organizationUnblockUser
	// ---
	// summary: Unblock a user
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: username
	//   in: path
	//   description: username of the user to unblock
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"
	//   "422":
	//     "$ref": "#/responses/validationError"

	shared.UnblockUser(ctx, ctx.Doer, ctx.Org.Organization.AsUser())
}
