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

package auth_test

import (
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
)

func TestGetWebAuthnCredentialByID(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	res, err := auth_model.GetWebAuthnCredentialByID(t.Context(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "WebAuthn credential", res.Name)

	_, err = auth_model.GetWebAuthnCredentialByID(t.Context(), 342432)
	assert.Error(t, err)
	assert.True(t, auth_model.IsErrWebAuthnCredentialNotExist(err))
}

func TestGetWebAuthnCredentialsByUID(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	res, err := auth_model.GetWebAuthnCredentialsByUID(t.Context(), 32)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "WebAuthn credential", res[0].Name)
}

func TestWebAuthnCredential_TableName(t *testing.T) {
	assert.Equal(t, "webauthn_credential", auth_model.WebAuthnCredential{}.TableName())
}

func TestWebAuthnCredential_UpdateSignCount(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	cred := unittest.AssertExistsAndLoadBean(t, &auth_model.WebAuthnCredential{ID: 1})
	cred.SignCount = 1
	assert.NoError(t, cred.UpdateSignCount(t.Context()))
	unittest.AssertExistsAndLoadBean(t, &auth_model.WebAuthnCredential{ID: 1, SignCount: 1})
}

func TestWebAuthnCredential_UpdateLargeCounter(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	cred := unittest.AssertExistsAndLoadBean(t, &auth_model.WebAuthnCredential{ID: 1})
	cred.SignCount = 0xffffffff
	assert.NoError(t, cred.UpdateSignCount(t.Context()))
	unittest.AssertExistsAndLoadBean(t, &auth_model.WebAuthnCredential{ID: 1, SignCount: 0xffffffff})
}

func TestCreateCredential(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	res, err := auth_model.CreateCredential(t.Context(), 1, "WebAuthn Created Credential", &webauthn.Credential{ID: []byte("Test")})
	assert.NoError(t, err)
	assert.Equal(t, "WebAuthn Created Credential", res.Name)
	assert.Equal(t, []byte("Test"), res.CredentialID)

	unittest.AssertExistsAndLoadBean(t, &auth_model.WebAuthnCredential{Name: "WebAuthn Created Credential", UserID: 1})
}
