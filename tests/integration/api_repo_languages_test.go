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
	"net/url"
	"testing"
	"time"

	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestRepoLanguages(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		session := loginUser(t, "user2")

		// Request editor page
		req := NewRequest(t, "GET", "/user2/repo1/_new/master/")
		resp := session.MakeRequest(t, req, http.StatusOK)

		doc := NewHTMLParser(t, resp.Body)
		lastCommit := doc.GetInputValueByName("last_commit")
		assert.NotEmpty(t, lastCommit)

		// Save new file to master branch
		req = NewRequestWithValues(t, "POST", "/user2/repo1/_new/master/", map[string]string{
			"_csrf":         doc.GetCSRF(),
			"last_commit":   lastCommit,
			"tree_path":     "test.go",
			"content":       "package main",
			"commit_choice": "direct",
		})
		resp = session.MakeRequest(t, req, http.StatusOK)
		assert.NotEmpty(t, test.RedirectURL(resp))

		// let kmup calculate language stats
		time.Sleep(time.Second)

		// Save new file to master branch
		req = NewRequest(t, "GET", "/api/v1/repos/user2/repo1/languages")
		resp = MakeRequest(t, req, http.StatusOK)

		var languages map[string]int64
		DecodeJSON(t, resp, &languages)

		assert.InDeltaMapValues(t, map[string]int64{"Go": 12}, languages, 0)
	})
}
