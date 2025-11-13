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

package doctor

import (
	"testing"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/log"

	"github.com/stretchr/testify/assert"
)

func Test_fixUnfinishedRunStatus(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	fixUnfinishedRunStatus(t.Context(), log.GetLogger(log.DEFAULT), true)

	// check if the run is cancelled by id
	run := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRun{ID: 805})
	assert.Equal(t, actions_model.StatusCancelled, run.Status)
}
