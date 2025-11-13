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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	cases := map[string]string{
		// all old cases from the legacy package
		"HTTPServer":         "http_server",
		"_camelCase":         "_camel_case",
		"NoHTTPS":            "no_https",
		"Wi_thF":             "wi_th_f",
		"_AnotherTES_TCaseP": "_another_tes_t_case_p",
		"ALL":                "all",
		"_HELLO_WORLD_":      "_hello_world_",
		"HELLO_WORLD":        "hello_world",
		"HELLO____WORLD":     "hello____world",
		"TW":                 "tw",
		"_C":                 "_c",

		"  sentence case  ": "__sentence_case__",
		" Mixed-hyphen case _and SENTENCE_case and UPPER-case": "_mixed_hyphen_case__and_sentence_case_and_upper_case",

		// new cases
		" ":        "_",
		"A":        "a",
		"A0":       "a0",
		"a0":       "a0",
		"Aa0":      "aa0",
		"啊":        "啊",
		"A啊":       "a啊",
		"Aa啊b":     "aa啊b",
		"A啊B":      "a啊_b",
		"Aa啊B":     "aa啊_b",
		"TheCase2": "the_case2",
		"ObjIDs":   "obj_i_ds", // the strange database column name which already exists
	}
	for input, expected := range cases {
		assert.Equal(t, expected, ToSnakeCase(input))
	}
}

func TestSplitTrimSpace(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, SplitTrimSpace("a\nb\nc", "\n"))
	assert.Equal(t, []string{"a", "b"}, SplitTrimSpace("\r\na\n\r\nb\n\n", "\n"))
}
