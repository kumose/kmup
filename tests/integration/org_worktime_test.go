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

package integration_test

import (
	"testing"

	"github.com/kumose/kmup/models/organization"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTimesByRepos tests TimesByRepos functionality
func testTimesByRepos(t *testing.T) {
	kases := []struct {
		name     string
		unixfrom int64
		unixto   int64
		orgname  int64
		expected []organization.WorktimeSumByRepos
	}{
		{
			name:     "Full sum for org 1",
			unixfrom: 0,
			unixto:   9223372036854775807,
			orgname:  1,
			expected: []organization.WorktimeSumByRepos(nil),
		},
		{
			name:     "Full sum for org 2",
			unixfrom: 0,
			unixto:   9223372036854775807,
			orgname:  2,
			expected: []organization.WorktimeSumByRepos{
				{
					RepoName: "repo1",
					SumTime:  4083,
				},
				{
					RepoName: "repo2",
					SumTime:  75,
				},
			},
		},
		{
			name:     "Simple time bound",
			unixfrom: 946684801,
			unixto:   946684802,
			orgname:  2,
			expected: []organization.WorktimeSumByRepos{
				{
					RepoName: "repo1",
					SumTime:  3662,
				},
			},
		},
		{
			name:     "Both times inclusive",
			unixfrom: 946684801,
			unixto:   946684801,
			orgname:  2,
			expected: []organization.WorktimeSumByRepos{
				{
					RepoName: "repo1",
					SumTime:  3661,
				},
			},
		},
		{
			name:     "Should ignore deleted",
			unixfrom: 947688814,
			unixto:   947688815,
			orgname:  2,
			expected: []organization.WorktimeSumByRepos{
				{
					RepoName: "repo2",
					SumTime:  71,
				},
			},
		},
	}

	// Run test kases
	for _, kase := range kases {
		t.Run(kase.name, func(t *testing.T) {
			org, err := organization.GetOrgByID(t.Context(), kase.orgname)
			assert.NoError(t, err)
			results, err := organization.GetWorktimeByRepos(t.Context(), org, kase.unixfrom, kase.unixto)
			assert.NoError(t, err)
			assert.Equal(t, kase.expected, results)
		})
	}
}

// TestTimesByMilestones tests TimesByMilestones functionality
func testTimesByMilestones(t *testing.T) {
	kases := []struct {
		name     string
		unixfrom int64
		unixto   int64
		orgname  int64
		expected []organization.WorktimeSumByMilestones
	}{
		{
			name:     "Full sum for org 1",
			unixfrom: 0,
			unixto:   9223372036854775807,
			orgname:  1,
			expected: []organization.WorktimeSumByMilestones(nil),
		},
		{
			name:     "Full sum for org 2",
			unixfrom: 0,
			unixto:   9223372036854775807,
			orgname:  2,
			expected: []organization.WorktimeSumByMilestones{
				{
					RepoName:      "repo1",
					MilestoneName: "",
					MilestoneID:   0,
					SumTime:       401,
					HideRepoName:  false,
				},
				{
					RepoName:      "repo1",
					MilestoneName: "milestone1",
					MilestoneID:   1,
					SumTime:       3682,
					HideRepoName:  true,
				},
				{
					RepoName:      "repo2",
					MilestoneName: "",
					MilestoneID:   0,
					SumTime:       75,
					HideRepoName:  false,
				},
			},
		},
		{
			name:     "Simple time bound",
			unixfrom: 946684801,
			unixto:   946684802,
			orgname:  2,
			expected: []organization.WorktimeSumByMilestones{
				{
					RepoName:      "repo1",
					MilestoneName: "milestone1",
					MilestoneID:   1,
					SumTime:       3662,
					HideRepoName:  false,
				},
			},
		},
		{
			name:     "Both times inclusive",
			unixfrom: 946684801,
			unixto:   946684801,
			orgname:  2,
			expected: []organization.WorktimeSumByMilestones{
				{
					RepoName:      "repo1",
					MilestoneName: "milestone1",
					MilestoneID:   1,
					SumTime:       3661,
					HideRepoName:  false,
				},
			},
		},
		{
			name:     "Should ignore deleted",
			unixfrom: 947688814,
			unixto:   947688815,
			orgname:  2,
			expected: []organization.WorktimeSumByMilestones{
				{
					RepoName:      "repo2",
					MilestoneName: "",
					MilestoneID:   0,
					SumTime:       71,
					HideRepoName:  false,
				},
			},
		},
	}

	// Run test kases
	for _, kase := range kases {
		t.Run(kase.name, func(t *testing.T) {
			org, err := organization.GetOrgByID(t.Context(), kase.orgname)
			require.NoError(t, err)
			results, err := organization.GetWorktimeByMilestones(t.Context(), org, kase.unixfrom, kase.unixto)
			if assert.NoError(t, err) {
				assert.Equal(t, kase.expected, results)
			}
		})
	}
}

// TestTimesByMembers tests TimesByMembers functionality
func testTimesByMembers(t *testing.T) {
	kases := []struct {
		name     string
		unixfrom int64
		unixto   int64
		orgname  int64
		expected []organization.WorktimeSumByMembers
	}{
		{
			name:     "Full sum for org 1",
			unixfrom: 0,
			unixto:   9223372036854775807,
			orgname:  1,
			expected: []organization.WorktimeSumByMembers(nil),
		},
		{
			// Test case: Sum of times forever in org no. 2
			name:     "Full sum for org 2",
			unixfrom: 0,
			unixto:   9223372036854775807,
			orgname:  2,
			expected: []organization.WorktimeSumByMembers{
				{
					UserName: "user2",
					SumTime:  3666,
				},
				{
					UserName: "user1",
					SumTime:  491,
				},
			},
		},
		{
			name:     "Simple time bound",
			unixfrom: 946684801,
			unixto:   946684802,
			orgname:  2,
			expected: []organization.WorktimeSumByMembers{
				{
					UserName: "user2",
					SumTime:  3662,
				},
			},
		},
		{
			name:     "Both times inclusive",
			unixfrom: 946684801,
			unixto:   946684801,
			orgname:  2,
			expected: []organization.WorktimeSumByMembers{
				{
					UserName: "user2",
					SumTime:  3661,
				},
			},
		},
		{
			name:     "Should ignore deleted",
			unixfrom: 947688814,
			unixto:   947688815,
			orgname:  2,
			expected: []organization.WorktimeSumByMembers{
				{
					UserName: "user1",
					SumTime:  71,
				},
			},
		},
	}

	// Run test kases
	for _, kase := range kases {
		t.Run(kase.name, func(t *testing.T) {
			org, err := organization.GetOrgByID(t.Context(), kase.orgname)
			assert.NoError(t, err)
			results, err := organization.GetWorktimeByMembers(t.Context(), org, kase.unixfrom, kase.unixto)
			assert.NoError(t, err)
			assert.Equal(t, kase.expected, results)
		})
	}
}

func TestOrgWorktime(t *testing.T) {
	// we need to run these tests in integration test because there are complex SQL queries
	assert.NoError(t, unittest.PrepareTestDatabase())
	t.Run("ByRepos", testTimesByRepos)
	t.Run("ByMilestones", testTimesByMilestones)
	t.Run("ByMembers", testTimesByMembers)
}
