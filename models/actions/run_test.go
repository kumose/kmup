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

package actions

import (
	"testing"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestUpdateRepoRunsNumbers(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// update the number to a wrong one, the original is 3
	_, err := db.GetEngine(t.Context()).ID(4).Cols("num_closed_action_runs").Update(&repo_model.Repository{
		NumClosedActionRuns: 2,
	})
	assert.NoError(t, err)

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 4})
	assert.Equal(t, 4, repo.NumActionRuns)
	assert.Equal(t, 2, repo.NumClosedActionRuns)

	// now update will correct them, only num_actionr_runs and num_closed_action_runs should be updated
	err = UpdateRepoRunsNumbers(t.Context(), repo)
	assert.NoError(t, err)
	repo = unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 4})
	assert.Equal(t, 5, repo.NumActionRuns)
	assert.Equal(t, 3, repo.NumClosedActionRuns)
}
