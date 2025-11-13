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
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"

	"github.com/stretchr/testify/assert"
)

func TestAPIIssueTemplateList(t *testing.T) {
	onKmupRun(t, func(*testing.T, *url.URL) {
		var issueTemplates []*api.IssueTemplate

		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{Name: "user2"})
		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{OwnerName: "user2", Name: "repo1"})

		// no issue template
		req := NewRequest(t, "GET", "/api/v1/repos/user2/repo1/issue_templates")
		resp := MakeRequest(t, req, http.StatusOK)
		issueTemplates = nil
		DecodeJSON(t, resp, &issueTemplates)
		assert.Empty(t, issueTemplates)

		// one correct issue template and some incorrect issue templates
		err := createOrReplaceFileInBranch(user, repo, ".kmup/ISSUE_TEMPLATE/tmpl-ok.md", repo.DefaultBranch, `----
name: foo
about: bar
----
`)
		assert.NoError(t, err)

		err = createOrReplaceFileInBranch(user, repo, ".kmup/ISSUE_TEMPLATE/tmpl-err1.yml", repo.DefaultBranch, `name: '`)
		assert.NoError(t, err)

		err = createOrReplaceFileInBranch(user, repo, ".kmup/ISSUE_TEMPLATE/tmpl-err2.yml", repo.DefaultBranch, `other: `)
		assert.NoError(t, err)

		req = NewRequest(t, "GET", "/api/v1/repos/user2/repo1/issue_templates")
		resp = MakeRequest(t, req, http.StatusOK)
		issueTemplates = nil
		DecodeJSON(t, resp, &issueTemplates)
		assert.Len(t, issueTemplates, 1)
		assert.Equal(t, "foo", issueTemplates[0].Name)
		assert.Equal(t, "error occurs when parsing issue template: count=2", resp.Header().Get("X-Kmup-Warning"))
	})
}
