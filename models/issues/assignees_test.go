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

package issues_test

import (
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestUpdateAssignee(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// Fake issue with assignees
	issue, err := issues_model.GetIssueByID(t.Context(), 1)
	assert.NoError(t, err)

	err = issue.LoadAttributes(t.Context())
	assert.NoError(t, err)

	// Assign multiple users
	user2, err := user_model.GetUserByID(t.Context(), 2)
	assert.NoError(t, err)
	_, _, err = issues_model.ToggleIssueAssignee(t.Context(), issue, &user_model.User{ID: 1}, user2.ID)
	assert.NoError(t, err)

	org3, err := user_model.GetUserByID(t.Context(), 3)
	assert.NoError(t, err)
	_, _, err = issues_model.ToggleIssueAssignee(t.Context(), issue, &user_model.User{ID: 1}, org3.ID)
	assert.NoError(t, err)

	user1, err := user_model.GetUserByID(t.Context(), 1) // This user is already assigned (see the definition in fixtures), so running  UpdateAssignee should unassign him
	assert.NoError(t, err)
	_, _, err = issues_model.ToggleIssueAssignee(t.Context(), issue, &user_model.User{ID: 1}, user1.ID)
	assert.NoError(t, err)

	// Check if he got removed
	isAssigned, err := issues_model.IsUserAssignedToIssue(t.Context(), issue, user1)
	assert.NoError(t, err)
	assert.False(t, isAssigned)

	// Check if they're all there
	err = issue.LoadAssignees(t.Context())
	assert.NoError(t, err)

	var expectedAssignees []*user_model.User
	expectedAssignees = append(expectedAssignees, user2, org3)

	for in, assignee := range issue.Assignees {
		assert.Equal(t, assignee.ID, expectedAssignees[in].ID)
	}

	// Check if the user is assigned
	isAssigned, err = issues_model.IsUserAssignedToIssue(t.Context(), issue, user2)
	assert.NoError(t, err)
	assert.True(t, isAssigned)

	// This user should not be assigned
	isAssigned, err = issues_model.IsUserAssignedToIssue(t.Context(), issue, &user_model.User{ID: 4})
	assert.NoError(t, err)
	assert.False(t, isAssigned)
}

func TestMakeIDsFromAPIAssigneesToAdd(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	_ = unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	_ = unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	IDs, err := issues_model.MakeIDsFromAPIAssigneesToAdd(t.Context(), "", []string{""})
	assert.NoError(t, err)
	assert.Equal(t, []int64{}, IDs)

	_, err = issues_model.MakeIDsFromAPIAssigneesToAdd(t.Context(), "", []string{"none_existing_user"})
	assert.Error(t, err)

	IDs, err = issues_model.MakeIDsFromAPIAssigneesToAdd(t.Context(), "user1", []string{"user1"})
	assert.NoError(t, err)
	assert.Equal(t, []int64{1}, IDs)

	IDs, err = issues_model.MakeIDsFromAPIAssigneesToAdd(t.Context(), "user2", []string{""})
	assert.NoError(t, err)
	assert.Equal(t, []int64{2}, IDs)

	IDs, err = issues_model.MakeIDsFromAPIAssigneesToAdd(t.Context(), "", []string{"user1", "user2"})
	assert.NoError(t, err)
	assert.Equal(t, []int64{1, 2}, IDs)
}
