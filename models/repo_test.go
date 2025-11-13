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

package models

import (
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestCheckRepoStats(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	assert.NoError(t, CheckRepoStats(t.Context()))
}

func TestDoctorUserStarNum(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	assert.NoError(t, DoctorUserStarNum(t.Context()))
}

func Test_repoStatsCorrectIssueNumComments(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	issue2 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 2})
	assert.NotNil(t, issue2)
	assert.Equal(t, 0, issue2.NumComments) // the fixture data is wrong, but we don't fix it here

	assert.NoError(t, repoStatsCorrectIssueNumComments(t.Context(), 2))
	// reload the issue
	issue2 = unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 2})
	assert.Equal(t, 1, issue2.NumComments)
}
