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

package user

import (
	"github.com/kumose/kmup/routers/api/v1/shared"
	"github.com/kumose/kmup/services/context"
)

// https://docs.github.com/en/rest/actions/self-hosted-runners?apiVersion=2022-11-28#create-a-registration-token-for-an-organization

// GetRegistrationToken returns the token to register user runners
func GetRegistrationToken(ctx *context.APIContext) {
	// swagger:operation GET /user/actions/runners/registration-token user userGetRunnerRegistrationToken
	// ---
	// summary: Get an user's actions runner registration token
	// produces:
	// - application/json
	// parameters:
	// responses:
	//   "200":
	//     "$ref": "#/responses/RegistrationToken"

	shared.GetRegistrationToken(ctx, ctx.Doer.ID, 0)
}

// CreateRegistrationToken returns the token to register user runners
func CreateRegistrationToken(ctx *context.APIContext) {
	// swagger:operation POST /user/actions/runners/registration-token user userCreateRunnerRegistrationToken
	// ---
	// summary: Get an user's actions runner registration token
	// produces:
	// - application/json
	// parameters:
	// responses:
	//   "200":
	//     "$ref": "#/responses/RegistrationToken"

	shared.GetRegistrationToken(ctx, ctx.Doer.ID, 0)
}

// ListRunners get user-level runners
func ListRunners(ctx *context.APIContext) {
	// swagger:operation GET /user/actions/runners user getUserRunners
	// ---
	// summary: Get user-level runners
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/definitions/ActionRunnersResponse"
	//   "400":
	//     "$ref": "#/responses/error"
	//   "404":
	//     "$ref": "#/responses/notFound"
	shared.ListRunners(ctx, ctx.Doer.ID, 0)
}

// GetRunner get an user-level runner
func GetRunner(ctx *context.APIContext) {
	// swagger:operation GET /user/actions/runners/{runner_id} user getUserRunner
	// ---
	// summary: Get an user-level runner
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
	shared.GetRunner(ctx, ctx.Doer.ID, 0, ctx.PathParamInt64("runner_id"))
}

// DeleteRunner delete an user-level runner
func DeleteRunner(ctx *context.APIContext) {
	// swagger:operation DELETE /user/actions/runners/{runner_id} user deleteUserRunner
	// ---
	// summary: Delete an user-level runner
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
	shared.DeleteRunner(ctx, ctx.Doer.ID, 0, ctx.PathParamInt64("runner_id"))
}
