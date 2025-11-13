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

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateIssueWatch(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	assert.NoError(t, issues_model.CreateOrUpdateIssueWatch(t.Context(), 3, 1, true))
	iw := unittest.AssertExistsAndLoadBean(t, &issues_model.IssueWatch{UserID: 3, IssueID: 1})
	assert.True(t, iw.IsWatching)

	assert.NoError(t, issues_model.CreateOrUpdateIssueWatch(t.Context(), 1, 1, false))
	iw = unittest.AssertExistsAndLoadBean(t, &issues_model.IssueWatch{UserID: 1, IssueID: 1})
	assert.False(t, iw.IsWatching)
}

func TestGetIssueWatch(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	_, exists, err := issues_model.GetIssueWatch(t.Context(), 9, 1)
	assert.True(t, exists)
	assert.NoError(t, err)

	iw, exists, err := issues_model.GetIssueWatch(t.Context(), 2, 2)
	assert.True(t, exists)
	assert.NoError(t, err)
	assert.False(t, iw.IsWatching)

	_, exists, err = issues_model.GetIssueWatch(t.Context(), 3, 1)
	assert.False(t, exists)
	assert.NoError(t, err)
}

func TestGetIssueWatchers(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	iws, err := issues_model.GetIssueWatchers(t.Context(), 1, db.ListOptions{})
	assert.NoError(t, err)
	// Watcher is inactive, thus 0
	assert.Empty(t, iws)

	iws, err = issues_model.GetIssueWatchers(t.Context(), 2, db.ListOptions{})
	assert.NoError(t, err)
	// Watcher is explicit not watching
	assert.Empty(t, iws)

	iws, err = issues_model.GetIssueWatchers(t.Context(), 5, db.ListOptions{})
	assert.NoError(t, err)
	// Issue has no Watchers
	assert.Empty(t, iws)

	iws, err = issues_model.GetIssueWatchers(t.Context(), 7, db.ListOptions{})
	assert.NoError(t, err)
	// Issue has one watcher
	assert.Len(t, iws, 1)
}
