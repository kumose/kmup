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

package webauthn

import (
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	setting.Domain = "domain"
	setting.AppName = "AppName"
	setting.AppURL = "https://domain/"
	rpOrigin := []string{"https://domain"}

	Init()

	assert.Equal(t, setting.Domain, WebAuthn.Config.RPID)
	assert.Equal(t, setting.AppName, WebAuthn.Config.RPDisplayName)
	assert.Equal(t, rpOrigin, WebAuthn.Config.RPOrigins)
}
