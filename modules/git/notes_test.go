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

func TestGetNotes(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := OpenRepository(t.Context(), bareRepo1Path)
	assert.NoError(t, err)
	defer bareRepo1.Close()

	note := Note{}
	err = GetNote(t.Context(), bareRepo1, "95bb4d39648ee7e325106df01a621c530863a653", &note)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Note contents\n"), note.Message)
	assert.Equal(t, "Vladimir Panteleev", note.Commit.Author.Name)
}

func TestGetNestedNotes(t *testing.T) {
	repoPath := filepath.Join(testReposDir, "repo3_notes")
	repo, err := OpenRepository(t.Context(), repoPath)
	assert.NoError(t, err)
	defer repo.Close()

	note := Note{}
	err = GetNote(t.Context(), repo, "3e668dbfac39cbc80a9ff9c61eb565d944453ba4", &note)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Note 2"), note.Message)
	err = GetNote(t.Context(), repo, "ba0a96fa63532d6c5087ecef070b0250ed72fa47", &note)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Note 1"), note.Message)
}

func TestGetNonExistentNotes(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := OpenRepository(t.Context(), bareRepo1Path)
	assert.NoError(t, err)
	defer bareRepo1.Close()

	note := Note{}
	err = GetNote(t.Context(), bareRepo1, "non_existent_sha", &note)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &ErrNotExist{})
}
