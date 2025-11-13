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

package git

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubTree_Issue29101(t *testing.T) {
	repo, err := OpenRepository(t.Context(), filepath.Join(testReposDir, "repo1_bare"))
	assert.NoError(t, err)
	defer repo.Close()

	commit, err := repo.GetCommit("ce064814f4a0d337b333e646ece456cd39fab612")
	assert.NoError(t, err)

	// old code could produce a different error if called multiple times
	for range 10 {
		_, err = commit.SubTree("file1.txt")
		assert.Error(t, err)
		assert.True(t, IsErrNotExist(err))
	}
}

func Test_GetTreePathLatestCommit(t *testing.T) {
	repo, err := OpenRepository(t.Context(), filepath.Join(testReposDir, "repo6_blame"))
	assert.NoError(t, err)
	defer repo.Close()

	commitID, err := repo.GetBranchCommitID("master")
	assert.NoError(t, err)
	assert.Equal(t, "544d8f7a3b15927cddf2299b4b562d6ebd71b6a7", commitID)

	commit, err := repo.GetTreePathLatestCommit("master", "blame.txt")
	assert.NoError(t, err)
	assert.NotNil(t, commit)
	assert.Equal(t, "45fb6cbc12f970b04eacd5cd4165edd11c8d7376", commit.ID.String())
}
