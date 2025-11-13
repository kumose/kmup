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

package cron

import (
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTaskToScheduler(t *testing.T) {
	assert.Empty(t, scheduler.Jobs())
	defer scheduler.Clear()

	// no seconds
	err := addTaskToScheduler(&Task{
		Name: "task 1",
		config: &BaseConfig{
			Schedule: "5 4 * * *",
		},
	})
	assert.NoError(t, err)
	jobs := scheduler.Jobs()
	assert.Len(t, jobs, 1)
	assert.Equal(t, "task 1", jobs[0].Tags()[0])
	assert.Equal(t, "5 4 * * *", jobs[0].Tags()[1])

	// with seconds
	err = addTaskToScheduler(&Task{
		Name: "task 2",
		config: &BaseConfig{
			Schedule: "30 5 4 * * *",
		},
	})
	assert.NoError(t, err)
	jobs = scheduler.Jobs() // the item order is not guaranteed, so we need to sort it before "assert"
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Tags()[0] < jobs[j].Tags()[0]
	})
	assert.Len(t, jobs, 2)
	assert.Equal(t, "task 2", jobs[1].Tags()[0])
	assert.Equal(t, "30 5 4 * * *", jobs[1].Tags()[1])
}

func TestScheduleHasSeconds(t *testing.T) {
	tests := []struct {
		schedule  string
		hasSecond bool
	}{
		{"* * * * * *", true},
		{"* * * * *", false},
		{"5 4 * * *", false},
		{"5 4 * * *", false},
		{"5,8 4 * * *", false},
		{"*   *   *  * * *", true},
		{"5,8 4   *  *   *", false},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.hasSecond, scheduleHasSeconds(test.schedule))
		})
	}
}
