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

package convert

import (
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func Test_ToPullReview(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	reviewer := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	review := unittest.AssertExistsAndLoadBean(t, &issues_model.Review{ID: 6})
	assert.Equal(t, reviewer.ID, review.ReviewerID)
	assert.Equal(t, issues_model.ReviewTypePending, review.Type)

	reviewList := []*issues_model.Review{review}

	t.Run("Anonymous User", func(t *testing.T) {
		prList, err := ToPullReviewList(t.Context(), reviewList, nil)
		assert.NoError(t, err)
		assert.Empty(t, prList)
	})

	t.Run("Reviewer Himself", func(t *testing.T) {
		prList, err := ToPullReviewList(t.Context(), reviewList, reviewer)
		assert.NoError(t, err)
		assert.Len(t, prList, 1)
	})

	t.Run("Other User", func(t *testing.T) {
		user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
		prList, err := ToPullReviewList(t.Context(), reviewList, user4)
		assert.NoError(t, err)
		assert.Empty(t, prList)
	})

	t.Run("Admin User", func(t *testing.T) {
		adminUser := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
		prList, err := ToPullReviewList(t.Context(), reviewList, adminUser)
		assert.NoError(t, err)
		assert.Len(t, prList, 1)
	})
}
