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

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/gitrepo"

	"github.com/stretchr/testify/assert"
)

func TestGetDirectorySize(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	repo, err := repo_model.GetRepositoryByID(t.Context(), 1)
	assert.NoError(t, err)
	size, err := gitrepo.CalcRepositorySize(repo)
	assert.NoError(t, err)
	repo.Size = 8165 // real size on the disk
	assert.Equal(t, repo.Size, size)
}
