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
	"net/url"
	"strings"
	"testing"

	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func testRepoCommitsSearch(t *testing.T, query, commit string) {
	session := loginUser(t, "user2")

	// Request repository commits page
	req := NewRequestf(t, "GET", "/user2/commits_search_test/commits/branch/master/search?q=%s", url.QueryEscape(query))
	resp := session.MakeRequest(t, req, http.StatusOK)

	doc := NewHTMLParser(t, resp.Body)
	sel := doc.doc.Find("#commits-table tbody tr td.sha a")
	assert.Equal(t, commit, strings.TrimSpace(sel.Text()))
}

func TestRepoCommitsSearch(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	testRepoCommitsSearch(t, "e8eabd", "")
	testRepoCommitsSearch(t, "38a9cb", "")
	testRepoCommitsSearch(t, "6e8e", "6e8eabd9a7")
	testRepoCommitsSearch(t, "58e97", "58e97d1a24")
	testRepoCommitsSearch(t, "[build]", "")
	testRepoCommitsSearch(t, "author:alice", "6e8eabd9a7")
	testRepoCommitsSearch(t, "author:alice 6e8ea", "6e8eabd9a7")
	testRepoCommitsSearch(t, "committer:Tom", "58e97d1a24")
	testRepoCommitsSearch(t, "author:bob commit-4", "58e97d1a24")
	testRepoCommitsSearch(t, "author:bob commit after:2019-03-03", "58e97d1a24")
	testRepoCommitsSearch(t, "committer:alice 6e8e before:2019-03-02", "6e8eabd9a7")
	testRepoCommitsSearch(t, "committer:alice commit before:2019-03-02", "6e8eabd9a7")
	testRepoCommitsSearch(t, "committer:alice author:tom commit before:2019-03-04 after:2019-03-02", "0a8499a22a")
}
