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

package markup

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/models/user"
	kmup_context "github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestRenderHelperMention(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	userPublic := "user1"
	userPrivate := "user31"
	userLimited := "user33"
	userNoSuch := "no-such-user"

	unittest.AssertCount(t, &user.User{Name: userPublic}, 1)
	unittest.AssertCount(t, &user.User{Name: userPrivate}, 1)
	unittest.AssertCount(t, &user.User{Name: userLimited}, 1)
	unittest.AssertCount(t, &user.User{Name: userNoSuch}, 0)

	// when using general context, use user's visibility to check
	assert.True(t, FormalRenderHelperFuncs().IsUsernameMentionable(t.Context(), userPublic))
	assert.False(t, FormalRenderHelperFuncs().IsUsernameMentionable(t.Context(), userLimited))
	assert.False(t, FormalRenderHelperFuncs().IsUsernameMentionable(t.Context(), userPrivate))
	assert.False(t, FormalRenderHelperFuncs().IsUsernameMentionable(t.Context(), userNoSuch))

	// when using web context, use user.IsUserVisibleToViewer to check
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	base := kmup_context.NewBaseContextForTest(httptest.NewRecorder(), req)
	kmupCtx := kmup_context.NewWebContext(base, &contexttest.MockRender{}, nil)

	assert.True(t, FormalRenderHelperFuncs().IsUsernameMentionable(kmupCtx, userPublic))
	assert.False(t, FormalRenderHelperFuncs().IsUsernameMentionable(kmupCtx, userPrivate))

	kmupCtx.Doer, err = user.GetUserByName(t.Context(), userPrivate)
	assert.NoError(t, err)
	assert.True(t, FormalRenderHelperFuncs().IsUsernameMentionable(kmupCtx, userPublic))
	assert.True(t, FormalRenderHelperFuncs().IsUsernameMentionable(kmupCtx, userPrivate))
}
