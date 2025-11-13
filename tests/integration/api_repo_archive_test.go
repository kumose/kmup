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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/perm"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIDownloadArchive(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	session := loginUser(t, user2.LowerName)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)

	link, _ := url.Parse(fmt.Sprintf("/api/v1/repos/%s/%s/archive/master.zip", user2.Name, repo.Name))
	resp := MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
	bs, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, bs, 320)

	link, _ = url.Parse(fmt.Sprintf("/api/v1/repos/%s/%s/archive/master.tar.gz", user2.Name, repo.Name))
	resp = MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
	bs, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, bs, 266)

	// Must return a link to a commit ID as the "immutable" archive link
	linkHeaderRe := regexp.MustCompile(`^<(https?://.*/api/v1/repos/user2/repo1/archive/[a-f0-9]+\.tar\.gz.*)>; rel="immutable"$`)
	m := linkHeaderRe.FindStringSubmatch(resp.Header().Get("Link"))
	assert.NotEmpty(t, m[1])
	resp = MakeRequest(t, NewRequest(t, "GET", m[1]).AddTokenAuth(token), http.StatusOK)
	bs2, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	// The locked URL should give the same bytes as the non-locked one
	assert.Equal(t, bs, bs2)

	link, _ = url.Parse(fmt.Sprintf("/api/v1/repos/%s/%s/archive/master.bundle", user2.Name, repo.Name))
	resp = MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
	bs, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, bs, 382)

	link, _ = url.Parse(fmt.Sprintf("/api/v1/repos/%s/%s/archive/master", user2.Name, repo.Name))
	MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusBadRequest)

	t.Run("GitHubStyle", testAPIDownloadArchiveGitHubStyle)
	t.Run("PrivateRepo", testAPIDownloadArchivePrivateRepo)
}

func testAPIDownloadArchiveGitHubStyle(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	session := loginUser(t, user2.LowerName)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)

	link, _ := url.Parse(fmt.Sprintf("/api/v1/repos/%s/%s/zipball/master", user2.Name, repo.Name))
	resp := MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
	bs, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, bs, 320)

	link, _ = url.Parse(fmt.Sprintf("/api/v1/repos/%s/%s/tarball/master", user2.Name, repo.Name))
	resp = MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
	bs, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, bs, 266)

	// Must return a link to a commit ID as the "immutable" archive link
	linkHeaderRe := regexp.MustCompile(`^<(https?://.*/api/v1/repos/user2/repo1/archive/[a-f0-9]+\.tar\.gz.*)>; rel="immutable"$`)
	m := linkHeaderRe.FindStringSubmatch(resp.Header().Get("Link"))
	assert.NotEmpty(t, m[1])
	resp = MakeRequest(t, NewRequest(t, "GET", m[1]).AddTokenAuth(token), http.StatusOK)
	bs2, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	// The locked URL should give the same bytes as the non-locked one
	assert.Equal(t, bs, bs2)

	link, _ = url.Parse(fmt.Sprintf("/api/v1/repos/%s/%s/bundle/master", user2.Name, repo.Name))
	resp = MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
	bs, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, bs, 382)
}

func testAPIDownloadArchivePrivateRepo(t *testing.T) {
	_ = repo_model.UpdateRepositoryColsNoAutoTime(t.Context(), &repo_model.Repository{ID: 1, IsPrivate: true}, "is_private")
	MakeRequest(t, NewRequest(t, "HEAD", "/api/v1/repos/user2/repo1/archive/master.zip"), http.StatusNotFound)
	MakeRequest(t, NewRequest(t, "HEAD", "/api/v1/repos/user2/repo1/zipball/master"), http.StatusNotFound)
	_ = repo_model.UpdateRepoUnitPublicAccess(t.Context(), &repo_model.RepoUnit{RepoID: 1, Type: unit.TypeCode, AnonymousAccessMode: perm.AccessModeRead})
	MakeRequest(t, NewRequest(t, "HEAD", "/api/v1/repos/user2/repo1/archive/master.zip"), http.StatusOK)
	MakeRequest(t, NewRequest(t, "HEAD", "/api/v1/repos/user2/repo1/zipball/master"), http.StatusOK)
}
