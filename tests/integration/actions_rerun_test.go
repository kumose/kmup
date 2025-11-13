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
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	runnerv1 "github.com/kumose/actions-proto-go/runner/v1"
)

func TestActionsRerun(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		session := loginUser(t, user2.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)

		apiRepo := createActionsTestRepo(t, token, "actions-rerun", false)
		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: apiRepo.ID})
		httpContext := NewAPITestContext(t, user2.Name, repo.Name, auth_model.AccessTokenScopeWriteRepository)
		defer doAPIDeleteRepository(httpContext)(t)

		runner := newMockRunner()
		runner.registerAsRepoRunner(t, repo.OwnerName, repo.Name, "mock-runner", []string{"ubuntu-latest"}, false)

		wfTreePath := ".kmup/workflows/actions-rerun-workflow-1.yml"
		wfFileContent := `name: actions-rerun-workflow-1
on: 
  push:
    paths:
      - '.kmup/workflows/actions-rerun-workflow-1.yml'
jobs:
  job1:
    runs-on: ubuntu-latest
    steps:
      - run: echo 'job1'
  job2:
    runs-on: ubuntu-latest
    needs: [job1]
    steps:
      - run: echo 'job2'
`

		opts := getWorkflowCreateFileOptions(user2, repo.DefaultBranch, "create"+wfTreePath, wfFileContent)
		createWorkflowFile(t, token, user2.Name, repo.Name, wfTreePath, opts)

		// fetch and exec job1
		job1Task := runner.fetchTask(t)
		_, _, run := getTaskAndJobAndRunByTaskID(t, job1Task.Id)
		runner.execTask(t, job1Task, &mockTaskOutcome{
			result: runnerv1.Result_RESULT_SUCCESS,
		})
		// RERUN-FAILURE: the run is not done
		req := NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/actions/runs/%d/rerun", user2.Name, repo.Name, run.Index), map[string]string{
			"_csrf": GetUserCSRFToken(t, session),
		})
		session.MakeRequest(t, req, http.StatusBadRequest)
		// fetch and exec job2
		job2Task := runner.fetchTask(t)
		runner.execTask(t, job2Task, &mockTaskOutcome{
			result: runnerv1.Result_RESULT_SUCCESS,
		})

		// RERUN-1: rerun the run
		req = NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/actions/runs/%d/rerun", user2.Name, repo.Name, run.Index), map[string]string{
			"_csrf": GetUserCSRFToken(t, session),
		})
		session.MakeRequest(t, req, http.StatusOK)
		// fetch and exec job1
		job1TaskR1 := runner.fetchTask(t)
		runner.execTask(t, job1TaskR1, &mockTaskOutcome{
			result: runnerv1.Result_RESULT_SUCCESS,
		})
		// fetch and exec job2
		job2TaskR1 := runner.fetchTask(t)
		runner.execTask(t, job2TaskR1, &mockTaskOutcome{
			result: runnerv1.Result_RESULT_SUCCESS,
		})

		// RERUN-2: rerun job1
		req = NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/actions/runs/%d/jobs/%d/rerun", user2.Name, repo.Name, run.Index, 0), map[string]string{
			"_csrf": GetUserCSRFToken(t, session),
		})
		session.MakeRequest(t, req, http.StatusOK)
		// job2 needs job1, so rerunning job1 will also rerun job2
		// fetch and exec job1
		job1TaskR2 := runner.fetchTask(t)
		runner.execTask(t, job1TaskR2, &mockTaskOutcome{
			result: runnerv1.Result_RESULT_SUCCESS,
		})
		// fetch and exec job2
		job2TaskR2 := runner.fetchTask(t)
		runner.execTask(t, job2TaskR2, &mockTaskOutcome{
			result: runnerv1.Result_RESULT_SUCCESS,
		})

		// RERUN-3: rerun job2
		req = NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/actions/runs/%d/jobs/%d/rerun", user2.Name, repo.Name, run.Index, 1), map[string]string{
			"_csrf": GetUserCSRFToken(t, session),
		})
		session.MakeRequest(t, req, http.StatusOK)
		// only job2 will rerun
		// fetch and exec job2
		job2TaskR3 := runner.fetchTask(t)
		runner.execTask(t, job2TaskR3, &mockTaskOutcome{
			result: runnerv1.Result_RESULT_SUCCESS,
		})
		runner.fetchNoTask(t)
	})
}
