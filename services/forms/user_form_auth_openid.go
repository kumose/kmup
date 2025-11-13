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

// SignInOpenIDForm form for signing in with OpenID
type SignInOpenIDForm struct {
	Openid   string `binding:"Required;MaxSize(256)"`
	Remember bool
}

// Validate validates the fields
func (f *SignInOpenIDForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

// SignUpOpenIDForm form for signin up with OpenID
type SignUpOpenIDForm struct {
	UserName string `binding:"Required;Username;MaxSize(40)"`
	Email    string `binding:"Required;Email;MaxSize(254)"`
}

// Validate validates the fields
func (f *SignUpOpenIDForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

// ConnectOpenIDForm form for connecting an existing account to an OpenID URI
type ConnectOpenIDForm struct {
	UserName string `binding:"Required;MaxSize(254)"`
	Password string `binding:"Required;MaxSize(255)"`
}

// Validate validates the fields
func (f *ConnectOpenIDForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}
