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
package feed_test

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/routers/web/feed"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

func TestCheckGetOrgFeedsAsOrgMember(t *testing.T) {
	unittest.PrepareTestEnv(t)
	t.Run("OrgMember", func(t *testing.T) {
		ctx, resp := contexttest.MockContext(t, "org3.atom")
		ctx.ContextUser = unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 3})
		contexttest.LoadUser(t, ctx, 2)
		ctx.IsSigned = true
		feed.ShowUserFeedAtom(ctx)
		assert.Contains(t, resp.Body.String(), "<entry>") // Should contain 1 private entry
	})
	t.Run("NonOrgMember", func(t *testing.T) {
		ctx, resp := contexttest.MockContext(t, "org3.atom")
		ctx.ContextUser = unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 3})
		contexttest.LoadUser(t, ctx, 5)
		ctx.IsSigned = true
		feed.ShowUserFeedAtom(ctx)
		assert.NotContains(t, resp.Body.String(), "<entry>") // Should not contain any entries
	})
}
