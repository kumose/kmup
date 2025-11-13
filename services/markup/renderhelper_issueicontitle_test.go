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

package markup

import (
	"testing"

	"github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestRenderHelperIssueIconTitle(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	ctx, _ := contexttest.MockContext(t, "/", contexttest.MockContextOption{Render: templates.HTMLRenderer()})
	ctx.Repo.Repository = unittest.AssertExistsAndLoadBean(t, &repo.Repository{ID: 1})
	htm, err := renderRepoIssueIconTitle(ctx, markup.RenderIssueIconTitleOptions{
		LinkHref:   "/link",
		IssueIndex: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, `<a href="/link"><span>octicon-issue-opened(16/text green)</span> issue1 (#1)</a>`, string(htm))

	ctx, _ = contexttest.MockContext(t, "/", contexttest.MockContextOption{Render: templates.HTMLRenderer()})
	htm, err = renderRepoIssueIconTitle(ctx, markup.RenderIssueIconTitleOptions{
		OwnerName:  "user2",
		RepoName:   "repo1",
		LinkHref:   "/link",
		IssueIndex: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, `<a href="/link"><span>octicon-issue-opened(16/text green)</span> issue1 (user2/repo1#1)</a>`, string(htm))

	ctx, _ = contexttest.MockContext(t, "/", contexttest.MockContextOption{Render: templates.HTMLRenderer()})
	_, err = renderRepoIssueIconTitle(ctx, markup.RenderIssueIconTitleOptions{
		OwnerName:  "user2",
		RepoName:   "repo2",
		LinkHref:   "/link",
		IssueIndex: 2,
	})
	assert.ErrorIs(t, err, util.ErrPermissionDenied)
}
