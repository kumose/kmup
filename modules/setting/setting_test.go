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

package setting

import (
	"testing"

	"github.com/kumose/kmup/modules/json"

	"github.com/stretchr/testify/assert"
)

func TestMakeAbsoluteAssetURL(t *testing.T) {
	assert.Equal(t, "https://localhost:2345", MakeAbsoluteAssetURL("https://localhost:1234", "https://localhost:2345"))
	assert.Equal(t, "https://localhost:2345", MakeAbsoluteAssetURL("https://localhost:1234/", "https://localhost:2345"))
	assert.Equal(t, "https://localhost:2345", MakeAbsoluteAssetURL("https://localhost:1234/", "https://localhost:2345/"))
	assert.Equal(t, "https://localhost:1234/foo", MakeAbsoluteAssetURL("https://localhost:1234", "/foo"))
	assert.Equal(t, "https://localhost:1234/foo", MakeAbsoluteAssetURL("https://localhost:1234/", "/foo"))
	assert.Equal(t, "https://localhost:1234/foo", MakeAbsoluteAssetURL("https://localhost:1234/", "/foo/"))
	assert.Equal(t, "https://localhost:1234/foo", MakeAbsoluteAssetURL("https://localhost:1234/foo", "/foo"))
	assert.Equal(t, "https://localhost:1234/foo", MakeAbsoluteAssetURL("https://localhost:1234/foo/", "/foo"))
	assert.Equal(t, "https://localhost:1234/foo", MakeAbsoluteAssetURL("https://localhost:1234/foo/", "/foo/"))
	assert.Equal(t, "https://localhost:1234/bar", MakeAbsoluteAssetURL("https://localhost:1234/foo", "/bar"))
	assert.Equal(t, "https://localhost:1234/bar", MakeAbsoluteAssetURL("https://localhost:1234/foo/", "/bar"))
	assert.Equal(t, "https://localhost:1234/bar", MakeAbsoluteAssetURL("https://localhost:1234/foo/", "/bar/"))
}

func TestMakeManifestData(t *testing.T) {
	jsonBytes := MakeManifestData(`Example App '\"`, "https://example.com", "https://example.com/foo/bar")
	assert.True(t, json.Valid(jsonBytes))
}
