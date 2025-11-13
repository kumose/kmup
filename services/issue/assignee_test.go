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

package issue

import (
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestDeleteNotPassedAssignee(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// Fake issue with assignees
	issue, err := issues_model.GetIssueByID(t.Context(), 1)
	assert.NoError(t, err)

	err = issue.LoadAttributes(t.Context())
	assert.NoError(t, err)

	assert.Len(t, issue.Assignees, 1)

	user1, err := user_model.GetUserByID(t.Context(), 1) // This user is already assigned (see the definition in fixtures), so running  UpdateAssignee should unassign him
	assert.NoError(t, err)

	// Check if he got removed
	isAssigned, err := issues_model.IsUserAssignedToIssue(t.Context(), issue, user1)
	assert.NoError(t, err)
	assert.True(t, isAssigned)

	// Clean everyone
	err = DeleteNotPassedAssignee(t.Context(), issue, user1, []*user_model.User{})
	assert.NoError(t, err)
	assert.Empty(t, issue.Assignees)

	// Reload to check they're gone
	issue.ResetAttributesLoaded()
	assert.NoError(t, issue.LoadAssignees(t.Context()))
	assert.Empty(t, issue.Assignees)
	assert.Empty(t, issue.Assignee)
}
