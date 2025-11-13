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
	"time"

	actions_model "github.com/kumose/kmup/models/actions"
	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/routers/web/repo/actions"

	runnerv1 "github.com/kumose/actions-proto-go/runner/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestActionsDeleteRun(t *testing.T) {
	now := time.Now()
	testCase := struct {
		treePath         string
		fileContent      string
		outcomes         map[string]*mockTaskOutcome
		expectedStatuses map[string]string
	}{
		treePath: ".kmup/workflows/test1.yml",
		fileContent: `name: test1
on:
  push:
    paths:
      - .kmup/workflows/test1.yml
jobs:
  job1:
    runs-on: ubuntu-latest
    steps:
      - run: echo job1
  job2:
    runs-on: ubuntu-latest
    steps:
      - run: echo job2
  job3:
    runs-on: ubuntu-latest
    steps:
      - run: echo job3
`,
		outcomes: map[string]*mockTaskOutcome{
			"job1": {
				result: runnerv1.Result_RESULT_SUCCESS,
				logRows: []*runnerv1.LogRow{
					{
						Time:    timestamppb.New(now.Add(4 * time.Second)),
						Content: "  \U0001F433  docker create image",
					},
					{
						Time:    timestamppb.New(now.Add(5 * time.Second)),
						Content: "job1",
					},
					{
						Time:    timestamppb.New(now.Add(6 * time.Second)),
						Content: "\U0001F3C1  Job succeeded",
					},
				},
			},
			"job2": {
				result: runnerv1.Result_RESULT_SUCCESS,
				logRows: []*runnerv1.LogRow{
					{
						Time:    timestamppb.New(now.Add(4 * time.Second)),
						Content: "  \U0001F433  docker create image",
					},
					{
						Time:    timestamppb.New(now.Add(5 * time.Second)),
						Content: "job2",
					},
					{
						Time:    timestamppb.New(now.Add(6 * time.Second)),
						Content: "\U0001F3C1  Job succeeded",
					},
				},
			},
			"job3": {
				result: runnerv1.Result_RESULT_SUCCESS,
				logRows: []*runnerv1.LogRow{
					{
						Time:    timestamppb.New(now.Add(4 * time.Second)),
						Content: "  \U0001F433  docker create image",
					},
					{
						Time:    timestamppb.New(now.Add(5 * time.Second)),
						Content: "job3",
					},
					{
						Time:    timestamppb.New(now.Add(6 * time.Second)),
						Content: "\U0001F3C1  Job succeeded",
					},
				},
			},
		},
		expectedStatuses: map[string]string{
			"job1": actions_model.StatusSuccess.String(),
			"job2": actions_model.StatusSuccess.String(),
			"job3": actions_model.StatusSuccess.String(),
		},
	}
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		session := loginUser(t, user2.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)

		apiRepo := createActionsTestRepo(t, token, "actions-delete-run-test", false)
		runner := newMockRunner()
		runner.registerAsRepoRunner(t, user2.Name, apiRepo.Name, "mock-runner", []string{"ubuntu-latest"}, false)

		opts := getWorkflowCreateFileOptions(user2, apiRepo.DefaultBranch, "create "+testCase.treePath, testCase.fileContent)
		createWorkflowFile(t, token, user2.Name, apiRepo.Name, testCase.treePath, opts)

		runIndex := ""
		for i := 0; i < len(testCase.outcomes); i++ {
			task := runner.fetchTask(t)
			jobName := getTaskJobNameByTaskID(t, token, user2.Name, apiRepo.Name, task.Id)
			outcome := testCase.outcomes[jobName]
			assert.NotNil(t, outcome)
			runner.execTask(t, task, outcome)
			runIndex = task.Context.GetFields()["run_number"].GetStringValue()
			assert.Equal(t, "1", runIndex)
		}

		for i := 0; i < len(testCase.outcomes); i++ {
			req := NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/actions/runs/%s/jobs/%d", user2.Name, apiRepo.Name, runIndex, i), map[string]string{
				"_csrf": GetUserCSRFToken(t, session),
			})
			resp := session.MakeRequest(t, req, http.StatusOK)
			var listResp actions.ViewResponse
			err := json.Unmarshal(resp.Body.Bytes(), &listResp)
			assert.NoError(t, err)
			assert.Len(t, listResp.State.Run.Jobs, 3)

			req = NewRequest(t, "GET", fmt.Sprintf("/%s/%s/actions/runs/%s/jobs/%d/logs", user2.Name, apiRepo.Name, runIndex, i)).
				AddTokenAuth(token)
			MakeRequest(t, req, http.StatusOK)
		}

		req := NewRequestWithValues(t, "GET", fmt.Sprintf("/%s/%s/actions/runs/%s", user2.Name, apiRepo.Name, runIndex), map[string]string{
			"_csrf": GetUserCSRFToken(t, session),
		})
		session.MakeRequest(t, req, http.StatusOK)

		req = NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/actions/runs/%s/delete", user2.Name, apiRepo.Name, runIndex), map[string]string{
			"_csrf": GetUserCSRFToken(t, session),
		})
		session.MakeRequest(t, req, http.StatusOK)

		req = NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/actions/runs/%s/delete", user2.Name, apiRepo.Name, runIndex), map[string]string{
			"_csrf": GetUserCSRFToken(t, session),
		})
		session.MakeRequest(t, req, http.StatusNotFound)

		req = NewRequestWithValues(t, "GET", fmt.Sprintf("/%s/%s/actions/runs/%s", user2.Name, apiRepo.Name, runIndex), map[string]string{
			"_csrf": GetUserCSRFToken(t, session),
		})
		session.MakeRequest(t, req, http.StatusNotFound)

		for i := 0; i < len(testCase.outcomes); i++ {
			req := NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/actions/runs/%s/jobs/%d", user2.Name, apiRepo.Name, runIndex, i), map[string]string{
				"_csrf": GetUserCSRFToken(t, session),
			})
			session.MakeRequest(t, req, http.StatusNotFound)

			req = NewRequest(t, "GET", fmt.Sprintf("/%s/%s/actions/runs/%s/jobs/%d/logs", user2.Name, apiRepo.Name, runIndex, i)).
				AddTokenAuth(token)
			MakeRequest(t, req, http.StatusNotFound)
		}
	})
}
