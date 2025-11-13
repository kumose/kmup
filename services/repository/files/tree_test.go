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

package files

import (
	"html/template"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/fileicon"
	"github.com/kumose/kmup/modules/git"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestGetTreeBySHA(t *testing.T) {
	unittest.PrepareTestEnv(t)
	ctx, _ := contexttest.MockContext(t, "user2/repo1")
	contexttest.LoadRepo(t, ctx, 1)
	contexttest.LoadRepoCommit(t, ctx)
	contexttest.LoadUser(t, ctx, 2)
	contexttest.LoadGitRepo(t, ctx)
	defer ctx.Repo.GitRepo.Close()

	sha := ctx.Repo.Repository.DefaultBranch
	page := 1
	perPage := 10
	ctx.SetPathParam("id", "1")
	ctx.SetPathParam("sha", sha)

	tree, err := GetTreeBySHA(ctx, ctx.Repo.Repository, ctx.Repo.GitRepo, ctx.PathParam("sha"), page, perPage, true)
	assert.NoError(t, err)
	expectedTree := &api.GitTreeResponse{
		SHA: "65f1bf27bc3bf70f64657658635e66094edbcb4d",
		URL: "https://try.kmup.io/api/v1/repos/user2/repo1/git/trees/65f1bf27bc3bf70f64657658635e66094edbcb4d",
		Entries: []api.GitEntry{
			{
				Path: "README.md",
				Mode: "100644",
				Type: "blob",
				Size: 30,
				SHA:  "4b4851ad51df6a7d9f25c979345979eaeb5b349f",
				URL:  "https://try.kmup.io/api/v1/repos/user2/repo1/git/blobs/4b4851ad51df6a7d9f25c979345979eaeb5b349f",
			},
		},
		Truncated:  false,
		Page:       1,
		TotalCount: 1,
	}

	assert.Equal(t, expectedTree, tree)
}

func TestGetTreeViewNodes(t *testing.T) {
	unittest.PrepareTestEnv(t)
	ctx, _ := contexttest.MockContext(t, "user2/repo1")
	ctx.Repo.RefFullName = git.RefNameFromBranch("sub-home-md-img-check")
	contexttest.LoadRepo(t, ctx, 1)
	contexttest.LoadRepoCommit(t, ctx)
	contexttest.LoadUser(t, ctx, 2)
	contexttest.LoadGitRepo(t, ctx)
	defer ctx.Repo.GitRepo.Close()

	curRepoLink := "/any/repo-link"
	renderedIconPool := fileicon.NewRenderedIconPool()
	mockIconForFile := func(id string) template.HTML {
		return template.HTML(`<svg class="svg git-entry-icon octicon-file" width="16" height="16" aria-hidden="true"><use xlink:href="#` + id + `"></use></svg>`)
	}
	mockIconForFolder := func(id string) template.HTML {
		return template.HTML(`<svg class="svg git-entry-icon octicon-file-directory-fill" width="16" height="16" aria-hidden="true"><use xlink:href="#` + id + `"></use></svg>`)
	}
	mockOpenIconForFolder := func(id string) template.HTML {
		return template.HTML(`<svg class="svg git-entry-icon octicon-file-directory-open-fill" width="16" height="16" aria-hidden="true"><use xlink:href="#` + id + `"></use></svg>`)
	}
	treeNodes, err := GetTreeViewNodes(ctx, curRepoLink, renderedIconPool, ctx.Repo.Commit, "", "")
	assert.NoError(t, err)
	assert.Equal(t, []*TreeViewNode{
		{
			EntryName:     "docs",
			EntryMode:     "tree",
			FullPath:      "docs",
			EntryIcon:     mockIconForFolder(`svg-mfi-folder-docs`),
			EntryIconOpen: mockOpenIconForFolder(`svg-mfi-folder-docs`),
		},
	}, treeNodes)

	treeNodes, err = GetTreeViewNodes(ctx, curRepoLink, renderedIconPool, ctx.Repo.Commit, "", "docs/README.md")
	assert.NoError(t, err)
	assert.Equal(t, []*TreeViewNode{
		{
			EntryName:     "docs",
			EntryMode:     "tree",
			FullPath:      "docs",
			EntryIcon:     mockIconForFolder(`svg-mfi-folder-docs`),
			EntryIconOpen: mockOpenIconForFolder(`svg-mfi-folder-docs`),
			Children: []*TreeViewNode{
				{
					EntryName: "README.md",
					EntryMode: "blob",
					FullPath:  "docs/README.md",
					EntryIcon: mockIconForFile(`svg-mfi-readme`),
				},
			},
		},
	}, treeNodes)

	treeNodes, err = GetTreeViewNodes(ctx, curRepoLink, renderedIconPool, ctx.Repo.Commit, "docs", "README.md")
	assert.NoError(t, err)
	assert.Equal(t, []*TreeViewNode{
		{
			EntryName: "README.md",
			EntryMode: "blob",
			FullPath:  "docs/README.md",
			EntryIcon: mockIconForFile(`svg-mfi-readme`),
		},
	}, treeNodes)
}
