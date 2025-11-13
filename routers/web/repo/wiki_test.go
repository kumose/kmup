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
	"io"
	"net/http"
	"net/url"
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/contexttest"
	"github.com/kumose/kmup/services/forms"
	repo_service "github.com/kumose/kmup/services/repository"
	wiki_service "github.com/kumose/kmup/services/wiki"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	content = "Wiki contents for unit tests"
	message = "Wiki commit message for unit tests"
)

func wikiEntry(t *testing.T, repo *repo_model.Repository, wikiName wiki_service.WebPath) *git.TreeEntry {
	wikiRepo, err := gitrepo.OpenRepository(t.Context(), repo.WikiStorageRepo())
	assert.NoError(t, err)
	defer wikiRepo.Close()
	commit, err := wikiRepo.GetBranchCommit("master")
	assert.NoError(t, err)
	entries, err := commit.ListEntries()
	assert.NoError(t, err)
	for _, entry := range entries {
		if entry.Name() == wiki_service.WebPathToGitPath(wikiName) {
			return entry
		}
	}
	return nil
}

func wikiContent(t *testing.T, repo *repo_model.Repository, wikiName wiki_service.WebPath) string {
	entry := wikiEntry(t, repo, wikiName)
	if !assert.NotNil(t, entry) {
		return ""
	}
	reader, err := entry.Blob().DataAsync()
	assert.NoError(t, err)
	defer reader.Close()
	bytes, err := io.ReadAll(reader)
	assert.NoError(t, err)
	return string(bytes)
}

func assertWikiExists(t *testing.T, repo *repo_model.Repository, wikiName wiki_service.WebPath) {
	assert.NotNil(t, wikiEntry(t, repo, wikiName))
}

func assertWikiNotExists(t *testing.T, repo *repo_model.Repository, wikiName wiki_service.WebPath) {
	assert.Nil(t, wikiEntry(t, repo, wikiName))
}

func assertPagesMetas(t *testing.T, expectedNames []string, metas any) {
	pageMetas, ok := metas.([]PageMeta)
	require.True(t, ok)
	require.Len(t, pageMetas, len(expectedNames))

	for i, pageMeta := range pageMetas {
		assert.Equal(t, expectedNames[i], pageMeta.Name)
	}
}

func TestWiki(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki")
	ctx.SetPathParam("*", "Home")
	contexttest.LoadRepo(t, ctx, 1)
	Wiki(ctx)
	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())
	assert.EqualValues(t, "Home", ctx.Data["Title"])
	assertPagesMetas(t, []string{"Home", "Page With Image", "Page With Spaced Name", "Unescaped File"}, ctx.Data["Pages"])

	ctx, _ = contexttest.MockContext(t, "user2/repo1/jpeg.jpg")
	ctx.SetPathParam("*", "jpeg.jpg")
	contexttest.LoadRepo(t, ctx, 1)
	Wiki(ctx)
	assert.Equal(t, http.StatusSeeOther, ctx.Resp.WrittenStatus())
	assert.Equal(t, "/user2/repo1/wiki/raw/jpeg.jpg", ctx.Resp.Header().Get("Location"))
}

func TestWikiPages(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki/?action=_pages")
	contexttest.LoadRepo(t, ctx, 1)
	WikiPages(ctx)
	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())
	assertPagesMetas(t, []string{"Home", "Page With Image", "Page With Spaced Name", "Unescaped File"}, ctx.Data["Pages"])
}

func TestNewWiki(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki/?action=_new")
	contexttest.LoadUser(t, ctx, 2)
	contexttest.LoadRepo(t, ctx, 1)
	NewWiki(ctx)
	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())
	assert.EqualValues(t, ctx.Tr("repo.wiki.new_page"), ctx.Data["Title"])
}

func TestNewWikiPost(t *testing.T) {
	for _, title := range []string{
		"New page",
		"&&&&",
	} {
		unittest.PrepareTestEnv(t)

		ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki/?action=_new")
		contexttest.LoadUser(t, ctx, 2)
		contexttest.LoadRepo(t, ctx, 1)
		web.SetForm(ctx, &forms.NewWikiForm{
			Title:   title,
			Content: content,
			Message: message,
		})
		NewWikiPost(ctx)
		assert.Equal(t, http.StatusSeeOther, ctx.Resp.WrittenStatus())
		assertWikiExists(t, ctx.Repo.Repository, wiki_service.UserTitleToWebPath("", title))
		assert.Equal(t, content, wikiContent(t, ctx.Repo.Repository, wiki_service.UserTitleToWebPath("", title)))
	}
}

func TestNewWikiPost_ReservedName(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki/?action=_new")
	contexttest.LoadUser(t, ctx, 2)
	contexttest.LoadRepo(t, ctx, 1)
	web.SetForm(ctx, &forms.NewWikiForm{
		Title:   "_edit",
		Content: content,
		Message: message,
	})
	NewWikiPost(ctx)
	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())
	assert.EqualValues(t, ctx.Tr("repo.wiki.reserved_page", "_edit"), ctx.Flash.ErrorMsg)
	assertWikiNotExists(t, ctx.Repo.Repository, "_edit")
}

