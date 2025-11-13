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

package markup

import (
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestCamoHandleLink(t *testing.T) {
	setting.AppURL = "https://kmup.com"
	// Test media proxy
	setting.Camo.Enabled = true
	setting.Camo.ServerURL = "https://image.proxy"
	setting.Camo.HMACKey = "geheim"

	assert.Equal(t,
		"https://kmup.com/img.jpg",
		camoHandleLink("https://kmup.com/img.jpg"))
	assert.Equal(t,
		"https://testimages.org/img.jpg",
		camoHandleLink("https://testimages.org/img.jpg"))
	assert.Equal(t,
		"https://image.proxy/eivin43gJwGVIjR9MiYYtFIk0mw/aHR0cDovL3Rlc3RpbWFnZXMub3JnL2ltZy5qcGc",
		camoHandleLink("http://testimages.org/img.jpg"))

	setting.Camo.Always = true
	assert.Equal(t,
		"https://kmup.com/img.jpg",
		camoHandleLink("https://kmup.com/img.jpg"))
	assert.Equal(t,
		"https://image.proxy/tkdlvmqpbIr7SjONfHNgEU622y0/aHR0cHM6Ly90ZXN0aW1hZ2VzLm9yZy9pbWcuanBn",
		camoHandleLink("https://testimages.org/img.jpg"))
	assert.Equal(t,
		"https://image.proxy/eivin43gJwGVIjR9MiYYtFIk0mw/aHR0cDovL3Rlc3RpbWFnZXMub3JnL2ltZy5qcGc",
		camoHandleLink("http://testimages.org/img.jpg"))

	// Restore previous settings
	setting.Camo.Enabled = false
}
