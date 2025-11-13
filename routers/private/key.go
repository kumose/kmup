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

package private

import (
	"net/http"

	asymkey_model "github.com/kumose/kmup/models/asymkey"
	"github.com/kumose/kmup/modules/private"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/services/context"
)

// UpdatePublicKeyInRepo update public key and deploy key updates
func UpdatePublicKeyInRepo(ctx *context.PrivateContext) {
	keyID := ctx.PathParamInt64("id")
	repoID := ctx.PathParamInt64("repoid")
	if err := asymkey_model.UpdatePublicKeyUpdated(ctx, keyID); err != nil {
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: err.Error(),
		})
		return
	}

	deployKey, err := asymkey_model.GetDeployKeyByRepo(ctx, keyID, repoID)
	if err != nil {
		if asymkey_model.IsErrDeployKeyNotExist(err) {
			ctx.PlainText(http.StatusOK, "success")
			return
		}
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: err.Error(),
		})
		return
	}
	deployKey.UpdatedUnix = timeutil.TimeStampNow()
	if err = asymkey_model.UpdateDeployKeyCols(ctx, deployKey, "updated_unix"); err != nil {
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: err.Error(),
		})
		return
	}

	ctx.PlainText(http.StatusOK, "success")
}

// AuthorizedPublicKeyByContent searches content as prefix (without comment part)
// and returns public key found.
func AuthorizedPublicKeyByContent(ctx *context.PrivateContext) {
	content := ctx.FormString("content")

	publicKey, err := asymkey_model.SearchPublicKeyByContent(ctx, content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: err.Error(),
		})
		return
	}

	authorizedString, err := asymkey_model.AuthorizedStringForKey(publicKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err:     err.Error(),
			UserMsg: "invalid public key",
		})
		return
	}
	ctx.PlainText(http.StatusOK, authorizedString)
}
