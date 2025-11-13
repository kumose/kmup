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

package arch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareVersions(t *testing.T) {
	// https://man.archlinux.org/man/vercmp.8.en
	checks := [][]string{
		{"1.0a", "1.0b", "1.0beta", "1.0p", "1.0pre", "1.0rc", "1.0", "1.0.a", "1.0.1"},
		{"1", "1.0", "1.1", "1.1.1", "1.2", "2.0", "3.0.0"},
	}
	for _, check := range checks {
		for i := 0; i < len(check)-1; i++ {
			require.Equal(t, -1, compareVersions(check[i], check[i+1]))
			require.Equal(t, 1, compareVersions(check[i+1], check[i]))
		}
	}
	require.Equal(t, 1, compareVersions("1.0-2", "1.0"))
	require.Equal(t, 0, compareVersions("0:1.0-1", "1.0"))
	require.Equal(t, 1, compareVersions("1:1.0-1", "2.0"))
}
