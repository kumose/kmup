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
	"fmt"
	"net/url"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
)

func TestGitSSHRedirect(t *testing.T) {
	onKmupRun(t, testGitSSHRedirect)
}

func testGitSSHRedirect(t *testing.T, u *url.URL) {
	apiTestContext := NewAPITestContext(t, "user2", "repo1", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)

	withKeyFile(t, "my-testing-key", func(keyFile string) {
		t.Run("CreateUserKey", doAPICreateUserKey(apiTestContext, "test-key", keyFile))

		testCases := []struct {
			testName string
			userName string
			repoName string
		}{
			{"Test untouched", "user2", "repo1"},
			{"Test renamed user", "olduser2", "repo1"},
			{"Test renamed repo", "user2", "oldrepo1"},
			{"Test renamed user and repo", "olduser2", "oldrepo1"},
		}

		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				cloneURL := createSSHUrl(fmt.Sprintf("%s/%s.git", tc.userName, tc.repoName), u)
				t.Run("Clone", doGitClone(t.TempDir(), cloneURL))
			})
		}
	})
}
