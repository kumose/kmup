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

package activitypub

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kumose/kmup/modules/activitypub"
	"github.com/kumose/kmup/modules/httplib"
	"github.com/kumose/kmup/modules/setting"
	kmup_context "github.com/kumose/kmup/services/context"

	"github.com/42wim/httpsig"
	ap "github.com/go-ap/activitypub"
)

func getPublicKeyFromResponse(b []byte, keyID *url.URL) (p crypto.PublicKey, err error) {
	person := ap.PersonNew(ap.IRI(keyID.String()))
	err = person.UnmarshalJSON(b)
	if err != nil {
		return nil, fmt.Errorf("ActivityStreams type cannot be converted to one known to have publicKey property: %w", err)
	}
	pubKey := person.PublicKey
	if pubKey.ID.String() != keyID.String() {
		return nil, fmt.Errorf("cannot find publicKey with id: %s in %s", keyID, string(b))
	}
	pubKeyPem := pubKey.PublicKeyPem
	block, _ := pem.Decode([]byte(pubKeyPem))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("could not decode publicKeyPem to PUBLIC KEY pem block type")
	}
	p, err = x509.ParsePKIXPublicKey(block.Bytes)
	return p, err
}

func fetch(iri *url.URL) (b []byte, err error) {
	req := httplib.NewRequest(iri.String(), http.MethodGet)
	req.Header("Accept", activitypub.ActivityStreamsContentType)
	req.Header("User-Agent", "Kmup/"+setting.AppVer)
	resp, err := req.Response()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("url IRI fetch [%s] failed with status (%d): %s", iri, resp.StatusCode, resp.Status)
	}
	b, err = io.ReadAll(io.LimitReader(resp.Body, setting.Federation.MaxSize))
	return b, err
}

func verifyHTTPSignatures(ctx *kmup_context.APIContext) (authenticated bool, err error) {
	r := ctx.Req

	// 1. Figure out what key we need to verify
	v, err := httpsig.NewVerifier(r)
	if err != nil {
		return false, err
	}
	ID := v.KeyId()
	idIRI, err := url.Parse(ID)
	if err != nil {
		return false, err
	}
	// 2. Fetch the public key of the other actor
	b, err := fetch(idIRI)
	if err != nil {
		return false, err
	}
	pubKey, err := getPublicKeyFromResponse(b, idIRI)
	if err != nil {
		return false, err
	}
	// 3. Verify the other actor's key
	algo := httpsig.Algorithm(setting.Federation.Algorithms[0])
	authenticated = v.Verify(pubKey, algo) == nil
	return authenticated, err
}

// ReqHTTPSignature function
func ReqHTTPSignature() func(ctx *kmup_context.APIContext) {
	return func(ctx *kmup_context.APIContext) {
		if authenticated, err := verifyHTTPSignatures(ctx); err != nil {
			ctx.APIErrorInternal(err)
		} else if !authenticated {
			ctx.APIError(http.StatusForbidden, "request signature verification failed")
		}
	}
}