func TestEditWiki(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki/Home?action=_edit")
	ctx.SetPathParam("*", "Home")
	contexttest.LoadUser(t, ctx, 2)
	contexttest.LoadRepo(t, ctx, 1)
	EditWiki(ctx)
	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())
	assert.EqualValues(t, "Home", ctx.Data["Title"])
	assert.Equal(t, wikiContent(t, ctx.Repo.Repository, "Home"), ctx.Data["WikiEditContent"])

	ctx, _ = contexttest.MockContext(t, "user2/repo1/wiki/jpeg.jpg?action=_edit")
	ctx.SetPathParam("*", "jpeg.jpg")
	contexttest.LoadUser(t, ctx, 2)
	contexttest.LoadRepo(t, ctx, 1)
	EditWiki(ctx)
	assert.Equal(t, http.StatusForbidden, ctx.Resp.WrittenStatus())
}

func TestEditWikiPost(t *testing.T) {
	for _, title := range []string{
		"Home",
		"New/<page>",
	} {
		unittest.PrepareTestEnv(t)
		ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki/Home?action=_new")
		ctx.SetPathParam("*", "Home")
		contexttest.LoadUser(t, ctx, 2)
		contexttest.LoadRepo(t, ctx, 1)
		web.SetForm(ctx, &forms.NewWikiForm{
			Title:   title,
			Content: content,
			Message: message,
		})
		EditWikiPost(ctx)
		assert.Equal(t, http.StatusSeeOther, ctx.Resp.WrittenStatus())
		assertWikiExists(t, ctx.Repo.Repository, wiki_service.UserTitleToWebPath("", title))
		assert.Equal(t, content, wikiContent(t, ctx.Repo.Repository, wiki_service.UserTitleToWebPath("", title)))
		if title != "Home" {
			assertWikiNotExists(t, ctx.Repo.Repository, "Home")
		}
	}
}

func TestDeleteWikiPagePost(t *testing.T) {
	unittest.PrepareTestEnv(t)

	ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki/Home?action=_delete")
	contexttest.LoadUser(t, ctx, 2)
	contexttest.LoadRepo(t, ctx, 1)
	DeleteWikiPagePost(ctx)
	assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus())
	assertWikiNotExists(t, ctx.Repo.Repository, "Home")
}

func TestWikiRaw(t *testing.T) {
	for filepath, filetype := range map[string]string{
		"jpeg.jpg":                      "image/jpeg",
		"images/jpeg.jpg":               "image/jpeg",
		"files/Non-Renderable-File.zip": "application/octet-stream",
		"Page With Spaced Name":         "text/plain; charset=utf-8",
		"Page-With-Spaced-Name":         "text/plain; charset=utf-8",
		"Page With Spaced Name.md":      "", // there is no "Page With Spaced Name.md" in repo
		"Page-With-Spaced-Name.md":      "text/plain; charset=utf-8",
	} {
		unittest.PrepareTestEnv(t)

		ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki/raw/"+url.PathEscape(filepath))
		ctx.SetPathParam("*", filepath)
		contexttest.LoadUser(t, ctx, 2)
		contexttest.LoadRepo(t, ctx, 1)
		WikiRaw(ctx)
		if filetype == "" {
			assert.Equal(t, http.StatusNotFound, ctx.Resp.WrittenStatus(), "filepath: %s", filepath)
		} else {
			assert.Equal(t, http.StatusOK, ctx.Resp.WrittenStatus(), "filepath: %s", filepath)
			assert.Equal(t, filetype, ctx.Resp.Header().Get("Content-Type"), "filepath: %s", filepath)
		}
	}
}

func TestDefaultWikiBranch(t *testing.T) {
	unittest.PrepareTestEnv(t)

	// repo with no wiki
	repoWithNoWiki := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	assert.False(t, repo_service.HasWiki(t.Context(), repoWithNoWiki))
	assert.NoError(t, wiki_service.ChangeDefaultWikiBranch(t.Context(), repoWithNoWiki, "master"))

	// repo with wiki
	assert.NoError(t, repo_model.UpdateRepositoryColsNoAutoTime(
		t.Context(),
		&repo_model.Repository{ID: 1, DefaultWikiBranch: "wrong-branch"},
		"default_wiki_branch",
	),
	)

	ctx, _ := contexttest.MockContext(t, "user2/repo1/wiki")
	ctx.SetPathParam("*", "Home")
	contexttest.LoadRepo(t, ctx, 1)
	assert.Equal(t, "wrong-branch", ctx.Repo.Repository.DefaultWikiBranch)
	Wiki(ctx) // after the visiting, the out-of-sync database record will update the branch name to "master"
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.Equal(t, "master", ctx.Repo.Repository.DefaultWikiBranch)

	// invalid branch name should fail
	assert.Error(t, wiki_service.ChangeDefaultWikiBranch(t.Context(), repo, "the bad name"))
	repo = unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.Equal(t, "master", repo.DefaultWikiBranch)

	// the same branch name, should succeed (actually a no-op)
	assert.NoError(t, wiki_service.ChangeDefaultWikiBranch(t.Context(), repo, "master"))
	repo = unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.Equal(t, "master", repo.DefaultWikiBranch)

	// change to another name
	assert.NoError(t, wiki_service.ChangeDefaultWikiBranch(t.Context(), repo, "master"))
	repo = unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.Equal(t, "master", repo.DefaultWikiBranch)
}
