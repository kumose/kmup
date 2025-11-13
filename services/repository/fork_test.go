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
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
)

func TestForkRepository(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// user 13 has already forked repo10
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 13})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 10})

	fork, err := ForkRepository(t.Context(), user, user, ForkRepoOptions{
		BaseRepo:    repo,
		Name:        "test",
		Description: "test",
	})
	assert.Nil(t, fork)
	assert.Error(t, err)
	assert.True(t, IsErrForkAlreadyExist(err))

	// user not reached maximum limit of repositories
	assert.False(t, repo_model.IsErrReachLimitOfRepo(err))

	// change AllowForkWithoutMaximumLimit to false for the test
	defer test.MockVariableValue(&setting.Repository.AllowForkWithoutMaximumLimit, false)()
	// user has reached maximum limit of repositories
	user.MaxRepoCreation = 0
	fork2, err := ForkRepository(t.Context(), user, user, ForkRepoOptions{
		BaseRepo:    repo,
		Name:        "test",
		Description: "test",
	})
	assert.Nil(t, fork2)
	assert.True(t, repo_model.IsErrReachLimitOfRepo(err))
}

func TestForkRepositoryCleanup(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// a successful fork
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo10 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 10})

	fork, err := ForkRepository(t.Context(), user2, user2, ForkRepoOptions{
		BaseRepo: repo10,
		Name:     "test",
	})
	assert.NoError(t, err)
	assert.NotNil(t, fork)

	exist, err := util.IsExist(repo_model.RepoPath(user2.Name, "test"))
	assert.NoError(t, err)
	assert.True(t, exist)

	err = DeleteRepositoryDirectly(t.Context(), fork.ID)
	assert.NoError(t, err)

	// a failed creating because some mock data
	// create the repository directory so that the creation will fail after database record created.
	assert.NoError(t, os.MkdirAll(repo_model.RepoPath(user2.Name, "test"), os.ModePerm))

	fork2, err := ForkRepository(t.Context(), user2, user2, ForkRepoOptions{
		BaseRepo: repo10,
		Name:     "test",
	})
	assert.Nil(t, fork2)
	assert.Error(t, err)

	// assert the cleanup is successful
	unittest.AssertNotExistsBean(t, &repo_model.Repository{OwnerName: user2.Name, Name: "test"})

	exist, err = util.IsExist(repo_model.RepoPath(user2.Name, "test"))
	assert.NoError(t, err)
	assert.False(t, exist)
}
