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

package repo

import (
	"net/http"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	git_module "github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestViewHomeSubmoduleRedirect(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockContext(t, "/user2/repo1/src/branch/master/test-submodule")
	submodule := git_module.NewCommitSubmoduleFile("/user2/repo1", "test-submodule", "../repo-other", "any-ref-id")
	handleRepoViewSubmodule(ctx, submodule)
	assert.Equal(t, http.StatusSeeOther, ctx.Resp.WrittenStatus())
	assert.Equal(t, "/user2/repo-other/tree/any-ref-id", ctx.Resp.Header().Get("Location"))

	ctx, _ = contexttest.MockContext(t, "/user2/repo1/src/branch/master/test-submodule")
	submodule = git_module.NewCommitSubmoduleFile("/user2/repo1", "test-submodule", "https://other/user2/repo-other.git", "any-ref-id")
	handleRepoViewSubmodule(ctx, submodule)
	// do not auto-redirect for external URLs, to avoid open redirect or phishing
	assert.Equal(t, http.StatusNotFound, ctx.Resp.WrittenStatus())

	ctx, respWriter := contexttest.MockContext(t, "/user2/repo1/src/branch/master/test-submodule?only_content=true")
	submodule = git_module.NewCommitSubmoduleFile("/user2/repo1", "test-submodule", "../repo-other", "any-ref-id")
	handleRepoViewSubmodule(ctx, submodule)
	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())
	assert.Equal(t, `<a href="/user2/repo-other/tree/any-ref-id">/user2/repo-other/tree/any-ref-id</a>`, respWriter.Body.String())
}
