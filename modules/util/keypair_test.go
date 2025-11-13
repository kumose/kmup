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

package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeygen(t *testing.T) {
	priv, pub, err := GenerateKeyPair(2048)
	assert.NoError(t, err)

	assert.NotEmpty(t, priv)
	assert.NotEmpty(t, pub)

	assert.Regexp(t, "^-----BEGIN RSA PRIVATE KEY-----.*", priv)
	assert.Regexp(t, "^-----BEGIN PUBLIC KEY-----.*", pub)
}

func TestSignUsingKeys(t *testing.T) {
	priv, pub, err := GenerateKeyPair(2048)
	assert.NoError(t, err)

	privPem, _ := pem.Decode([]byte(priv))
	if privPem == nil || privPem.Type != "RSA PRIVATE KEY" {
		t.Fatal("key is wrong type")
	}

	privParsed, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	assert.NoError(t, err)

	pubPem, _ := pem.Decode([]byte(pub))
	if pubPem == nil || pubPem.Type != "PUBLIC KEY" {
		t.Fatal("key failed to decode")
	}

	pubParsed, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	assert.NoError(t, err)

	// Sign
	msg := "activity pub is great!"
	h := sha256.New()
	h.Write([]byte(msg))
	d := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privParsed, crypto.SHA256, d)
	assert.NoError(t, err)

	// Verify
	err = rsa.VerifyPKCS1v15(pubParsed.(*rsa.PublicKey), crypto.SHA256, d, sig)
	assert.NoError(t, err)
}
