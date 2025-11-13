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

package setting

import (
	"net/http"
	"strconv"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/services/context"
)

func UpdatePreferences(ctx *context.Context) {
	type preferencesForm struct {
		CodeViewShowFileTree bool `json:"codeViewShowFileTree"`
	}
	form := &preferencesForm{}
	if err := json.NewDecoder(ctx.Req.Body).Decode(&form); err != nil {
		ctx.HTTPError(http.StatusBadRequest, "json decode failed")
		return
	}
	_ = user_model.SetUserSetting(ctx, ctx.Doer.ID, user_model.SettingsKeyCodeViewShowFileTree, strconv.FormatBool(form.CodeViewShowFileTree))
	ctx.JSONOK()
}
