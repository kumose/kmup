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
	"strconv"
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestChangeDefaultBranch(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	branchesURL := fmt.Sprintf("/%s/%s/settings/branches", owner.Name, repo.Name)

	csrf := GetUserCSRFToken(t, session)
	req := NewRequestWithValues(t, "POST", branchesURL, map[string]string{
		"_csrf":  csrf,
		"action": "default_branch",
		"branch": "DefaultBranch",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	csrf = GetUserCSRFToken(t, session)
	req = NewRequestWithValues(t, "POST", branchesURL, map[string]string{
		"_csrf":  csrf,
		"action": "default_branch",
		"branch": "does_not_exist",
	})
	session.MakeRequest(t, req, http.StatusNotFound)
}

func checkDivergence(t *testing.T, session *TestSession, branchesURL, expectedDefaultBranch string, expectedBranchToDivergence map[string]*gitrepo.DivergeObject) {
	req := NewRequest(t, "GET", branchesURL)
	resp := session.MakeRequest(t, req, http.StatusOK)

	htmlDoc := NewHTMLParser(t, resp.Body)

	branchNodes := htmlDoc.doc.Find(".branch-name").Nodes
	branchNames := []string{}
	for _, node := range branchNodes {
		branchNames = append(branchNames, node.FirstChild.Data)
	}

	expectBranchCount := len(expectedBranchToDivergence)

	assert.Len(t, branchNames, expectBranchCount+1)
	assert.Equal(t, expectedDefaultBranch, branchNames[0])

	allCountBehindNodes := htmlDoc.doc.Find(".count-behind").Nodes
	allCountAheadNodes := htmlDoc.doc.Find(".count-ahead").Nodes

	assert.Len(t, allCountAheadNodes, expectBranchCount)
	assert.Len(t, allCountBehindNodes, expectBranchCount)

	for i := range expectBranchCount {
		branchName := branchNames[i+1]
		assert.Contains(t, expectedBranchToDivergence, branchName)

		expectedCountAhead := expectedBranchToDivergence[branchName].Ahead
		expectedCountBehind := expectedBranchToDivergence[branchName].Behind
		countAhead, err := strconv.Atoi(allCountAheadNodes[i].FirstChild.Data)
		assert.NoError(t, err)
		countBehind, err := strconv.Atoi(allCountBehindNodes[i].FirstChild.Data)
		assert.NoError(t, err)

		assert.Equal(t, expectedCountAhead, countAhead)
		assert.Equal(t, expectedCountBehind, countBehind)
	}
}

func TestChangeDefaultBranchDivergence(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 16})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	branchesURL := fmt.Sprintf("/%s/%s/branches", owner.Name, repo.Name)
	settingsBranchesURL := fmt.Sprintf("/%s/%s/settings/branches", owner.Name, repo.Name)

	// check branch divergence before switching default branch
	expectedBranchToDivergenceBefore := map[string]*gitrepo.DivergeObject{
		"not-signed": {
			Ahead:  0,
			Behind: 0,
		},
		"good-sign-not-yet-validated": {
			Ahead:  0,
			Behind: 1,
		},
		"good-sign": {
			Ahead:  1,
			Behind: 3,
		},
	}
	checkDivergence(t, session, branchesURL, "master", expectedBranchToDivergenceBefore)

	// switch default branch
	newDefaultBranch := "good-sign-not-yet-validated"
	csrf := GetUserCSRFToken(t, session)
	req := NewRequestWithValues(t, "POST", settingsBranchesURL, map[string]string{
		"_csrf":  csrf,
		"action": "default_branch",
		"branch": newDefaultBranch,
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	// check branch divergence after switching default branch
	expectedBranchToDivergenceAfter := map[string]*gitrepo.DivergeObject{
		"master": {
			Ahead:  1,
			Behind: 0,
		},
		"not-signed": {
			Ahead:  1,
			Behind: 0,
		},
		"good-sign": {
			Ahead:  1,
			Behind: 2,
		},
	}
	checkDivergence(t, session, branchesURL, newDefaultBranch, expectedBranchToDivergenceAfter)
}
