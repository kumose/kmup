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

package git

import (
	"testing"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestBranchRuleMatch(t *testing.T) {
	kases := []struct {
		Rule          string
		BranchName    string
		ExpectedMatch bool
	}{
		{
			Rule:          "release/*",
			BranchName:    "release/v1.17",
			ExpectedMatch: true,
		},
		{
			Rule:          "release/**/v1.17",
			BranchName:    "release/test/v1.17",
			ExpectedMatch: true,
		},
		{
			Rule:          "release/**/v1.17",
			BranchName:    "release/test/1/v1.17",
			ExpectedMatch: true,
		},
		{
			Rule:          "release/*/v1.17",
			BranchName:    "release/test/1/v1.17",
			ExpectedMatch: false,
		},
		{
			Rule:          "release/v*",
			BranchName:    "release/v1.16",
			ExpectedMatch: true,
		},
		{
			Rule:          "*",
			BranchName:    "release/v1.16",
			ExpectedMatch: false,
		},
		{
			Rule:          "**",
			BranchName:    "release/v1.16",
			ExpectedMatch: true,
		},
		{
			Rule:          "master",
			BranchName:    "master",
			ExpectedMatch: true,
		},
		{
			Rule:          "master",
			BranchName:    "master",
			ExpectedMatch: false,
		},
	}

	for _, kase := range kases {
		pb := ProtectedBranch{RuleName: kase.Rule}
		var should, infact string
		if !kase.ExpectedMatch {
			should = " not"
		} else {
			infact = " not"
		}
		assert.Equal(t, kase.ExpectedMatch, pb.Match(kase.BranchName),
			"%s should%s match %s but it is%s", kase.BranchName, should, kase.Rule, infact,
		)
	}
}

func TestUpdateProtectBranchPriorities(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	// Create some test protected branches with initial priorities
	protectedBranches := []*ProtectedBranch{
		{
			RepoID:   repo.ID,
			RuleName: "master",
			Priority: 1,
		},
		{
			RepoID:   repo.ID,
			RuleName: "develop",
			Priority: 2,
		},
		{
			RepoID:   repo.ID,
			RuleName: "feature/*",
			Priority: 3,
		},
	}

	for _, pb := range protectedBranches {
		_, err := db.GetEngine(t.Context()).Insert(pb)
		assert.NoError(t, err)
	}

	// Test updating priorities
	newPriorities := []int64{protectedBranches[2].ID, protectedBranches[0].ID, protectedBranches[1].ID}
	err := UpdateProtectBranchPriorities(t.Context(), repo, newPriorities)
	assert.NoError(t, err)

	// Verify new priorities
	pbs, err := FindRepoProtectedBranchRules(t.Context(), repo.ID)
	assert.NoError(t, err)

	expectedPriorities := map[string]int64{
		"feature/*": 1,
		"master":    2,
		"develop":   3,
	}

	for _, pb := range pbs {
		assert.Equal(t, expectedPriorities[pb.RuleName], pb.Priority)
	}
}

func TestNewProtectBranchPriority(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	err := UpdateProtectBranch(t.Context(), repo, &ProtectedBranch{
		RepoID:   repo.ID,
		RuleName: "branch-1",
		Priority: 1,
	}, WhitelistOptions{})
	assert.NoError(t, err)

	newPB := &ProtectedBranch{
		RepoID:   repo.ID,
		RuleName: "branch-2",
		// Priority intentionally omitted
	}

	err = UpdateProtectBranch(t.Context(), repo, newPB, WhitelistOptions{})
	assert.NoError(t, err)

	savedPB2, err := GetFirstMatchProtectedBranchRule(t.Context(), repo.ID, "branch-2")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), savedPB2.Priority)
}
