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

func TestRepository_GetRefs(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := OpenRepository(t.Context(), bareRepo1Path)
	assert.NoError(t, err)
	defer bareRepo1.Close()

	refs, err := bareRepo1.GetRefs()

	assert.NoError(t, err)
	assert.Len(t, refs, 6)

	expectedRefs := []string{
		BranchPrefix + "branch1",
		BranchPrefix + "branch2",
		BranchPrefix + "master",
		TagPrefix + "test",
		TagPrefix + "signed-tag",
		NotesRef,
	}

	for _, ref := range refs {
		assert.Contains(t, expectedRefs, ref.Name)
	}
}

func TestRepository_GetRefsFiltered(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := OpenRepository(t.Context(), bareRepo1Path)
	assert.NoError(t, err)
	defer bareRepo1.Close()

	refs, err := bareRepo1.GetRefsFiltered(TagPrefix)

	assert.NoError(t, err)
	if assert.Len(t, refs, 2) {
		assert.Equal(t, TagPrefix+"signed-tag", refs[0].Name)
		assert.Equal(t, "tag", refs[0].Type)
		assert.Equal(t, "36f97d9a96457e2bab511db30fe2db03893ebc64", refs[0].Object.String())
		assert.Equal(t, TagPrefix+"test", refs[1].Name)
		assert.Equal(t, "tag", refs[1].Type)
		assert.Equal(t, "3ad28a9149a2864384548f3d17ed7f38014c9e8a", refs[1].Object.String())
	}
}
