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

	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestExploreRepos(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	req := NewRequest(t, "GET", "/explore/repos?q=TheKeyword&topic=1&language=TheLang")
	resp := MakeRequest(t, req, http.StatusOK)
	respStr := resp.Body.String()

	assert.Contains(t, respStr, `<input type="hidden" name="topic" value="true">`)
	assert.Contains(t, respStr, `<input type="hidden" name="language" value="TheLang">`)
	assert.Contains(t, respStr, `<input type="search" name="q" value="TheKeyword"`)
}
