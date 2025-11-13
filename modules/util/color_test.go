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
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HexToRBGColor(t *testing.T) {
	cases := []struct {
		colorString string
		expectedR   float64
		expectedG   float64
		expectedB   float64
	}{
		{"2b8685", 43, 134, 133},
		{"1e1", 17, 238, 17},
		{"#1e1", 17, 238, 17},
		{"1e16", 17, 238, 17},
		{"3bb6b3", 59, 182, 179},
		{"#3bb6b399", 59, 182, 179},
		{"#0", 0, 0, 0},
		{"#00000", 0, 0, 0},
		{"#1234567", 0, 0, 0},
	}
	for n, c := range cases {
		r, g, b := HexToRBGColor(c.colorString)
		assert.InDelta(t, c.expectedR, r, 0, "case %d: error R should match: expected %f, but get %f", n, c.expectedR, r)
		assert.InDelta(t, c.expectedG, g, 0, "case %d: error G should match: expected %f, but get %f", n, c.expectedG, g)
		assert.InDelta(t, c.expectedB, b, 0, "case %d: error B should match: expected %f, but get %f", n, c.expectedB, b)
	}
}

func Test_UseLightText(t *testing.T) {
	cases := []struct {
		color    string
		expected string
	}{
		{"#d73a4a", "#fff"},
		{"#0075ca", "#fff"},
		{"#cfd3d7", "#000"},
		{"#a2eeef", "#000"},
		{"#7057ff", "#fff"},
		{"#008672", "#fff"},
		{"#e4e669", "#000"},
		{"#d876e3", "#000"},
		{"#ffffff", "#000"},
		{"#2b8684", "#fff"},
		{"#2b8786", "#fff"},
		{"#2c8786", "#000"},
		{"#3bb6b3", "#000"},
		{"#7c7268", "#fff"},
		{"#7e716c", "#fff"},
		{"#81706d", "#fff"},
		{"#807070", "#fff"},
		{"#84b6eb", "#000"},
	}
	for n, c := range cases {
		assert.Equal(t, c.expected, ContrastColor(c.color), "case %d: error should match", n)
	}
}
