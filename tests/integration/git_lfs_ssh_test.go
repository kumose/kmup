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
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/routers/common"
	"github.com/kumose/kmup/services/context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitLFSSSH(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		localRepoForUpload := filepath.Join(t.TempDir(), "test-upload")
		localRepoForDownload := filepath.Join(t.TempDir(), "test-download")
		apiTestContext := NewAPITestContext(t, "user2", "repo1", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)

		var mu sync.Mutex
		var routerCalls []string
		web.RouteMock(common.RouterMockPointCommonLFS, func(ctx *context.Base) {
			mu.Lock()
			routerCalls = append(routerCalls, ctx.Req.Method+" "+ctx.Req.URL.Path)
			mu.Unlock()
		})

		withKeyFile(t, "my-testing-key", func(keyFile string) {
			t.Run("CreateUserKey", doAPICreateUserKey(apiTestContext, "test-key", keyFile))
			cloneURL := createSSHUrl(apiTestContext.GitPath(), u)
			t.Run("CloneOrigin", doGitClone(localRepoForUpload, cloneURL))

			cfg, err := setting.CfgProvider.PrepareSaving()
			require.NoError(t, err)
			cfg.Section("server").Key("LFS_ALLOW_PURE_SSH").SetValue("true")
			setting.LFS.AllowPureSSH = true
			require.NoError(t, cfg.Save())

			_, _, cmdErr := gitcmd.NewCommand("config", "lfs.sshtransfer", "always").
				WithDir(localRepoForUpload).
				RunStdString(t.Context())
			assert.NoError(t, cmdErr)
			pushedFiles := lfsCommitAndPushTest(t, localRepoForUpload, 10)

			t.Run("CloneLFS", doGitClone(localRepoForDownload, cloneURL))
			content, err := os.ReadFile(filepath.Join(localRepoForDownload, pushedFiles[0]))
			assert.NoError(t, err)
			assert.Len(t, content, 10)
		})

		countBatch := slices.ContainsFunc(routerCalls, func(s string) bool {
			return strings.Contains(s, "POST /api/internal/repo/user2/repo1.git/info/lfs/objects/batch")
		})
		countUpload := slices.ContainsFunc(routerCalls, func(s string) bool {
			return strings.Contains(s, "PUT /api/internal/repo/user2/repo1.git/info/lfs/objects/")
		})
		countDownload := slices.ContainsFunc(routerCalls, func(s string) bool {
			return strings.Contains(s, "GET /api/internal/repo/user2/repo1.git/info/lfs/objects/")
		})
		nonAPIRequests := slices.ContainsFunc(routerCalls, func(s string) bool {
			fields := strings.Fields(s)
			return !strings.HasPrefix(fields[1], "/api/")
		})
		assert.NotZero(t, countBatch)
		assert.NotZero(t, countUpload)
		assert.NotZero(t, countDownload)
		assert.Zero(t, nonAPIRequests)
	})
}
