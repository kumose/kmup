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

package v1_18

import (
	"testing"

	"github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/migrations/base"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateOpenMilestoneCounts(t *testing.T) {
	type ExpectedMilestone issues.Milestone

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(issues.Milestone), new(ExpectedMilestone), new(issues.Issue))
	defer deferable()
	if x == nil || t.Failed() {
		return
	}

	if err := UpdateOpenMilestoneCounts(x); err != nil {
		assert.NoError(t, err)
		return
	}

	expected := []ExpectedMilestone{}
	if err := x.Table("expected_milestone").Asc("id").Find(&expected); !assert.NoError(t, err) {
		return
	}

	got := []issues.Milestone{}
	if err := x.Table("milestone").Asc("id").Find(&got); !assert.NoError(t, err) {
		return
	}

	for i, e := range expected {
		got := got[i]
		assert.Equal(t, e.ID, got.ID)
		assert.Equal(t, e.NumIssues, got.NumIssues)
		assert.Equal(t, e.NumClosedIssues, got.NumClosedIssues)
	}
}
