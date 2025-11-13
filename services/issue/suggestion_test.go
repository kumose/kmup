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

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/optional"

	"github.com/stretchr/testify/assert"
)

func Test_Suggestion(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	testCases := []struct {
		keyword         string
		isPull          optional.Option[bool]
		expectedIndexes []int64
	}{
		{
			keyword:         "",
			expectedIndexes: []int64{5, 1, 4, 2, 3},
		},
		{
			keyword:         "1",
			expectedIndexes: []int64{1},
		},
		{
			keyword:         "issue",
			expectedIndexes: []int64{4, 1, 2, 3},
		},
		{
			keyword:         "pull",
			expectedIndexes: []int64{5},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.keyword, func(t *testing.T) {
			issues, err := GetSuggestion(t.Context(), repo1, testCase.isPull, testCase.keyword)
			assert.NoError(t, err)

			issueIndexes := make([]int64, 0, len(issues))
			for _, issue := range issues {
				issueIndexes = append(issueIndexes, issue.Index)
			}
			assert.Equal(t, testCase.expectedIndexes, issueIndexes)
		})
	}
}
