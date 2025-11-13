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

package common

import (
	"net/http"
	"testing"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/web/middleware"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRequestPriority(t *testing.T) {
	type test struct {
		Name         string
		User         *user_model.User
		RoutePattern string
		Expected     Priority
	}

	cases := []test{
		{
			Name:     "Logged In",
			User:     &user_model.User{},
			Expected: HighPriority,
		},
		{
			Name:         "Sign In",
			RoutePattern: "/user/login",
			Expected:     DefaultPriority,
		},
		{
			Name:         "Repo Home",
			RoutePattern: "/{username}/{reponame}",
			Expected:     DefaultPriority,
		},
		{
			Name:         "User Repo",
			RoutePattern: "/{username}/{reponame}/src/branch/main",
			Expected:     LowPriority,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, _ := contexttest.MockContext(t, "")

			if tc.User != nil {
				data := middleware.GetContextData(ctx)
				data[middleware.ContextDataKeySignedUser] = tc.User
			}

			rctx := chi.RouteContext(ctx)
			rctx.RoutePatterns = []string{tc.RoutePattern}

			assert.Exactly(t, tc.Expected, requestPriority(ctx))
		})
	}
}

func TestRenderServiceUnavailable(t *testing.T) {
	t.Run("HTML", func(t *testing.T) {
		ctx, resp := contexttest.MockContext(t, "")
		ctx.Req.Header.Set("Accept", "text/html")

		renderServiceUnavailable(resp, ctx.Req)
		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
		assert.Contains(t, resp.Header().Get("Content-Type"), "text/html")

		body := resp.Body.String()
		assert.Contains(t, body, `lang="en-US"`)
		assert.Contains(t, body, "503 Service Unavailable")
	})

	t.Run("plain", func(t *testing.T) {
		ctx, resp := contexttest.MockContext(t, "")
		ctx.Req.Header.Set("Accept", "text/plain")

		renderServiceUnavailable(resp, ctx.Req)
		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
		assert.Contains(t, resp.Header().Get("Content-Type"), "text/plain")

		body := resp.Body.String()
		assert.Contains(t, body, "503 Service Unavailable")
	})
}
