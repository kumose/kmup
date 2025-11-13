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

	project_model "github.com/kumose/kmup/models/project"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestPrivateRepoProject(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// not logged in user
	req := NewRequest(t, "GET", "/user31/-/projects")
	MakeRequest(t, req, http.StatusNotFound)

	sess := loginUser(t, "user1")
	req = NewRequest(t, "GET", "/user31/-/projects")
	sess.MakeRequest(t, req, http.StatusOK)
}

func TestMoveRepoProjectColumns(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo2 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})

	projectsUnit := repo2.MustGetUnit(t.Context(), unit.TypeProjects)
	assert.True(t, projectsUnit.ProjectsConfig().IsProjectsAllowed(repo_model.ProjectsModeRepo))

	project1 := project_model.Project{
		Title:        "new created project",
		RepoID:       repo2.ID,
		Type:         project_model.TypeRepository,
		TemplateType: project_model.TemplateTypeNone,
	}
	err := project_model.NewProject(t.Context(), &project1)
	assert.NoError(t, err)

	for i := range 3 {
		err = project_model.NewColumn(t.Context(), &project_model.Column{
			Title:     fmt.Sprintf("column %d", i+1),
			ProjectID: project1.ID,
		})
		assert.NoError(t, err)
	}

	columns, err := project1.GetColumns(t.Context())
	assert.NoError(t, err)
	assert.Len(t, columns, 3)
	assert.EqualValues(t, 0, columns[0].Sorting)
	assert.EqualValues(t, 1, columns[1].Sorting)
	assert.EqualValues(t, 2, columns[2].Sorting)

	sess := loginUser(t, "user1")
	req := NewRequest(t, "GET", fmt.Sprintf("/%s/projects/%d", repo2.FullName(), project1.ID))
	resp := sess.MakeRequest(t, req, http.StatusOK)
	htmlDoc := NewHTMLParser(t, resp.Body)

	req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/%s/projects/%d/move?_csrf="+htmlDoc.GetCSRF(), repo2.FullName(), project1.ID), map[string]any{
		"columns": []map[string]any{
			{"columnID": columns[1].ID, "sorting": 0},
			{"columnID": columns[2].ID, "sorting": 1},
			{"columnID": columns[0].ID, "sorting": 2},
		},
	})
	sess.MakeRequest(t, req, http.StatusOK)

	columnsAfter, err := project1.GetColumns(t.Context())
	assert.NoError(t, err)
	assert.Len(t, columnsAfter, 3)
	assert.Equal(t, columns[1].ID, columnsAfter[0].ID)
	assert.Equal(t, columns[2].ID, columnsAfter[1].ID)
	assert.Equal(t, columns[0].ID, columnsAfter[2].ID)

	assert.NoError(t, project_model.DeleteProjectByID(t.Context(), project1.ID))
}
