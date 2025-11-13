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

package integration

import (
	"net/http"
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIPullCommits(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	pr := unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: 2})
	assert.NoError(t, pr.LoadIssue(t.Context()))
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: pr.HeadRepoID})

	req := NewRequestf(t, http.MethodGet, "/api/v1/repos/%s/%s/pulls/%d/commits", repo.OwnerName, repo.Name, pr.Index)
	resp := MakeRequest(t, req, http.StatusOK)

	var commits []*api.Commit
	DecodeJSON(t, resp, &commits)

	require.Len(t, commits, 2)

	assert.Equal(t, "985f0301dba5e7b34be866819cd15ad3d8f508ee", commits[0].SHA)
	assert.Equal(t, "5c050d3b6d2db231ab1f64e324f1b6b9a0b181c2", commits[1].SHA)

	assert.NotEmpty(t, commits[0].Files)
	assert.NotEmpty(t, commits[1].Files)
	assert.NotNil(t, commits[0].RepoCommit.Verification)
	assert.NotNil(t, commits[1].RepoCommit.Verification)
}

// TODO add tests for already merged PR and closed PR
