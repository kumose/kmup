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

package gitrepo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	path string
}

func (r *mockRepository) RelativePath() string {
	return r.path
}

func TestRepoGetDivergingCommits(t *testing.T) {
	repo := &mockRepository{path: "repo1_bare"}
	do, err := GetDivergingCommits(t.Context(), repo, "master", "branch2")
	assert.NoError(t, err)
	assert.Equal(t, &DivergeObject{
		Ahead:  1,
		Behind: 5,
	}, do)

	do, err = GetDivergingCommits(t.Context(), repo, "master", "master")
	assert.NoError(t, err)
	assert.Equal(t, &DivergeObject{
		Ahead:  0,
		Behind: 0,
	}, do)

	do, err = GetDivergingCommits(t.Context(), repo, "master", "test")
	assert.NoError(t, err)
	assert.Equal(t, &DivergeObject{
		Ahead:  0,
		Behind: 2,
	}, do)
}
