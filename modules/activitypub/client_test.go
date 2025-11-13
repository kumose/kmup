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
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestActivityPubSignedPost(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	pubID := "https://example.com/pubID"
	c, err := NewClient(t.Context(), user, pubID)
	assert.NoError(t, err)

	expected := "BODY"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Regexp(t, "^"+setting.Federation.DigestAlgorithm, r.Header.Get("Digest"))
		assert.Contains(t, r.Header.Get("Signature"), pubID)
		assert.Equal(t, ActivityStreamsContentType, r.Header.Get("Content-Type"))
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(body))
		fmt.Fprint(w, expected)
	}))
	defer srv.Close()

	r, err := c.Post([]byte(expected), srv.URL)
	assert.NoError(t, err)
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(body))
}
