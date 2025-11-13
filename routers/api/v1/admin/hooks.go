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
	"errors"
	"net/http"

	"github.com/kumose/kmup/models/webhook"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
	webhook_service "github.com/kumose/kmup/services/webhook"
)

// ListHooks list system's webhooks
func ListHooks(ctx *context.APIContext) {
	// swagger:operation GET /admin/hooks admin adminListHooks
	// ---
	// summary: List system's webhooks
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
	// - type: string
	//   enum:
	//     - system
	//     - default
	//     - all
	//   description: system, default or both kinds of webhooks
	//   name: type
	//   default: system
	//   in: query
	//
	// responses:
	//   "200":
	//     "$ref": "#/responses/HookList"

	// for compatibility the default value is true
	isSystemWebhook := optional.Some(true)
	typeValue := ctx.FormString("type")
	switch typeValue {
	case "default":
		isSystemWebhook = optional.Some(false)
	case "all":
		isSystemWebhook = optional.None[bool]()
	}

	sysHooks, err := webhook.GetSystemOrDefaultWebhooks(ctx, isSystemWebhook)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	hooks := make([]*api.Hook, len(sysHooks))
	for i, hook := range sysHooks {
		h, err := webhook_service.ToHook(setting.AppURL+"/-/admin", hook)
		if err != nil {
			ctx.APIErrorInternal(err)
			return
		}
		hooks[i] = h
	}
	ctx.JSON(http.StatusOK, hooks)
}

// GetHook get an organization's hook by id
func GetHook(ctx *context.APIContext) {
	// swagger:operation GET /admin/hooks/{id} admin adminGetHook
	// ---
	// summary: Get a hook
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the hook to get
	//   type: integer
	//   format: int64
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/Hook"

	hookID := ctx.PathParamInt64("id")
	hook, err := webhook.GetSystemOrDefaultWebhook(ctx, hookID)
	if err != nil {
		if errors.Is(err, util.ErrNotExist) {
			ctx.APIErrorNotFound()
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}
	h, err := webhook_service.ToHook("/-/admin/", hook)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	ctx.JSON(http.StatusOK, h)
}

// CreateHook create a hook for an organization
func CreateHook(ctx *context.APIContext) {
	// swagger:operation POST /admin/hooks admin adminCreateHook
	// ---
	// summary: Create a hook
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreateHookOption"
	// responses:
	//   "201":
	//     "$ref": "#/responses/Hook"

	form := web.GetForm(ctx).(*api.CreateHookOption)

	utils.AddSystemHook(ctx, form)
}

// EditHook modify a hook of a repository
func EditHook(ctx *context.APIContext) {
	// swagger:operation PATCH /admin/hooks/{id} admin adminEditHook
	// ---
	// summary: Update a hook
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
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

	form := web.GetForm(ctx).(*api.EditHookOption)

	// TODO in body params
	hookID := ctx.PathParamInt64("id")
	utils.EditSystemHook(ctx, form, hookID)
}

// DeleteHook delete a system hook
func DeleteHook(ctx *context.APIContext) {
	// swagger:operation DELETE /admin/hooks/{id} admin adminDeleteHook
	// ---
	// summary: Delete a hook
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the hook to delete
	//   type: integer
	//   format: int64
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"

	hookID := ctx.PathParamInt64("id")
	if err := webhook.DeleteDefaultSystemWebhook(ctx, hookID); err != nil {
		if errors.Is(err, util.ErrNotExist) {
			ctx.APIErrorNotFound()
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}
