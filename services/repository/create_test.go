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
	"os"
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepositoryDirectly(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// a successful creating repository
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	createdRepo, err := CreateRepositoryDirectly(t.Context(), user2, user2, CreateRepoOptions{
		Name: "created-repo",
	}, true)
	assert.NoError(t, err)
	assert.NotNil(t, createdRepo)

	exist, err := util.IsExist(repo_model.RepoPath(user2.Name, createdRepo.Name))
	assert.NoError(t, err)
	assert.True(t, exist)

	unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{OwnerName: user2.Name, Name: createdRepo.Name})

	err = DeleteRepositoryDirectly(t.Context(), createdRepo.ID)
	assert.NoError(t, err)

	// a failed creating because some mock data
	// create the repository directory so that the creation will fail after database record created.
	assert.NoError(t, os.MkdirAll(repo_model.RepoPath(user2.Name, createdRepo.Name), os.ModePerm))

	createdRepo2, err := CreateRepositoryDirectly(t.Context(), user2, user2, CreateRepoOptions{
		Name: "created-repo",
	}, true)
	assert.Nil(t, createdRepo2)
	assert.Error(t, err)

	// assert the cleanup is successful
	unittest.AssertNotExistsBean(t, &repo_model.Repository{OwnerName: user2.Name, Name: createdRepo.Name})

	exist, err = util.IsExist(repo_model.RepoPath(user2.Name, createdRepo.Name))
	assert.NoError(t, err)
	assert.False(t, exist)
}
