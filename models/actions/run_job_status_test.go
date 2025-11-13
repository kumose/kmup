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

	"github.com/stretchr/testify/assert"
)

func TestAggregateJobStatus(t *testing.T) {
	testStatuses := func(expected Status, statuses []Status) {
		t.Helper()
		var jobs []*ActionRunJob
		for _, v := range statuses {
			jobs = append(jobs, &ActionRunJob{Status: v})
		}
		actual := AggregateJobStatus(jobs)
		if !assert.Equal(t, expected, actual) {
			var statusStrings []string
			for _, s := range statuses {
				statusStrings = append(statusStrings, s.String())
			}
			t.Errorf("AggregateJobStatus(%v) = %v, want %v", statusStrings, statusNames[actual], statusNames[expected])
		}
	}

	cases := []struct {
		statuses []Status
		expected Status
	}{
		// unknown cases, maybe it shouldn't happen in real world
		{[]Status{}, StatusUnknown},
		{[]Status{StatusUnknown, StatusSuccess}, StatusUnknown},
		{[]Status{StatusUnknown, StatusSkipped}, StatusUnknown},
		{[]Status{StatusUnknown, StatusFailure}, StatusFailure},
		{[]Status{StatusUnknown, StatusCancelled}, StatusCancelled},
		{[]Status{StatusUnknown, StatusWaiting}, StatusWaiting},
		{[]Status{StatusUnknown, StatusRunning}, StatusRunning},
		{[]Status{StatusUnknown, StatusBlocked}, StatusBlocked},

		// success with other status
		{[]Status{StatusSuccess}, StatusSuccess},
		{[]Status{StatusSuccess, StatusSkipped}, StatusSuccess}, // skipped doesn't affect success
		{[]Status{StatusSuccess, StatusFailure}, StatusFailure},
		{[]Status{StatusSuccess, StatusCancelled}, StatusCancelled},
		{[]Status{StatusSuccess, StatusWaiting}, StatusWaiting},
		{[]Status{StatusSuccess, StatusRunning}, StatusRunning},
		{[]Status{StatusSuccess, StatusBlocked}, StatusBlocked},

		// any cancelled, then cancelled
		{[]Status{StatusCancelled}, StatusCancelled},
		{[]Status{StatusCancelled, StatusSuccess}, StatusCancelled},
		{[]Status{StatusCancelled, StatusSkipped}, StatusCancelled},
		{[]Status{StatusCancelled, StatusFailure}, StatusCancelled},
		{[]Status{StatusCancelled, StatusWaiting}, StatusCancelled},
		{[]Status{StatusCancelled, StatusRunning}, StatusCancelled},
		{[]Status{StatusCancelled, StatusBlocked}, StatusCancelled},

		// failure with other status, usually fail fast, but "running" wins to match GitHub's behavior
		// another reason that we can't make "failure" wins over "running": it would cause a weird behavior that user cannot cancel a workflow or get current running workflows correctly by filter after a job fail.
		{[]Status{StatusFailure}, StatusFailure},
		{[]Status{StatusFailure, StatusSuccess}, StatusFailure},
		{[]Status{StatusFailure, StatusSkipped}, StatusFailure},
		{[]Status{StatusFailure, StatusCancelled}, StatusCancelled},
		{[]Status{StatusFailure, StatusWaiting}, StatusWaiting},
		{[]Status{StatusFailure, StatusRunning}, StatusRunning},
		{[]Status{StatusFailure, StatusBlocked}, StatusFailure},

		// skipped with other status
		// "all skipped" is also considered as "mergeable" by "services/actions.toCommitStatus", the same as GitHub
		{[]Status{StatusSkipped}, StatusSkipped},
		{[]Status{StatusSkipped, StatusSuccess}, StatusSuccess},
		{[]Status{StatusSkipped, StatusFailure}, StatusFailure},
		{[]Status{StatusSkipped, StatusCancelled}, StatusCancelled},
		{[]Status{StatusSkipped, StatusWaiting}, StatusWaiting},
		{[]Status{StatusSkipped, StatusRunning}, StatusRunning},
		{[]Status{StatusSkipped, StatusBlocked}, StatusBlocked},
	}

	for _, c := range cases {
		testStatuses(c.expected, c.statuses)
	}
}
