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

	actions_model "github.com/kumose/kmup/models/actions"

	"github.com/stretchr/testify/assert"
)

func Test_jobStatusResolver_Resolve(t *testing.T) {
	tests := []struct {
		name string
		jobs actions_model.ActionJobList
		want map[int64]actions_model.Status
	}{
		{
			name: "no blocked",
			jobs: actions_model.ActionJobList{
				{ID: 1, JobID: "1", Status: actions_model.StatusWaiting, Needs: []string{}},
				{ID: 2, JobID: "2", Status: actions_model.StatusWaiting, Needs: []string{}},
				{ID: 3, JobID: "3", Status: actions_model.StatusWaiting, Needs: []string{}},
			},
			want: map[int64]actions_model.Status{},
		},
		{
			name: "single blocked",
			jobs: actions_model.ActionJobList{
				{ID: 1, JobID: "1", Status: actions_model.StatusSuccess, Needs: []string{}},
				{ID: 2, JobID: "2", Status: actions_model.StatusBlocked, Needs: []string{"1"}},
				{ID: 3, JobID: "3", Status: actions_model.StatusWaiting, Needs: []string{}},
			},
			want: map[int64]actions_model.Status{
				2: actions_model.StatusWaiting,
			},
		},
		{
			name: "multiple blocked",
			jobs: actions_model.ActionJobList{
				{ID: 1, JobID: "1", Status: actions_model.StatusSuccess, Needs: []string{}},
				{ID: 2, JobID: "2", Status: actions_model.StatusBlocked, Needs: []string{"1"}},
				{ID: 3, JobID: "3", Status: actions_model.StatusBlocked, Needs: []string{"1"}},
			},
			want: map[int64]actions_model.Status{
				2: actions_model.StatusWaiting,
				3: actions_model.StatusWaiting,
			},
		},
		{
			name: "chain blocked",
			jobs: actions_model.ActionJobList{
				{ID: 1, JobID: "1", Status: actions_model.StatusFailure, Needs: []string{}},
				{ID: 2, JobID: "2", Status: actions_model.StatusBlocked, Needs: []string{"1"}},
				{ID: 3, JobID: "3", Status: actions_model.StatusBlocked, Needs: []string{"2"}},
			},
			want: map[int64]actions_model.Status{
				2: actions_model.StatusSkipped,
				3: actions_model.StatusSkipped,
			},
		},
		{
			name: "loop need",
			jobs: actions_model.ActionJobList{
				{ID: 1, JobID: "1", Status: actions_model.StatusBlocked, Needs: []string{"3"}},
				{ID: 2, JobID: "2", Status: actions_model.StatusBlocked, Needs: []string{"1"}},
				{ID: 3, JobID: "3", Status: actions_model.StatusBlocked, Needs: []string{"2"}},
			},
			want: map[int64]actions_model.Status{},
		},
		{
			name: "`if` is not empty and all jobs in `needs` completed successfully",
			jobs: actions_model.ActionJobList{
				{ID: 1, JobID: "job1", Status: actions_model.StatusSuccess, Needs: []string{}},
				{ID: 2, JobID: "job2", Status: actions_model.StatusBlocked, Needs: []string{"job1"}, WorkflowPayload: []byte(
					`
name: test
on: push
jobs:
  job2:
    runs-on: ubuntu-latest
    needs: job1
    if: ${{ always() && needs.job1.result == 'success' }}
    steps:
      - run: echo "will be checked by act_runner"
`)},
			},
			want: map[int64]actions_model.Status{2: actions_model.StatusWaiting},
		},
		{
			name: "`if` is not empty and not all jobs in `needs` completed successfully",
			jobs: actions_model.ActionJobList{
				{ID: 1, JobID: "job1", Status: actions_model.StatusFailure, Needs: []string{}},
				{ID: 2, JobID: "job2", Status: actions_model.StatusBlocked, Needs: []string{"job1"}, WorkflowPayload: []byte(
					`
name: test
on: push
jobs:
  job2:
    runs-on: ubuntu-latest
    needs: job1
    if: ${{ always() && needs.job1.result == 'failure' }}
    steps:
      - run: echo "will be checked by act_runner"
`)},
			},
			want: map[int64]actions_model.Status{2: actions_model.StatusWaiting},
		},
		{
			name: "`if` is empty and not all jobs in `needs` completed successfully",
			jobs: actions_model.ActionJobList{
				{ID: 1, JobID: "job1", Status: actions_model.StatusFailure, Needs: []string{}},
				{ID: 2, JobID: "job2", Status: actions_model.StatusBlocked, Needs: []string{"job1"}, WorkflowPayload: []byte(
					`
name: test
on: push
jobs:
  job2:
    runs-on: ubuntu-latest
    needs: job1
    steps:
      - run: echo "should be skipped"
`)},
			},
			want: map[int64]actions_model.Status{2: actions_model.StatusSkipped},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newJobStatusResolver(tt.jobs, nil)
			assert.Equal(t, tt.want, r.Resolve(t.Context()))
		})
	}
}
