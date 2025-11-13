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
	"net/http"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIListEmails(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadUser)

	req := NewRequest(t, "GET", "/api/v1/user/emails").
		AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusOK)

	var emails []*api.Email
	DecodeJSON(t, resp, &emails)

	assert.Equal(t, []*api.Email{
		{
			Email:    "user2@example.com",
			Verified: true,
			Primary:  true,
		},
		{
			Email:    "user2-2@example.com",
			Verified: false,
			Primary:  false,
		},
	}, emails)
}

func TestAPIAddEmail(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser)

	opts := api.CreateEmailOption{
		Emails: []string{"user101@example.com"},
	}

	req := NewRequestWithJSON(t, "POST", "/api/v1/user/emails", &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusUnprocessableEntity)

	opts = api.CreateEmailOption{
		Emails: []string{"user2-3@example.com"},
	}
	req = NewRequestWithJSON(t, "POST", "/api/v1/user/emails", &opts).
		AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusCreated)

	var emails []*api.Email
	DecodeJSON(t, resp, &emails)
	assert.Equal(t, []*api.Email{
		{
			Email:    "user2@example.com",
			Verified: true,
			Primary:  true,
		},
		{
			Email:    "user2-2@example.com",
			Verified: false,
			Primary:  false,
		},
		{
			Email:    "user2-3@example.com",
			Verified: true,
			Primary:  false,
		},
	}, emails)

	opts = api.CreateEmailOption{
		Emails: []string{"notAEmail"},
	}
	req = NewRequestWithJSON(t, "POST", "/api/v1/user/emails", &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusUnprocessableEntity)
}

func TestAPIDeleteEmail(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser)

	opts := api.DeleteEmailOption{
		Emails: []string{"user2-3@example.com"},
	}
	req := NewRequestWithJSON(t, "DELETE", "/api/v1/user/emails", &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNotFound)

	opts = api.DeleteEmailOption{
		Emails: []string{"user2-2@example.com"},
	}
	req = NewRequestWithJSON(t, "DELETE", "/api/v1/user/emails", &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	req = NewRequest(t, "GET", "/api/v1/user/emails").
		AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusOK)

	var emails []*api.Email
	DecodeJSON(t, resp, &emails)
	assert.Equal(t, []*api.Email{
		{
			Email:    "user2@example.com",
			Verified: true,
			Primary:  true,
		},
	}, emails)
}
