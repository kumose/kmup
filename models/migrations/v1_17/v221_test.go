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

package v1_17

import (
	"encoding/base32"
	"testing"

	"github.com/kumose/kmup/models/migrations/base"

	"github.com/stretchr/testify/assert"
)

func Test_StoreWebauthnCredentialIDAsBytes(t *testing.T) {
	// Create webauthnCredential table
	type WebauthnCredential struct {
		ID              int64 `xorm:"pk autoincr"`
		Name            string
		LowerName       string `xorm:"unique(s)"`
		UserID          int64  `xorm:"INDEX unique(s)"`
		CredentialID    string `xorm:"INDEX VARCHAR(410)"`
		PublicKey       []byte
		AttestationType string
		AAGUID          []byte
		SignCount       uint32 `xorm:"BIGINT"`
		CloneWarning    bool
	}

	type ExpectedWebauthnCredential struct {
		ID           int64  `xorm:"pk autoincr"`
		CredentialID string // CredentialID is at most 1023 bytes as per spec released 20 July 2022
	}

	type ConvertedWebauthnCredential struct {
		ID                int64  `xorm:"pk autoincr"`
		CredentialIDBytes []byte `xorm:"VARBINARY(1024)"` // CredentialID is at most 1023 bytes as per spec released 20 July 2022
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(WebauthnCredential), new(ExpectedWebauthnCredential))
	defer deferable()
	if x == nil || t.Failed() {
		return
	}

	if err := StoreWebauthnCredentialIDAsBytes(x); err != nil {
		assert.NoError(t, err)
		return
	}

	expected := []ExpectedWebauthnCredential{}
	if err := x.Table("expected_webauthn_credential").Asc("id").Find(&expected); !assert.NoError(t, err) {
		return
	}

	got := []ConvertedWebauthnCredential{}
	if err := x.Table("webauthn_credential").Select("id, credential_id_bytes").Asc("id").Find(&got); !assert.NoError(t, err) {
		return
	}

	for i, e := range expected {
		credIDBytes, _ := base32.HexEncoding.DecodeString(e.CredentialID)
		assert.Equal(t, credIDBytes, got[i].CredentialIDBytes)
	}
}
