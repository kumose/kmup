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
	"net/http"

	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
	webhook_service "github.com/kumose/kmup/services/webhook"
)

// ListHooks list an organziation's webhooks
func ListHooks(ctx *context.APIContext) {
	// swagger:operation GET /orgs/{org}/hooks organization orgListHooks
	// ---
	// summary: List an organization's webhooks
	// produces:
	// - application/json
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
	// responses:
	//   "200":
	//     "$ref": "#/responses/HookList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	utils.ListOwnerHooks(
		ctx,
		ctx.ContextUser,
	)
}

// GetHook get an organization's hook by id
func GetHook(ctx *context.APIContext) {
	// swagger:operation GET /orgs/{org}/hooks/{id} organization orgGetHook
	// ---
	// summary: Get a hook
	// produces:
	// - application/json
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: id
	//   in: path
	//   description: id of the hook to get
	//   type: integer
	//   format: int64
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/Hook"
	//   "404":
	//     "$ref": "#/responses/notFound"

	hook, err := utils.GetOwnerHook(ctx, ctx.ContextUser.ID, ctx.PathParamInt64("id"))
	if err != nil {
		return
	}

	apiHook, err := webhook_service.ToHook(ctx.ContextUser.HomeLink(), hook)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	ctx.JSON(http.StatusOK, apiHook)
}

// CreateHook create a hook for an organization
func CreateHook(ctx *context.APIContext) {
	// swagger:operation POST /orgs/{org}/hooks organization orgCreateHook
	// ---
	// summary: Create a hook
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreateHookOption"
	// responses:
	//   "201":
	//     "$ref": "#/responses/Hook"
	//   "404":
	//     "$ref": "#/responses/notFound"

	utils.AddOwnerHook(
		ctx,
		ctx.ContextUser,
		web.GetForm(ctx).(*api.CreateHookOption),
	)
}

// EditHook modify a hook of an organization
func EditHook(ctx *context.APIContext) {
	// swagger:operation PATCH /orgs/{org}/hooks/{id} organization orgEditHook
	// ---
	// summary: Update a hook
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: id
	//   in: path
	//   description: id of the hook to update
	//   type: integer
	//   format: int64
	//   required: true
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/EditHookOption"
	// responses:
	//   "200":
	//     "$ref": "#/responses/Hook"
	//   "404":
	//     "$ref": "#/responses/notFound"

	utils.EditOwnerHook(
		ctx,
		ctx.ContextUser,
		web.GetForm(ctx).(*api.EditHookOption),
		ctx.PathParamInt64("id"),
	)
}

// DeleteHook delete a hook of an organization
func DeleteHook(ctx *context.APIContext) {
	// swagger:operation DELETE /orgs/{org}/hooks/{id} organization orgDeleteHook
	// ---
	// summary: Delete a hook
	// produces:
	// - application/json
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: id
	//   in: path
	//   description: id of the hook to delete
	//   type: integer
	//   format: int64
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"

	utils.DeleteOwnerHook(
		ctx,
		ctx.ContextUser,
		ctx.PathParamInt64("id"),
	)
}
