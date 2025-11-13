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

package forms

import (
	"net/http"

	"github.com/kumose/kmup/modules/web/middleware"
	"github.com/kumose/kmup/services/context"

	"github.com/kumose-go/chi/binding"
)

type PackageCleanupRuleForm struct {
	ID            int64
	Enabled       bool
	Type          string `binding:"Required;In(alpine,arch,cargo,chef,composer,conan,conda,container,cran,debian,generic,go,helm,maven,npm,nuget,pub,pypi,rpm,rubygems,swift,vagrant)"`
	KeepCount     int    `binding:"In(0,1,5,10,25,50,100)"`
	KeepPattern   string `binding:"RegexPattern"`
	RemoveDays    int    `binding:"In(0,7,14,30,60,90,180)"`
	RemovePattern string `binding:"RegexPattern"`
	MatchFullName bool
	Action        string `binding:"Required;In(save,remove)"`
}

func (f *PackageCleanupRuleForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}
