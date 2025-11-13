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
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/services/migrations"

	"github.com/stretchr/testify/assert"
)

func TestMigrateLocalPath(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	adminUser := unittest.AssertExistsAndLoadBean(t, &user_model.User{Name: "user1"})

	old := setting.ImportLocalPaths
	setting.ImportLocalPaths = true

	basePath := t.TempDir()

	lowercasePath := filepath.Join(basePath, "lowercase")
	err := os.Mkdir(lowercasePath, 0o700)
	assert.NoError(t, err)

	err = migrations.IsMigrateURLAllowed(lowercasePath, adminUser)
	assert.NoError(t, err, "case lowercase path")

	mixedcasePath := filepath.Join(basePath, "mIxeDCaSe")
	err = os.Mkdir(mixedcasePath, 0o700)
	assert.NoError(t, err)

	err = migrations.IsMigrateURLAllowed(mixedcasePath, adminUser)
	assert.NoError(t, err, "case mixedcase path")

	setting.ImportLocalPaths = old
}

func TestMigrateKmupForm(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		AllowLocalNetworks := setting.Migrations.AllowLocalNetworks
		setting.Migrations.AllowLocalNetworks = true
		AppVer := setting.AppVer
		// Kmup SDK (go-sdk) need to parse the AppVer from server response, so we must set it to a valid version string.
		setting.AppVer = "1.16.0"
		defer func() {
			setting.Migrations.AllowLocalNetworks = AllowLocalNetworks
			setting.AppVer = AppVer
			migrations.Init()
		}()
		assert.NoError(t, migrations.Init())

		ownerName := "user2"
		repoName := "repo1"
		repoOwner := unittest.AssertExistsAndLoadBean(t, &user_model.User{Name: ownerName})
		session := loginUser(t, ownerName)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeReadMisc)

		// Step 0: verify the repo is available
		req := NewRequestf(t, "GET", "/%s/%s", ownerName, repoName)
		_ = session.MakeRequest(t, req, http.StatusOK)
		// Step 1: get the Kmup migration form
		req = NewRequestf(t, "GET", "/repo/migrate/?service_type=%d", structs.KmupService)
		resp := session.MakeRequest(t, req, http.StatusOK)
		// Step 2: load the form
		htmlDoc := NewHTMLParser(t, resp.Body)
		form := htmlDoc.doc.Find(`form.ui.form[action^="/repo/migrate"]`)
		link, exists := form.Attr("action")
		assert.True(t, exists, "The template has changed")
		serviceInput, exists := form.Find(`input[name="service"]`).Attr("value")
		assert.True(t, exists)
		assert.Equal(t, fmt.Sprintf("%d", structs.KmupService), serviceInput)
		// Step 4: submit the migration to only migrate issues
		migratedRepoName := "otherrepo"
		req = NewRequestWithValues(t, "POST", link, map[string]string{
			"_csrf":       htmlDoc.GetCSRF(),
			"service":     fmt.Sprintf("%d", structs.KmupService),
			"clone_addr":  fmt.Sprintf("%s%s/%s", u, ownerName, repoName),
			"auth_token":  token,
			"issues":      "on",
			"repo_name":   migratedRepoName,
			"description": "",
			"uid":         strconv.FormatInt(repoOwner.ID, 10),
		})
		resp = session.MakeRequest(t, req, http.StatusSeeOther)
		// Step 5: a redirection displays the migrated repository
		loc := resp.Header().Get("Location")
		assert.Equal(t, fmt.Sprintf("/%s/%s", ownerName, migratedRepoName), loc)
		// Step 6: check the repo was created
		unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{Name: migratedRepoName})
	})
}

func Test_UpdateCommentsMigrationsByType(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	err := issues_model.UpdateCommentsMigrationsByType(t.Context(), structs.GithubService, "1", 1)
	assert.NoError(t, err)
}
