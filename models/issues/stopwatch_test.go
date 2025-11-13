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

func TestCancelStopwatch(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	issue1 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})

	ok, err := issues_model.CancelStopwatch(t.Context(), user1, issue1)
	assert.NoError(t, err)
	assert.True(t, ok)
	unittest.AssertNotExistsBean(t, &issues_model.Stopwatch{UserID: user1.ID, IssueID: issue1.ID})
	unittest.AssertExistsAndLoadBean(t, &issues_model.Comment{Type: issues_model.CommentTypeCancelTracking, PosterID: user1.ID, IssueID: issue1.ID})

	ok, err = issues_model.CancelStopwatch(t.Context(), user1, issue1)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestStopwatchExists(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	assert.True(t, issues_model.StopwatchExists(t.Context(), 1, 1))
	assert.False(t, issues_model.StopwatchExists(t.Context(), 1, 2))
}

func TestHasUserStopwatch(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	exists, sw, _, err := issues_model.HasUserStopwatch(t.Context(), 1)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, int64(1), sw.ID)

	exists, _, _, err = issues_model.HasUserStopwatch(t.Context(), 3)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCreateOrStopIssueStopwatch(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
	issue1 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
	issue3 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 3})

	// create a new stopwatch
	ok, err := issues_model.CreateIssueStopwatch(t.Context(), user4, issue1)
	assert.NoError(t, err)
	assert.True(t, ok)
	unittest.AssertExistsAndLoadBean(t, &issues_model.Stopwatch{UserID: user4.ID, IssueID: issue1.ID})
	// should not create a second stopwatch for the same issue
	ok, err = issues_model.CreateIssueStopwatch(t.Context(), user4, issue1)
	assert.NoError(t, err)
	assert.False(t, ok)
	// on a different issue, it will finish the existing stopwatch and create a new one
	ok, err = issues_model.CreateIssueStopwatch(t.Context(), user4, issue3)
	assert.NoError(t, err)
	assert.True(t, ok)
	unittest.AssertNotExistsBean(t, &issues_model.Stopwatch{UserID: user4.ID, IssueID: issue1.ID})
	unittest.AssertExistsAndLoadBean(t, &issues_model.Stopwatch{UserID: user4.ID, IssueID: issue3.ID})

	// user2 already has a stopwatch in test fixture
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	issue2 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 2})
	ok, err = issues_model.FinishIssueStopwatch(t.Context(), user2, issue2)
	assert.NoError(t, err)
	assert.True(t, ok)
	unittest.AssertNotExistsBean(t, &issues_model.Stopwatch{UserID: user2.ID, IssueID: issue2.ID})
	unittest.AssertExistsAndLoadBean(t, &issues_model.TrackedTime{UserID: user2.ID, IssueID: issue2.ID})
	ok, err = issues_model.FinishIssueStopwatch(t.Context(), user2, issue2)
	assert.NoError(t, err)
	assert.False(t, ok)
}
