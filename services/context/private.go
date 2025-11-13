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

package context

import (
	"context"
	"net/http"
	"time"

	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/web"
	web_types "github.com/kumose/kmup/modules/web/types"
)

// PrivateContext represents a context for private routes
type PrivateContext struct {
	*Base
	Override context.Context

	Repo *Repository
}

func init() {
	web.RegisterResponseStatusProvider[*PrivateContext](func(req *http.Request) web_types.ResponseStatusProvider {
		return req.Context().Value(privateContextKey).(*PrivateContext)
	})
}

func (ctx *PrivateContext) Deadline() (deadline time.Time, ok bool) {
	if ctx.Override != nil {
		return ctx.Override.Deadline()
	}
	return ctx.Base.Deadline()
}

func (ctx *PrivateContext) Done() <-chan struct{} {
	if ctx.Override != nil {
		return ctx.Override.Done()
	}
	return ctx.Base.Done()
}

func (ctx *PrivateContext) Err() error {
	if ctx.Override != nil {
		return ctx.Override.Err()
	}
	return ctx.Base.Err()
}

type privateContextKeyType struct{}

var privateContextKey privateContextKeyType

func GetPrivateContext(req *http.Request) *PrivateContext {
	return req.Context().Value(privateContextKey).(*PrivateContext)
}

func PrivateContexter() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			base := NewBaseContext(w, req)
			ctx := &PrivateContext{Base: base}
			ctx.SetContextValue(privateContextKey, ctx)
			next.ServeHTTP(ctx.Resp, ctx.Req)
		})
	}
}

// OverrideContext overrides the underlying request context for Done() etc.
// This function should be used when there is a need for work to continue even if the request has been cancelled.
// Primarily this affects hook/post-receive and hook/proc-receive both of which need to continue working even if
// the underlying request has timed out from the ssh/http push
func OverrideContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// We now need to override the request context as the base for our work because even if the request is cancelled we have to continue this work
			ctx := GetPrivateContext(req)
			var finished func()
			ctx.Override, _, finished = process.GetManager().AddTypedContext(graceful.GetManager().HammerContext(), "PrivateContext: "+ctx.Req.RequestURI, process.RequestProcessType, true)
			defer finished()
			next.ServeHTTP(ctx.Resp, ctx.Req)
		})
	}
}
