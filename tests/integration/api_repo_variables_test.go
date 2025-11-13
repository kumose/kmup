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
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"
)

func TestAPIRepoVariables(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

	t.Run("CreateRepoVariable", func(t *testing.T) {
		cases := []struct {
			Name           string
			ExpectedStatus int
		}{
			{
				Name:           "-",
				ExpectedStatus: http.StatusBadRequest,
			},
			{
				Name:           "_",
				ExpectedStatus: http.StatusCreated,
			},
			{
				Name:           "TEST_VAR",
				ExpectedStatus: http.StatusCreated,
			},
			{
				Name:           "test_var",
				ExpectedStatus: http.StatusConflict,
			},
			{
				Name:           "ci",
				ExpectedStatus: http.StatusBadRequest,
			},
			{
				Name:           "123var",
				ExpectedStatus: http.StatusBadRequest,
			},
			{
				Name:           "var@test",
				ExpectedStatus: http.StatusBadRequest,
			},
			{
				Name:           "github_var",
				ExpectedStatus: http.StatusBadRequest,
			},
			{
				Name:           "kmup_var",
				ExpectedStatus: http.StatusBadRequest,
			},
		}

		for _, c := range cases {
			req := NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/actions/variables/%s", repo.FullName(), c.Name), api.CreateVariableOption{
				Value: "value",
			}).AddTokenAuth(token)
			MakeRequest(t, req, c.ExpectedStatus)
		}
	})

	t.Run("UpdateRepoVariable", func(t *testing.T) {
		variableName := "test_update_var"
		url := fmt.Sprintf("/api/v1/repos/%s/actions/variables/%s", repo.FullName(), variableName)
		req := NewRequestWithJSON(t, "POST", url, api.CreateVariableOption{
			Value: "initial_val",
		}).AddTokenAuth(token)
		MakeRequest(t, req, http.StatusCreated)

		cases := []struct {
			Name           string
			UpdateName     string
			ExpectedStatus int
		}{
			{
				Name:           "not_found_var",
				ExpectedStatus: http.StatusNotFound,
			},
			{
				Name:           variableName,
				UpdateName:     "1invalid",
				ExpectedStatus: http.StatusBadRequest,
			},
			{
				Name:           variableName,
				UpdateName:     "invalid@name",
				ExpectedStatus: http.StatusBadRequest,
			},
			{
				Name:           variableName,
				UpdateName:     "ci",
				ExpectedStatus: http.StatusBadRequest,
			},
			{
				Name:           variableName,
				UpdateName:     "updated_var_name",
				ExpectedStatus: http.StatusNoContent,
			},
			{
				Name:           variableName,
				ExpectedStatus: http.StatusNotFound,
			},
			{
				Name:           "updated_var_name",
				ExpectedStatus: http.StatusNoContent,
			},
		}

		for _, c := range cases {
			req := NewRequestWithJSON(t, "PUT", fmt.Sprintf("/api/v1/repos/%s/actions/variables/%s", repo.FullName(), c.Name), api.UpdateVariableOption{
				Name:  c.UpdateName,
				Value: "updated_val",
			}).AddTokenAuth(token)
			MakeRequest(t, req, c.ExpectedStatus)
		}
	})

	t.Run("DeleteRepoVariable", func(t *testing.T) {
		variableName := "test_delete_var"
		url := fmt.Sprintf("/api/v1/repos/%s/actions/variables/%s", repo.FullName(), variableName)

		req := NewRequestWithJSON(t, "POST", url, api.CreateVariableOption{
			Value: "initial_val",
		}).AddTokenAuth(token)
		MakeRequest(t, req, http.StatusCreated)

		req = NewRequest(t, "DELETE", url).AddTokenAuth(token)
		MakeRequest(t, req, http.StatusNoContent)

		req = NewRequest(t, "DELETE", url).AddTokenAuth(token)
		MakeRequest(t, req, http.StatusNotFound)
	})
}
