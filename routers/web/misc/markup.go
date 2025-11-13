// Copyright 2014 The Gogs Authors. All rights reserved.
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

package misc

import (
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/routers/common"
	"github.com/kumose/kmup/services/context"
)

// Markup render markup document to HTML
func Markup(ctx *context.Context) {
	form := web.GetForm(ctx).(*api.MarkupOption)
	mode := util.Iif(form.Wiki, "wiki", form.Mode) //nolint:staticcheck // form.Wiki is deprecated
	common.RenderMarkup(ctx.Base, ctx.Repo, mode, form.Text, form.Context, form.FilePath)
}
