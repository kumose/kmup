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

package repository

import (
	"testing"

	"github.com/kumose/kmup/models/db"
	git_model "github.com/kumose/kmup/models/git"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestSyncRepoBranches(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	_, err := db.GetEngine(t.Context()).ID(1).Update(&repo_model.Repository{ObjectFormatName: "bad-fmt"})
	assert.NoError(t, db.TruncateBeans(t.Context(), &git_model.Branch{}))
	assert.NoError(t, err)
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.Equal(t, "bad-fmt", repo.ObjectFormatName)
	_, err = SyncRepoBranches(t.Context(), 1, 0)
	assert.NoError(t, err)
	repo = unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.Equal(t, "sha1", repo.ObjectFormatName)
	branch, err := git_model.GetBranch(t.Context(), 1, "master")
	assert.NoError(t, err)
	assert.Equal(t, "master", branch.Name)
}
