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

package webtheme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseThemeMetaInfo(t *testing.T) {
	m := parseThemeMetaInfoToMap(`kmup-theme-meta-info {
	--k1: "v1";
	--k2: "v\"2";
	--k3: 'v3';
	--k4: 'v\'4';
	--k5: v5;
}`)
	assert.Equal(t, map[string]string{
		"--k1": "v1",
		"--k2": `v"2`,
		"--k3": "v3",
		"--k4": "v'4",
		"--k5": "v5",
	}, m)

	// if an auto theme imports others, the meta info should be extracted from the last one
	// the meta in imported themes should be ignored to avoid incorrect overriding
	m = parseThemeMetaInfoToMap(`
@media (prefers-color-scheme: dark) { kmup-theme-meta-info { --k1: foo; } }
@media (prefers-color-scheme: light) { kmup-theme-meta-info { --k1: bar; } }
kmup-theme-meta-info {
	--k2: real;
}`)
	assert.Equal(t, map[string]string{"--k2": "real"}, m)

	// compressed CSS, no trailing semicolon
	m = parseThemeMetaInfoToMap(`kmup-theme-meta-info{--k1:"v1"}`)
	assert.Equal(t, map[string]string{"--k1": "v1"}, m)
	m = parseThemeMetaInfoToMap(`kmup-theme-meta-info{--k1:"v1";--k2:"v2"}`)
	assert.Equal(t, map[string]string{"--k1": "v1", "--k2": "v2"}, m)
}
