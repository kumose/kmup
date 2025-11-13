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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepository_GetCodeActivityStats(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := OpenRepository(t.Context(), bareRepo1Path)
	assert.NoError(t, err)
	defer bareRepo1.Close()

	timeFrom, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00+00:00")
	assert.NoError(t, err)

	code, err := bareRepo1.GetCodeActivityStats(timeFrom, "")
	assert.NoError(t, err)
	assert.NotNil(t, code)

	assert.EqualValues(t, 10, code.CommitCount)
	assert.EqualValues(t, 3, code.AuthorCount)
	assert.EqualValues(t, 10, code.CommitCountInAllBranches)
	assert.EqualValues(t, 10, code.Additions)
	assert.EqualValues(t, 1, code.Deletions)
	assert.Len(t, code.Authors, 3)
	assert.Equal(t, "tris.git@shoddynet.org", code.Authors[1].Email)
	assert.EqualValues(t, 3, code.Authors[1].Commits)
	assert.EqualValues(t, 5, code.Authors[0].Commits)
}
