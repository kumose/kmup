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

package integration

import (
	"net/http"
	"net/url"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIGetRawFileOrLFS(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// Test with raw file
	req := NewRequest(t, "GET", "/api/v1/repos/user2/repo1/media/README.md")
	resp := MakeRequest(t, req, http.StatusOK)
	assert.Equal(t, "# repo1\n\nDescription for repo1", resp.Body.String())

	// Test with LFS
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		createLFSTestRepository(t, "repo-lfs-test")
		httpContext := NewAPITestContext(t, "user2", "repo-lfs-test", auth_model.AccessTokenScopeWriteRepository)
		t.Run("repo-lfs-test", func(t *testing.T) {
			u.Path = httpContext.GitPath()
			dstPath := t.TempDir()

			u.Path = httpContext.GitPath()
			u.User = url.UserPassword("user2", userPassword)

			t.Run("Clone", doGitClone(dstPath, u))

			dstPath2 := t.TempDir()

			t.Run("Partial Clone", doPartialGitClone(dstPath2, u))

			lfs := lfsCommitAndPushTest(t, dstPath, testFileSizeSmall)[0]

			reqLFS := NewRequest(t, "GET", "/api/v1/repos/user2/repo-lfs-test/media/"+lfs).AddTokenAuth(httpContext.Token)
			respLFS := MakeRequestNilResponseRecorder(t, reqLFS, http.StatusOK)
			assert.Equal(t, testFileSizeSmall, respLFS.Length)
		})
	})
}
