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
	"net/http/httptest"
	"testing"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/activitypub"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/routers"
	"github.com/kumose/kmup/tests"

	ap "github.com/go-ap/activitypub"
	"github.com/stretchr/testify/assert"
)

func TestActivityPubPerson(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.Federation.Enabled, true)()
	defer test.MockVariableValue(&testWebRoutes, routers.NormalRoutes())()

	t.Run("ExistingPerson", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		userID := 2
		username := "user2"
		req := NewRequest(t, "GET", fmt.Sprintf("/api/v1/activitypub/user-id/%v", userID))
		resp := MakeRequest(t, req, http.StatusOK)
		body := resp.Body.Bytes()
		assert.Contains(t, string(body), "@context")

		var person ap.Person
		err := person.UnmarshalJSON(body)
		assert.NoError(t, err)

		assert.Equal(t, ap.PersonType, person.Type)
		assert.Equal(t, username, person.PreferredUsername.String())
		keyID := person.GetID().String()
		assert.Regexp(t, fmt.Sprintf("activitypub/user-id/%v$", userID), keyID)
		assert.Regexp(t, fmt.Sprintf("activitypub/user-id/%v/outbox$", userID), person.Outbox.GetID().String())
		assert.Regexp(t, fmt.Sprintf("activitypub/user-id/%v/inbox$", userID), person.Inbox.GetID().String())

		pubKey := person.PublicKey
		assert.NotNil(t, pubKey)
		publicKeyID := keyID + "#main-key"
		assert.Equal(t, pubKey.ID.String(), publicKeyID)

		pubKeyPem := pubKey.PublicKeyPem
		assert.NotNil(t, pubKeyPem)
		assert.Regexp(t, "^-----BEGIN PUBLIC KEY-----", pubKeyPem)
	})
	t.Run("MissingPerson", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		req := NewRequest(t, "GET", "/api/v1/activitypub/user-id/999999999")
		resp := MakeRequest(t, req, http.StatusNotFound)
		assert.Contains(t, resp.Body.String(), "user does not exist")
	})
	t.Run("MissingPersonInbox", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		srv := httptest.NewServer(testWebRoutes)
		defer srv.Close()
		defer test.MockVariableValue(&setting.AppURL, srv.URL+"/")()

		username1 := "user1"
		ctx := t.Context()
		user1, err := user_model.GetUserByName(ctx, username1)
		assert.NoError(t, err)
		user1url := srv.URL + "/api/v1/activitypub/user-id/1#main-key"
		c, err := activitypub.NewClient(t.Context(), user1, user1url)
		assert.NoError(t, err)
		user2inboxurl := srv.URL + "/api/v1/activitypub/user-id/2/inbox"

		// Signed request succeeds
		resp, err := c.Post([]byte{}, user2inboxurl)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Unsigned request fails
		req := NewRequest(t, "POST", user2inboxurl)
		MakeRequest(t, req, http.StatusInternalServerError)
	})
}
