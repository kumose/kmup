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

package integration

import (
	"testing"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/util"
	repo_service "github.com/kumose/kmup/services/repository"
	user_service "github.com/kumose/kmup/services/user"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestEphemeralActionsRunnerDeletion(t *testing.T) {
	t.Run("ByTaskCompletion", testEphemeralActionsRunnerDeletionByTaskCompletion)
	t.Run("ByRepository", testEphemeralActionsRunnerDeletionByRepository)
	t.Run("ByUser", testEphemeralActionsRunnerDeletionByUser)
}

// Test that the ephemeral runner is deleted when the task is finished
func testEphemeralActionsRunnerDeletionByTaskCompletion(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	_, err := actions_model.GetRunnerByID(t.Context(), 34350)
	assert.NoError(t, err)

	task := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionTask{ID: 52})
	assert.Equal(t, actions_model.StatusRunning, task.Status)

	task.Status = actions_model.StatusSuccess
	err = actions_model.UpdateTask(t.Context(), task, "status")
	assert.NoError(t, err)

	_, err = actions_model.GetRunnerByID(t.Context(), 34350)
	assert.ErrorIs(t, err, util.ErrNotExist)
}

func testEphemeralActionsRunnerDeletionByRepository(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	_, err := actions_model.GetRunnerByID(t.Context(), 34350)
	assert.NoError(t, err)

	task := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionTask{ID: 52})
	assert.Equal(t, actions_model.StatusRunning, task.Status)

	err = repo_service.DeleteRepositoryDirectly(t.Context(), task.RepoID, true)
	assert.NoError(t, err)

	_, err = actions_model.GetRunnerByID(t.Context(), 34350)
	assert.ErrorIs(t, err, util.ErrNotExist)
}

// Test that the ephemeral runner is deleted when a user is deleted
func testEphemeralActionsRunnerDeletionByUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	_, err := actions_model.GetRunnerByID(t.Context(), 34350)
	assert.NoError(t, err)

	task := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionTask{ID: 52})
	assert.Equal(t, actions_model.StatusRunning, task.Status)

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	err = user_service.DeleteUser(t.Context(), user, true)
	assert.NoError(t, err)

	_, err = actions_model.GetRunnerByID(t.Context(), 34350)
	assert.ErrorIs(t, err, util.ErrNotExist)
}
