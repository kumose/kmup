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

package pull_test

import (
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	pull_service "github.com/kumose/kmup/services/pull"

	"github.com/stretchr/testify/assert"
)

func TestRepoGetReviewers(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// test public repo
	repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	ctx := t.Context()
	reviewers, err := pull_service.GetReviewers(ctx, repo1, 2, 0)
	assert.NoError(t, err)
	if assert.Len(t, reviewers, 1) {
		assert.ElementsMatch(t, []int64{2}, []int64{reviewers[0].ID})
	}

	// should not include doer and remove the poster
	reviewers, err = pull_service.GetReviewers(ctx, repo1, 11, 2)
	assert.NoError(t, err)
	assert.Empty(t, reviewers)

	// should not include PR poster, if PR poster would be otherwise eligible
	reviewers, err = pull_service.GetReviewers(ctx, repo1, 11, 4)
	assert.NoError(t, err)
	assert.Len(t, reviewers, 1)

	// test private user repo
	repo2 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})

	reviewers, err = pull_service.GetReviewers(ctx, repo2, 2, 4)
	assert.NoError(t, err)
	assert.Len(t, reviewers, 1)
	assert.EqualValues(t, 2, reviewers[0].ID)

	// test private org repo
	repo3 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})

	reviewers, err = pull_service.GetReviewers(ctx, repo3, 2, 1)
	assert.NoError(t, err)
	assert.Len(t, reviewers, 2)

	reviewers, err = pull_service.GetReviewers(ctx, repo3, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, reviewers, 1)
}

func TestRepoGetReviewerTeams(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo2 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	teams, err := pull_service.GetReviewerTeams(t.Context(), repo2)
	assert.NoError(t, err)
	assert.Empty(t, teams)

	repo3 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})
	teams, err = pull_service.GetReviewerTeams(t.Context(), repo3)
	assert.NoError(t, err)
	assert.Len(t, teams, 2)
}
