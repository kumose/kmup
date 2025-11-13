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

package label

import (
	"github.com/kumose/kmup/modules/label"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/forms"
)

func GetLabelEditForm(ctx *context.Context) *forms.CreateLabelForm {
	form := web.GetForm(ctx).(*forms.CreateLabelForm)
	if ctx.HasError() {
		ctx.JSONError(ctx.Data["ErrorMsg"].(string))
		return nil
	}
	var err error
	form.Color, err = label.NormalizeColor(form.Color)
	if err != nil {
		ctx.JSONError(ctx.Tr("repo.issues.label_color_invalid"))
		return nil
	}
	return form
}
