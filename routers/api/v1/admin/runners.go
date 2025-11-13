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
	"github.com/kumose/kmup/routers/api/v1/shared"
	"github.com/kumose/kmup/services/context"
)

// https://docs.github.com/en/rest/actions/self-hosted-runners?apiVersion=2022-11-28#create-a-registration-token-for-an-organization

// GetRegistrationToken returns the token to register global runners
func GetRegistrationToken(ctx *context.APIContext) {
	// swagger:operation GET /admin/runners/registration-token admin adminGetRunnerRegistrationToken
	// ---
	// summary: Get an global actions runner registration token
	// produces:
	// - application/json
	// parameters:
	// responses:
	//   "200":
	//     "$ref": "#/responses/RegistrationToken"

	shared.GetRegistrationToken(ctx, 0, 0)
}

// CreateRegistrationToken returns the token to register global runners
func CreateRegistrationToken(ctx *context.APIContext) {
	// swagger:operation POST /admin/actions/runners/registration-token admin adminCreateRunnerRegistrationToken
	// ---
	// summary: Get an global actions runner registration token
	// produces:
	// - application/json
	// parameters:
	// responses:
	//   "200":
	//     "$ref": "#/responses/RegistrationToken"

	shared.GetRegistrationToken(ctx, 0, 0)
}

// ListRunners get all runners
func ListRunners(ctx *context.APIContext) {
	// swagger:operation GET /admin/actions/runners admin getAdminRunners
	// ---
	// summary: Get all runners
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/definitions/ActionRunnersResponse"
	//   "400":
	//     "$ref": "#/responses/error"
	//   "404":
	//     "$ref": "#/responses/notFound"
	shared.ListRunners(ctx, 0, 0)
}

// GetRunner get an global runner
func GetRunner(ctx *context.APIContext) {
	// swagger:operation GET /admin/actions/runners/{runner_id} admin getAdminRunner
	// ---
	// summary: Get an global runner
	// produces:
	// - application/json
	// parameters:
	// - name: runner_id
	//   in: path
	//   description: id of the runner
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/definitions/ActionRunner"
	//   "400":
	//     "$ref": "#/responses/error"
	//   "404":
	//     "$ref": "#/responses/notFound"
	shared.GetRunner(ctx, 0, 0, ctx.PathParamInt64("runner_id"))
}

// DeleteRunner delete an global runner
func DeleteRunner(ctx *context.APIContext) {
	// swagger:operation DELETE /admin/actions/runners/{runner_id} admin deleteAdminRunner
	// ---
	// summary: Delete an global runner
	// produces:
	// - application/json
	// parameters:
	// - name: runner_id
	//   in: path
	//   description: id of the runner
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     description: runner has been deleted
	//   "400":
	//     "$ref": "#/responses/error"
	//   "404":
	//     "$ref": "#/responses/notFound"
	shared.DeleteRunner(ctx, 0, 0, ctx.PathParamInt64("runner_id"))
}
