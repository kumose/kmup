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
	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestBlockUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	countStars := func(t *testing.T, repoOwnerID, starrerID int64) int64 {
		count, err := db.Count[repo_model.Repository](t.Context(), &repo_model.StarredReposOptions{
			StarrerID:      starrerID,
			RepoOwnerID:    repoOwnerID,
			IncludePrivate: true,
		})
		assert.NoError(t, err)
		return count
	}

	countWatches := func(t *testing.T, repoOwnerID, watcherID int64) int64 {
		count, err := db.Count[repo_model.Repository](t.Context(), &repo_model.WatchedReposOptions{
			WatcherID:   watcherID,
			RepoOwnerID: repoOwnerID,
		})
		assert.NoError(t, err)
		return count
	}

	countRepositoryTransfers := func(t *testing.T, senderID, recipientID int64) int64 {
		transfers, err := repo_model.GetPendingRepositoryTransfers(t.Context(), &repo_model.PendingRepositoryTransferOptions{
			SenderID:    senderID,
			RecipientID: recipientID,
		})
		assert.NoError(t, err)
		return int64(len(transfers))
	}

	countAssignedIssues := func(t *testing.T, repoOwnerID, assigneeID int64) int64 {
		_, count, err := issues_model.GetAssignedIssues(t.Context(), &issues_model.AssignedIssuesOptions{
			AssigneeID:  assigneeID,
			RepoOwnerID: repoOwnerID,
		})
		assert.NoError(t, err)
		return count
	}

	countCollaborations := func(t *testing.T, repoOwnerID, collaboratorID int64) int64 {
		count, err := db.Count[repo_model.Collaboration](t.Context(), &repo_model.FindCollaborationOptions{
			CollaboratorID: collaboratorID,
			RepoOwnerID:    repoOwnerID,
		})
		assert.NoError(t, err)
		return count
	}

	t.Run("User", func(t *testing.T) {
		var blockerID int64 = 16
		blockerName := "user16"
		blockerToken := getUserToken(t, blockerName, auth_model.AccessTokenScopeWriteUser)

		var blockeeID int64 = 10
		blockeeName := "user10"

		t.Run("Block", func(t *testing.T) {
			req := NewRequest(t, "PUT", "/api/v1/user/blocks/"+blockeeName)
			MakeRequest(t, req, http.StatusUnauthorized)

			assert.EqualValues(t, 1, countStars(t, blockerID, blockeeID))
			assert.EqualValues(t, 1, countWatches(t, blockerID, blockeeID))
			assert.EqualValues(t, 1, countRepositoryTransfers(t, blockerID, blockeeID))
			assert.EqualValues(t, 1, countCollaborations(t, blockerID, blockeeID))

			req = NewRequest(t, "GET", "/api/v1/user/blocks/"+blockeeName).
				AddTokenAuth(blockerToken)
			MakeRequest(t, req, http.StatusNotFound)

			req = NewRequest(t, "PUT", fmt.Sprintf("/api/v1/user/blocks/%s?reason=test", blockeeName)).
				AddTokenAuth(blockerToken)
			MakeRequest(t, req, http.StatusNoContent)

			assert.EqualValues(t, 0, countStars(t, blockerID, blockeeID))
			assert.EqualValues(t, 0, countWatches(t, blockerID, blockeeID))
			assert.EqualValues(t, 0, countRepositoryTransfers(t, blockerID, blockeeID))
			assert.EqualValues(t, 0, countCollaborations(t, blockerID, blockeeID))

			req = NewRequest(t, "GET", "/api/v1/user/blocks/"+blockeeName).
				AddTokenAuth(blockerToken)
			MakeRequest(t, req, http.StatusNoContent)

			req = NewRequest(t, "PUT", "/api/v1/user/blocks/"+blockeeName).
				AddTokenAuth(blockerToken)
			MakeRequest(t, req, http.StatusBadRequest) // can't block blocked user

			req = NewRequest(t, "PUT", "/api/v1/user/blocks/"+"org3").
				AddTokenAuth(blockerToken)
			MakeRequest(t, req, http.StatusBadRequest) // can't block organization

			req = NewRequest(t, "GET", "/api/v1/user/blocks")
			MakeRequest(t, req, http.StatusUnauthorized)

			req = NewRequest(t, "GET", "/api/v1/user/blocks").
				AddTokenAuth(blockerToken)
			resp := MakeRequest(t, req, http.StatusOK)

			var users []api.User
			DecodeJSON(t, resp, &users)

			assert.Len(t, users, 1)
			assert.Equal(t, blockeeName, users[0].UserName)
		})

		t.Run("Unblock", func(t *testing.T) {
			req := NewRequest(t, "DELETE", "/api/v1/user/blocks/"+blockeeName)
			MakeRequest(t, req, http.StatusUnauthorized)

			req = NewRequest(t, "DELETE", "/api/v1/user/blocks/"+blockeeName).
				AddTokenAuth(blockerToken)
			MakeRequest(t, req, http.StatusNoContent)

			req = NewRequest(t, "DELETE", "/api/v1/user/blocks/"+blockeeName).
				AddTokenAuth(blockerToken)
			MakeRequest(t, req, http.StatusBadRequest)

			req = NewRequest(t, "DELETE", "/api/v1/user/blocks/"+"org3").
				AddTokenAuth(blockerToken)
			MakeRequest(t, req, http.StatusBadRequest)

			req = NewRequest(t, "GET", "/api/v1/user/blocks").
				AddTokenAuth(blockerToken)
			resp := MakeRequest(t, req, http.StatusOK)

			var users []api.User
			DecodeJSON(t, resp, &users)

			assert.Empty(t, users)
		})
	})

	t.Run("Organization", func(t *testing.T) {
		var blockerID int64 = 3
		blockerName := "org3"

		doerToken := getUserToken(t, "user2", auth_model.AccessTokenScopeWriteUser, auth_model.AccessTokenScopeWriteOrganization)

		var blockeeID int64 = 10
		blockeeName := "user10"

		t.Run("Block", func(t *testing.T) {
			req := NewRequest(t, "PUT", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, blockeeName))
			MakeRequest(t, req, http.StatusUnauthorized)

			req = NewRequest(t, "PUT", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, "user4")).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusBadRequest) // can't block member

			assert.EqualValues(t, 1, countStars(t, blockerID, blockeeID))
			assert.EqualValues(t, 1, countWatches(t, blockerID, blockeeID))
			assert.EqualValues(t, 1, countRepositoryTransfers(t, blockerID, blockeeID))
			assert.EqualValues(t, 1, countAssignedIssues(t, blockerID, blockeeID))
			assert.EqualValues(t, 1, countCollaborations(t, blockerID, blockeeID))

			req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, blockeeName)).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusNotFound)

			req = NewRequest(t, "PUT", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s?reason=test", blockerName, blockeeName)).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusNoContent)

			assert.EqualValues(t, 0, countStars(t, blockerID, blockeeID))
			assert.EqualValues(t, 0, countWatches(t, blockerID, blockeeID))
			assert.EqualValues(t, 0, countRepositoryTransfers(t, blockerID, blockeeID))
			assert.EqualValues(t, 0, countAssignedIssues(t, blockerID, blockeeID))
			assert.EqualValues(t, 0, countCollaborations(t, blockerID, blockeeID))

			req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, blockeeName)).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusNoContent)

			req = NewRequest(t, "PUT", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, blockeeName)).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusBadRequest) // can't block blocked user

			req = NewRequest(t, "PUT", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, "org3")).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusBadRequest) // can't block organization

			req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/orgs/%s/blocks", blockerName))
			MakeRequest(t, req, http.StatusUnauthorized)

			req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/orgs/%s/blocks", blockerName)).
				AddTokenAuth(doerToken)
			resp := MakeRequest(t, req, http.StatusOK)

			var users []api.User
			DecodeJSON(t, resp, &users)

			assert.Len(t, users, 1)
			assert.Equal(t, blockeeName, users[0].UserName)
		})

		t.Run("Unblock", func(t *testing.T) {
			req := NewRequest(t, "DELETE", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, blockeeName))
			MakeRequest(t, req, http.StatusUnauthorized)

			req = NewRequest(t, "DELETE", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, blockeeName)).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusNoContent)

			req = NewRequest(t, "DELETE", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, blockeeName)).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusBadRequest)

			req = NewRequest(t, "DELETE", fmt.Sprintf("/api/v1/orgs/%s/blocks/%s", blockerName, "org3")).
				AddTokenAuth(doerToken)
			MakeRequest(t, req, http.StatusBadRequest)

			req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/orgs/%s/blocks", blockerName)).
				AddTokenAuth(doerToken)
			resp := MakeRequest(t, req, http.StatusOK)

			var users []api.User
			DecodeJSON(t, resp, &users)

			assert.Empty(t, users)
		})
	})
}
