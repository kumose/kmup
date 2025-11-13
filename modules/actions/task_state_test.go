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

func TestFullSteps(t *testing.T) {
	tests := []struct {
		name string
		task *actions_model.ActionTask
		want []*actions_model.ActionTaskStep
	}{
		{
			name: "regular",
			task: &actions_model.ActionTask{
				Steps: []*actions_model.ActionTaskStep{
					{Status: actions_model.StatusSuccess, LogIndex: 10, LogLength: 80, Started: 10010, Stopped: 10090},
				},
				Status:    actions_model.StatusSuccess,
				Started:   10000,
				Stopped:   10100,
				LogLength: 100,
			},
			want: []*actions_model.ActionTaskStep{
				{Name: preStepName, Status: actions_model.StatusSuccess, LogIndex: 0, LogLength: 10, Started: 10000, Stopped: 10010},
				{Status: actions_model.StatusSuccess, LogIndex: 10, LogLength: 80, Started: 10010, Stopped: 10090},
				{Name: postStepName, Status: actions_model.StatusSuccess, LogIndex: 90, LogLength: 10, Started: 10090, Stopped: 10100},
			},
		},
		{
			name: "failed step",
			task: &actions_model.ActionTask{
				Steps: []*actions_model.ActionTaskStep{
					{Status: actions_model.StatusSuccess, LogIndex: 10, LogLength: 20, Started: 10010, Stopped: 10020},
					{Status: actions_model.StatusFailure, LogIndex: 30, LogLength: 60, Started: 10020, Stopped: 10090},
					{Status: actions_model.StatusCancelled, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
				},
				Status:    actions_model.StatusFailure,
				Started:   10000,
				Stopped:   10100,
				LogLength: 100,
			},
			want: []*actions_model.ActionTaskStep{
				{Name: preStepName, Status: actions_model.StatusSuccess, LogIndex: 0, LogLength: 10, Started: 10000, Stopped: 10010},
				{Status: actions_model.StatusSuccess, LogIndex: 10, LogLength: 20, Started: 10010, Stopped: 10020},
				{Status: actions_model.StatusFailure, LogIndex: 30, LogLength: 60, Started: 10020, Stopped: 10090},
				{Status: actions_model.StatusCancelled, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
				{Name: postStepName, Status: actions_model.StatusFailure, LogIndex: 90, LogLength: 10, Started: 10090, Stopped: 10100},
			},
		},
		{
			name: "first step is running",
			task: &actions_model.ActionTask{
				Steps: []*actions_model.ActionTaskStep{
					{Status: actions_model.StatusRunning, LogIndex: 10, LogLength: 80, Started: 10010, Stopped: 0},
				},
				Status:    actions_model.StatusRunning,
				Started:   10000,
				Stopped:   10100,
				LogLength: 100,
			},
			want: []*actions_model.ActionTaskStep{
				{Name: preStepName, Status: actions_model.StatusSuccess, LogIndex: 0, LogLength: 10, Started: 10000, Stopped: 10010},
				{Status: actions_model.StatusRunning, LogIndex: 10, LogLength: 80, Started: 10010, Stopped: 0},
				{Name: postStepName, Status: actions_model.StatusWaiting, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
			},
		},
		{
			name: "first step has canceled",
			task: &actions_model.ActionTask{
				Steps: []*actions_model.ActionTaskStep{
					{Status: actions_model.StatusCancelled, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
				},
				Status:    actions_model.StatusFailure,
				Started:   10000,
				Stopped:   10100,
				LogLength: 100,
			},
			want: []*actions_model.ActionTaskStep{
				{Name: preStepName, Status: actions_model.StatusFailure, LogIndex: 0, LogLength: 100, Started: 10000, Stopped: 10100},
				{Status: actions_model.StatusCancelled, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
				{Name: postStepName, Status: actions_model.StatusFailure, LogIndex: 100, LogLength: 0, Started: 10100, Stopped: 10100},
			},
		},
		{
			name: "empty steps",
			task: &actions_model.ActionTask{
				Steps:     []*actions_model.ActionTaskStep{},
				Status:    actions_model.StatusSuccess,
				Started:   10000,
				Stopped:   10100,
				LogLength: 100,
			},
			want: []*actions_model.ActionTaskStep{
				{Name: preStepName, Status: actions_model.StatusSuccess, LogIndex: 0, LogLength: 100, Started: 10000, Stopped: 10100},
				{Name: postStepName, Status: actions_model.StatusSuccess, LogIndex: 100, LogLength: 0, Started: 10100, Stopped: 10100},
			},
		},
		{
			name: "all steps finished but task is running",
			task: &actions_model.ActionTask{
				Steps: []*actions_model.ActionTaskStep{
					{Status: actions_model.StatusSuccess, LogIndex: 10, LogLength: 80, Started: 10010, Stopped: 10090},
				},
				Status:    actions_model.StatusRunning,
				Started:   10000,
				Stopped:   0,
				LogLength: 100,
			},
			want: []*actions_model.ActionTaskStep{
				{Name: preStepName, Status: actions_model.StatusSuccess, LogIndex: 0, LogLength: 10, Started: 10000, Stopped: 10010},
				{Status: actions_model.StatusSuccess, LogIndex: 10, LogLength: 80, Started: 10010, Stopped: 10090},
				{Name: postStepName, Status: actions_model.StatusRunning, LogIndex: 90, LogLength: 10, Started: 10090, Stopped: 0},
			},
		},
		{
			name: "skipped task",
			task: &actions_model.ActionTask{
				Steps: []*actions_model.ActionTaskStep{
					{Status: actions_model.StatusSkipped, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
				},
				Status:    actions_model.StatusSkipped,
				Started:   0,
				Stopped:   0,
				LogLength: 0,
			},
			want: []*actions_model.ActionTaskStep{
				{Name: preStepName, Status: actions_model.StatusSkipped, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
				{Status: actions_model.StatusSkipped, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
				{Name: postStepName, Status: actions_model.StatusSkipped, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
			},
		},
		{
			name: "first step is skipped",
			task: &actions_model.ActionTask{
				Steps: []*actions_model.ActionTaskStep{
					{Status: actions_model.StatusSkipped, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
					{Status: actions_model.StatusSuccess, LogIndex: 10, LogLength: 80, Started: 10010, Stopped: 10090},
				},
				Status:    actions_model.StatusSuccess,
				Started:   10000,
				Stopped:   10100,
				LogLength: 100,
			},
			want: []*actions_model.ActionTaskStep{
				{Name: preStepName, Status: actions_model.StatusSuccess, LogIndex: 0, LogLength: 10, Started: 10000, Stopped: 10010},
				{Status: actions_model.StatusSkipped, LogIndex: 0, LogLength: 0, Started: 0, Stopped: 0},
				{Status: actions_model.StatusSuccess, LogIndex: 10, LogLength: 80, Started: 10010, Stopped: 10090},
				{Name: postStepName, Status: actions_model.StatusSuccess, LogIndex: 90, LogLength: 10, Started: 10090, Stopped: 10100},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FullSteps(tt.task), "FullSteps(%v)", tt.task)
		})
	}
}
