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

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestLookupRedirect(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	repoID, err := repo_model.LookupRedirect(t.Context(), 2, "oldrepo1")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, repoID)

	_, err = repo_model.LookupRedirect(t.Context(), unittest.NonexistentID, "doesnotexist")
	assert.True(t, repo_model.IsErrRedirectNotExist(err))
}

func TestNewRedirect(t *testing.T) {
	// redirect to a completely new name
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.NoError(t, repo_model.NewRedirect(t.Context(), repo.OwnerID, repo.ID, repo.Name, "newreponame"))

	unittest.AssertExistsAndLoadBean(t, &repo_model.Redirect{
		OwnerID:        repo.OwnerID,
		LowerName:      repo.LowerName,
		RedirectRepoID: repo.ID,
	})
	unittest.AssertExistsAndLoadBean(t, &repo_model.Redirect{
		OwnerID:        repo.OwnerID,
		LowerName:      "oldrepo1",
		RedirectRepoID: repo.ID,
	})
}

func TestNewRedirect2(t *testing.T) {
	// redirect to previously used name
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.NoError(t, repo_model.NewRedirect(t.Context(), repo.OwnerID, repo.ID, repo.Name, "oldrepo1"))

	unittest.AssertExistsAndLoadBean(t, &repo_model.Redirect{
		OwnerID:        repo.OwnerID,
		LowerName:      repo.LowerName,
		RedirectRepoID: repo.ID,
	})
	unittest.AssertNotExistsBean(t, &repo_model.Redirect{
		OwnerID:        repo.OwnerID,
		LowerName:      "oldrepo1",
		RedirectRepoID: repo.ID,
	})
}

func TestNewRedirect3(t *testing.T) {
	// redirect for a previously-unredirected repo
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	assert.NoError(t, repo_model.NewRedirect(t.Context(), repo.OwnerID, repo.ID, repo.Name, "newreponame"))

	unittest.AssertExistsAndLoadBean(t, &repo_model.Redirect{
		OwnerID:        repo.OwnerID,
		LowerName:      repo.LowerName,
		RedirectRepoID: repo.ID,
	})
}
