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

package pull

import (
	"testing"

	git_model "github.com/kumose/kmup/models/git"
	"github.com/kumose/kmup/modules/commitstatus"

	"github.com/stretchr/testify/assert"
)

func TestMergeRequiredContextsCommitStatus(t *testing.T) {
	cases := []struct {
		commitStatuses   []*git_model.CommitStatus
		requiredContexts []string
		expected         commitstatus.CommitStatusState
	}{
		{
			commitStatuses:   []*git_model.CommitStatus{},
			requiredContexts: []string{},
			expected:         commitstatus.CommitStatusPending,
		},
		{
			commitStatuses: []*git_model.CommitStatus{
				{Context: "Build xxx", State: commitstatus.CommitStatusSkipped},
			},
			requiredContexts: []string{"Build*"},
			expected:         commitstatus.CommitStatusSuccess,
		},
		{
			commitStatuses: []*git_model.CommitStatus{
				{Context: "Build 1", State: commitstatus.CommitStatusSkipped},
				{Context: "Build 2", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 3", State: commitstatus.CommitStatusSuccess},
			},
			requiredContexts: []string{"Build*"},
			expected:         commitstatus.CommitStatusSuccess,
		},
		{
			commitStatuses: []*git_model.CommitStatus{
				{Context: "Build 1", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2t", State: commitstatus.CommitStatusPending},
			},
			requiredContexts: []string{"Build*", "Build 2t*"},
			expected:         commitstatus.CommitStatusPending,
		},
		{
			commitStatuses: []*git_model.CommitStatus{
				{Context: "Build 1", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2t", State: commitstatus.CommitStatusFailure},
			},
			requiredContexts: []string{"Build*", "Build 2t*"},
			expected:         commitstatus.CommitStatusFailure,
		},
		{
			commitStatuses: []*git_model.CommitStatus{
				{Context: "Build 1", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2t", State: commitstatus.CommitStatusFailure},
			},
			requiredContexts: []string{"Build*"},
			expected:         commitstatus.CommitStatusFailure,
		},
		{
			commitStatuses: []*git_model.CommitStatus{
				{Context: "Build 1", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2t", State: commitstatus.CommitStatusSuccess},
			},
			requiredContexts: []string{"Build*", "Build 2t*", "Build 3*"},
			expected:         commitstatus.CommitStatusPending,
		},
		{
			commitStatuses: []*git_model.CommitStatus{
				{Context: "Build 1", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2", State: commitstatus.CommitStatusSuccess},
				{Context: "Build 2t", State: commitstatus.CommitStatusSuccess},
			},
			requiredContexts: []string{"Build*", "Build *", "Build 2t*", "Build 1*"},
			expected:         commitstatus.CommitStatusSuccess,
		},
	}
	for i, c := range cases {
		assert.Equal(t, c.expected, MergeRequiredContextsCommitStatus(c.commitStatuses, c.requiredContexts), "case %d", i)
	}
}
