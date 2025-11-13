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

package repo_test

import (
	"testing"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestStarRepo(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	unittest.AssertNotExistsBean(t, &repo_model.Star{UID: user.ID, RepoID: repo.ID})
	assert.NoError(t, repo_model.StarRepo(t.Context(), user, repo, true))
	unittest.AssertExistsAndLoadBean(t, &repo_model.Star{UID: user.ID, RepoID: repo.ID})
	assert.NoError(t, repo_model.StarRepo(t.Context(), user, repo, true))
	unittest.AssertExistsAndLoadBean(t, &repo_model.Star{UID: user.ID, RepoID: repo.ID})
	assert.NoError(t, repo_model.StarRepo(t.Context(), user, repo, false))
	unittest.AssertNotExistsBean(t, &repo_model.Star{UID: user.ID, RepoID: repo.ID})
}

func TestIsStaring(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	assert.True(t, repo_model.IsStaring(t.Context(), 2, 4))
	assert.False(t, repo_model.IsStaring(t.Context(), 3, 4))
}

func TestRepository_GetStargazers(t *testing.T) {
	// repo with stargazers
	assert.NoError(t, unittest.PrepareTestDatabase())
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 4})
	gazers, err := repo_model.GetStargazers(t.Context(), repo, db.ListOptions{Page: 0})
	assert.NoError(t, err)
	if assert.Len(t, gazers, 1) {
		assert.Equal(t, int64(2), gazers[0].ID)
	}
}

func TestRepository_GetStargazers2(t *testing.T) {
	// repo with stargazers
	assert.NoError(t, unittest.PrepareTestDatabase())
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})
	gazers, err := repo_model.GetStargazers(t.Context(), repo, db.ListOptions{Page: 0})
	assert.NoError(t, err)
	assert.Empty(t, gazers)
}

func TestClearRepoStars(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	unittest.AssertNotExistsBean(t, &repo_model.Star{UID: user.ID, RepoID: repo.ID})
	assert.NoError(t, repo_model.StarRepo(t.Context(), user, repo, true))
	unittest.AssertExistsAndLoadBean(t, &repo_model.Star{UID: user.ID, RepoID: repo.ID})
	assert.NoError(t, repo_model.StarRepo(t.Context(), user, repo, false))
	unittest.AssertNotExistsBean(t, &repo_model.Star{UID: user.ID, RepoID: repo.ID})
	assert.NoError(t, repo_model.ClearRepoStars(t.Context(), repo.ID))
	unittest.AssertNotExistsBean(t, &repo_model.Star{UID: user.ID, RepoID: repo.ID})

	gazers, err := repo_model.GetStargazers(t.Context(), repo, db.ListOptions{Page: 0})
	assert.NoError(t, err)
	assert.Empty(t, gazers)
}
