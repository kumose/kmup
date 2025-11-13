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
	"github.com/kumose/kmup/modules/timeutil"

	"github.com/stretchr/testify/assert"
)

func TestContentHistory(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	dbCtx := t.Context()
	timeStampNow := timeutil.TimeStampNow()

	_ = issues_model.SaveIssueContentHistory(dbCtx, 1, 10, 0, timeStampNow, "i-a", true)
	_ = issues_model.SaveIssueContentHistory(dbCtx, 1, 10, 0, timeStampNow.Add(2), "i-b", false)
	_ = issues_model.SaveIssueContentHistory(dbCtx, 1, 10, 0, timeStampNow.Add(7), "i-c", false)

	_ = issues_model.SaveIssueContentHistory(dbCtx, 1, 10, 100, timeStampNow, "c-a", true)
	_ = issues_model.SaveIssueContentHistory(dbCtx, 1, 10, 100, timeStampNow.Add(5), "c-b", false)
	_ = issues_model.SaveIssueContentHistory(dbCtx, 1, 10, 100, timeStampNow.Add(20), "c-c", false)
	_ = issues_model.SaveIssueContentHistory(dbCtx, 1, 10, 100, timeStampNow.Add(50), "c-d", false)
	_ = issues_model.SaveIssueContentHistory(dbCtx, 1, 10, 100, timeStampNow.Add(51), "c-e", false)

	h1, _ := issues_model.GetIssueContentHistoryByID(dbCtx, 1)
	assert.EqualValues(t, 1, h1.ID)

	m, _ := issues_model.QueryIssueContentHistoryEditedCountMap(dbCtx, 10)
	assert.Equal(t, 3, m[0])
	assert.Equal(t, 5, m[100])

	/*
		we can not have this test with real `User` now, because we can not depend on `User` model (circle-import), so there is no `user` table
		when the refactor of models are done, this test will be possible to be run then with a real `User` model.
	*/
	type User struct {
		ID       int64
		Name     string
		FullName string
	}
	_ = db.GetEngine(dbCtx).Sync(&User{})

	list1, _ := issues_model.FetchIssueContentHistoryList(dbCtx, 10, 0)
	assert.Len(t, list1, 3)
	list2, _ := issues_model.FetchIssueContentHistoryList(dbCtx, 10, 100)
	assert.Len(t, list2, 5)

	hasHistory1, _ := issues_model.HasIssueContentHistory(dbCtx, 10, 0)
	assert.True(t, hasHistory1)
	hasHistory2, _ := issues_model.HasIssueContentHistory(dbCtx, 10, 1)
	assert.False(t, hasHistory2)

	h6, h6Prev, _ := issues_model.GetIssueContentHistoryAndPrev(dbCtx, 10, 6)
	assert.EqualValues(t, 6, h6.ID)
	assert.EqualValues(t, 5, h6Prev.ID)

	// soft-delete
	_ = issues_model.SoftDeleteIssueContentHistory(dbCtx, 5)
	h6, h6Prev, _ = issues_model.GetIssueContentHistoryAndPrev(dbCtx, 10, 6)
	assert.EqualValues(t, 6, h6.ID)
	assert.EqualValues(t, 4, h6Prev.ID)

	// only keep 3 history revisions for comment_id=100, the first and the last should never be deleted
	issues_model.KeepLimitedContentHistory(dbCtx, 10, 100, 3)
	list1, _ = issues_model.FetchIssueContentHistoryList(dbCtx, 10, 0)
	assert.Len(t, list1, 3)
	list2, _ = issues_model.FetchIssueContentHistoryList(dbCtx, 10, 100)
	assert.Len(t, list2, 3)
	assert.EqualValues(t, 8, list2[0].HistoryID)
	assert.EqualValues(t, 7, list2[1].HistoryID)
	assert.EqualValues(t, 4, list2[2].HistoryID)
}

func TestHasIssueContentHistoryForCommentOnly(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	_ = db.TruncateBeans(t.Context(), &issues_model.ContentHistory{})

	hasHistory1, _ := issues_model.HasIssueContentHistory(t.Context(), 10, 0)
	assert.False(t, hasHistory1)
	hasHistory2, _ := issues_model.HasIssueContentHistory(t.Context(), 10, 100)
	assert.False(t, hasHistory2)

	_ = issues_model.SaveIssueContentHistory(t.Context(), 1, 10, 100, timeutil.TimeStampNow(), "c-a", true)
	_ = issues_model.SaveIssueContentHistory(t.Context(), 1, 10, 100, timeutil.TimeStampNow().Add(5), "c-b", false)

	hasHistory1, _ = issues_model.HasIssueContentHistory(t.Context(), 10, 0)
	assert.False(t, hasHistory1)
	hasHistory2, _ = issues_model.HasIssueContentHistory(t.Context(), 10, 100)
	assert.True(t, hasHistory2)
}
