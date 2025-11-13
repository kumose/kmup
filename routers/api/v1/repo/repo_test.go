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

package repo

import (
	"net/http"
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestRepoEdit(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockAPIContext(t, "user2/repo1")
	contexttest.LoadRepo(t, ctx, 1)
	contexttest.LoadUser(t, ctx, 2)
	ctx.Repo.Owner = ctx.Doer
	description := "new description"
	website := "http://wwww.newwebsite.com"
	private := true
	hasIssues := false
	hasWiki := false
	defaultBranch := "master"
	hasPullRequests := true
	ignoreWhitespaceConflicts := true
	allowMerge := false
	allowRebase := false
	allowRebaseMerge := false
	allowSquashMerge := false
	allowFastForwardOnlyMerge := false
	archived := true
	opts := api.EditRepoOption{
		Name:                      &ctx.Repo.Repository.Name,
		Description:               &description,
		Website:                   &website,
		Private:                   &private,
		HasIssues:                 &hasIssues,
		HasWiki:                   &hasWiki,
		DefaultBranch:             &defaultBranch,
		HasPullRequests:           &hasPullRequests,
		IgnoreWhitespaceConflicts: &ignoreWhitespaceConflicts,
		AllowMerge:                &allowMerge,
		AllowRebase:               &allowRebase,
		AllowRebaseMerge:          &allowRebaseMerge,
		AllowSquash:               &allowSquashMerge,
		AllowFastForwardOnly:      &allowFastForwardOnlyMerge,
		Archived:                  &archived,
	}

	web.SetForm(ctx, &opts)
	Edit(ctx)

	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())
	unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{
		ID: 1,
	}, unittest.Cond("name = ? AND is_archived = 1", *opts.Name))
}

func TestRepoEditNameChange(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockAPIContext(t, "user2/repo1")
	contexttest.LoadRepo(t, ctx, 1)
	contexttest.LoadUser(t, ctx, 2)
	ctx.Repo.Owner = ctx.Doer
	name := "newname"
	opts := api.EditRepoOption{
		Name: &name,
	}

	web.SetForm(ctx, &opts)
	Edit(ctx)
	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())

	unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{
		ID: 1,
	}, unittest.Cond("name = ?", opts.Name))
}
