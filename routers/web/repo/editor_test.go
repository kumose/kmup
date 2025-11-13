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
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/gitrepo"

	"github.com/stretchr/testify/assert"
)

func TestEditorUtils(t *testing.T) {
	unittest.PrepareTestEnv(t)
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	t.Run("getUniquePatchBranchName", func(t *testing.T) {
		branchName := getUniquePatchBranchName(t.Context(), "user2", repo)
		assert.Equal(t, "user2-patch-1", branchName)
	})
	t.Run("getClosestParentWithFiles", func(t *testing.T) {
		gitRepo, _ := gitrepo.OpenRepository(t.Context(), repo)
		defer gitRepo.Close()
		treePath := getClosestParentWithFiles(gitRepo, "sub-home-md-img-check", "docs/foo/bar")
		assert.Equal(t, "docs", treePath)
		treePath = getClosestParentWithFiles(gitRepo, "sub-home-md-img-check", "any/other")
		assert.Empty(t, treePath)
	})
}
