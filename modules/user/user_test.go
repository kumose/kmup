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

package user

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getWhoamiOutput() (string, error) {
	output, err := exec.Command("whoami").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func TestCurrentUsername(t *testing.T) {
	user := CurrentUsername()
	require.NotEmpty(t, user)

	// Windows whoami is weird, so just skip remaining tests
	if runtime.GOOS == "windows" {
		t.Skip("skipped test because of weird whoami on Windows")
	}
	whoami, err := getWhoamiOutput()
	require.NoError(t, err)

	user = CurrentUsername()
	assert.Equal(t, whoami, user)

	t.Setenv("USER", "spoofed")
	user = CurrentUsername()
	assert.Equal(t, whoami, user)
}
