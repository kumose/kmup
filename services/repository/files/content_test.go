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

package files

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/contexttest"

	_ "github.com/kumose/kmup/models/actions"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

func TestGetContents(t *testing.T) {
	unittest.PrepareTestEnv(t)
	ctx, _ := contexttest.MockContext(t, "user2/repo1")
	ctx.SetPathParam("id", "1")
	contexttest.LoadRepo(t, ctx, 1)
	contexttest.LoadRepoCommit(t, ctx)
	contexttest.LoadUser(t, ctx, 2)
	contexttest.LoadGitRepo(t, ctx)

	// GetContentsOrList's behavior is fully tested in integration tests, so we don't need to test it here.

	t.Run("GetBlobBySHA", func(t *testing.T) {
		sha := "65f1bf27bc3bf70f64657658635e66094edbcb4d"
		ctx.SetPathParam("id", "1")
		ctx.SetPathParam("sha", sha)
		gbr, err := GetBlobBySHA(ctx.Repo.Repository, ctx.Repo.GitRepo, ctx.PathParam("sha"))
		expectedGBR := &api.GitBlobResponse{
			Content:  util.ToPointer("dHJlZSAyYTJmMWQ0NjcwNzI4YTJlMTAwNDllMzQ1YmQ3YTI3NjQ2OGJlYWI2CmF1dGhvciB1c2VyMSA8YWRkcmVzczFAZXhhbXBsZS5jb20+IDE0ODk5NTY0NzkgLTA0MDAKY29tbWl0dGVyIEV0aGFuIEtvZW5pZyA8ZXRoYW50a29lbmlnQGdtYWlsLmNvbT4gMTQ4OTk1NjQ3OSAtMDQwMAoKSW5pdGlhbCBjb21taXQK"),
			Encoding: util.ToPointer("base64"),
			URL:      "https://try.kmup.io/api/v1/repos/user2/repo1/git/blobs/65f1bf27bc3bf70f64657658635e66094edbcb4d",
			SHA:      "65f1bf27bc3bf70f64657658635e66094edbcb4d",
			Size:     180,
		}
		assert.NoError(t, err)
		assert.Equal(t, expectedGBR, gbr)
	})
}
