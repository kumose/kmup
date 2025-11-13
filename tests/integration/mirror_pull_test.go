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
	"slices"
	"testing"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/migration"
	mirror_service "github.com/kumose/kmup/services/mirror"
	release_service "github.com/kumose/kmup/services/release"
	repo_service "github.com/kumose/kmup/services/repository"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMirrorPull(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	ctx := t.Context()
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	repoPath := repo_model.RepoPath(user.Name, repo.Name)

	opts := migration.MigrateOptions{
		RepoName:    "test_mirror",
		Description: "Test mirror",
		Private:     false,
		Mirror:      true,
		CloneAddr:   repoPath,
		Wiki:        true,
		Releases:    true,
	}

	mirrorRepo, err := repo_service.CreateRepositoryDirectly(ctx, user, user, repo_service.CreateRepoOptions{
		Name:        opts.RepoName,
		Description: opts.Description,
		IsPrivate:   opts.Private,
		IsMirror:    opts.Mirror,
		Status:      repo_model.RepositoryBeingMigrated,
	}, false)
	assert.NoError(t, err)
	assert.True(t, mirrorRepo.IsMirror, "expected pull-mirror repo to be marked as a mirror immediately after its creation")

	mirrorRepo, err = repo_service.MigrateRepositoryGitData(ctx, user, mirrorRepo, opts, nil)
	assert.NoError(t, err)

	// these units should have been enabled
	mirrorRepo.Units = nil
	require.NoError(t, mirrorRepo.LoadUnits(ctx))
	assert.True(t, slices.ContainsFunc(mirrorRepo.Units, func(u *repo_model.RepoUnit) bool { return u.Type == unit.TypeReleases }))
	assert.True(t, slices.ContainsFunc(mirrorRepo.Units, func(u *repo_model.RepoUnit) bool { return u.Type == unit.TypeWiki }))

	gitRepo, err := gitrepo.OpenRepository(t.Context(), repo)
	assert.NoError(t, err)
	defer gitRepo.Close()

	findOptions := repo_model.FindReleasesOptions{
		IncludeDrafts: true,
		IncludeTags:   true,
		RepoID:        mirrorRepo.ID,
	}
	initCount, err := db.Count[repo_model.Release](t.Context(), findOptions)
	assert.NoError(t, err)
	assert.Zero(t, initCount) // no sync yet, so even though there is a tag in source repo, the mirror's release table is still empty

	assert.NoError(t, release_service.CreateRelease(gitRepo, &repo_model.Release{
		RepoID:       repo.ID,
		Repo:         repo,
		PublisherID:  user.ID,
		Publisher:    user,
		TagName:      "v0.2",
		Target:       "master",
		Title:        "v0.2 is released",
		Note:         "v0.2 is released",
		IsDraft:      false,
		IsPrerelease: false,
		IsTag:        true,
	}, nil, ""))

	_, err = repo_model.GetMirrorByRepoID(ctx, mirrorRepo.ID)
	assert.NoError(t, err)

	ok := mirror_service.SyncPullMirror(ctx, mirrorRepo.ID)
	assert.True(t, ok)

	// actually there is a tag in the source repo, so after "sync", that tag will also come into the mirror
	initCount++

	count, err := db.Count[repo_model.Release](t.Context(), findOptions)
	assert.NoError(t, err)
	assert.Equal(t, initCount+1, count)

	release, err := repo_model.GetRelease(t.Context(), repo.ID, "v0.2")
	assert.NoError(t, err)
	assert.NoError(t, release_service.DeleteReleaseByID(ctx, repo, release, user, true))

	ok = mirror_service.SyncPullMirror(ctx, mirrorRepo.ID)
	assert.True(t, ok)

	count, err = db.Count[repo_model.Release](t.Context(), findOptions)
	assert.NoError(t, err)
	assert.Equal(t, initCount, count)
}
