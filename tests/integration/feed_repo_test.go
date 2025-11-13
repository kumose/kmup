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
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestFeedRepo(t *testing.T) {
	t.Run("RSS", func(t *testing.T) {
		defer tests.PrepareTestEnv(t)()

		req := NewRequest(t, "GET", "/user2/repo1.rss")
		resp := MakeRequest(t, req, http.StatusOK)

		data := resp.Body.String()
		assert.Contains(t, data, `<rss version="2.0"`)

		var rss RSS
		err := xml.Unmarshal(resp.Body.Bytes(), &rss)
		assert.NoError(t, err)
		assert.Contains(t, rss.Channel.Link, "/user2/repo1")
		assert.NotEmpty(t, rss.Channel.PubDate)
		assert.Len(t, rss.Channel.Items, 1)
		assert.Equal(t, "issue5", rss.Channel.Items[0].Description)
		assert.NotEmpty(t, rss.Channel.Items[0].PubDate)
	})
}
