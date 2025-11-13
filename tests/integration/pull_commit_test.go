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

	pull_service "github.com/kumose/kmup/services/pull"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPullCommits(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user5")
	req := NewRequest(t, "GET", "/user2/repo1/pulls/3/commits/list")
	resp := session.MakeRequest(t, req, http.StatusOK)

	var pullCommitList struct {
		Commits             []pull_service.CommitInfo `json:"commits"`
		LastReviewCommitSha string                    `json:"last_review_commit_sha"`
	}
	DecodeJSON(t, resp, &pullCommitList)

	require.Len(t, pullCommitList.Commits, 2)
	assert.Equal(t, "985f0301dba5e7b34be866819cd15ad3d8f508ee", pullCommitList.Commits[0].ID)
	assert.Equal(t, "5c050d3b6d2db231ab1f64e324f1b6b9a0b181c2", pullCommitList.Commits[1].ID)
	assert.Equal(t, "4a357436d925b5c974181ff12a994538ddc5a269", pullCommitList.LastReviewCommitSha)

	t.Run("CommitBlobExcerpt", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		req = NewRequest(t, "GET", "/user2/repo1/blob_excerpt/985f0301dba5e7b34be866819cd15ad3d8f508ee?last_left=0&last_right=0&left=2&right=2&left_hunk_size=2&right_hunk_size=2&path=README.md&style=split&direction=up")
		resp = session.MakeRequest(t, req, http.StatusOK)
		assert.Contains(t, resp.Body.String(), `<td class="lines-code lines-code-new"><code class="code-inner"># repo1</code>`)
	})
}
