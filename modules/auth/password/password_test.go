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

package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComplexity_IsComplexEnough(t *testing.T) {
	matchComplexityOnce.Do(func() {})

	testlist := []struct {
		complexity  []string
		truevalues  []string
		falsevalues []string
	}{
		{[]string{"off"}, []string{"1", "-", "a", "A", "ñ", "日本語"}, []string{}},
		{[]string{"lower"}, []string{"abc", "abc!"}, []string{"ABC", "123", "=!$", ""}},
		{[]string{"upper"}, []string{"ABC"}, []string{"abc", "123", "=!$", "abc!", ""}},
		{[]string{"digit"}, []string{"123"}, []string{"abc", "ABC", "=!$", "abc!", ""}},
		{[]string{"spec"}, []string{"=!$", "abc!"}, []string{"abc", "ABC", "123", ""}},
		{[]string{"off"}, []string{"abc", "ABC", "123", "=!$", "abc!", ""}, nil},
		{[]string{"lower", "spec"}, []string{"abc!"}, []string{"abc", "ABC", "123", "=!$", "abcABC123", ""}},
		{[]string{"lower", "upper", "digit"}, []string{"abcABC123"}, []string{"abc", "ABC", "123", "=!$", "abc!", ""}},
		{[]string{""}, []string{"abC=1", "abc!9D"}, []string{"ABC", "123", "=!$", ""}},
	}

	for _, test := range testlist {
		testComplextity(test.complexity)
		for _, val := range test.truevalues {
			assert.True(t, IsComplexEnough(val))
		}
		for _, val := range test.falsevalues {
			assert.False(t, IsComplexEnough(val))
		}
	}

	// Remove settings for other tests
	testComplextity([]string{"off"})
}

func TestComplexity_Generate(t *testing.T) {
	matchComplexityOnce.Do(func() {})

	const maxCount = 50
	const pwdLen = 50

	test := func(t *testing.T, modes []string) {
		testComplextity(modes)
		for range maxCount {
			pwd, err := Generate(pwdLen)
			assert.NoError(t, err)
			assert.Len(t, pwd, pwdLen)
			assert.True(t, IsComplexEnough(pwd), "Failed complexities with modes %+v for generated: %s", modes, pwd)
		}
	}

	test(t, []string{"lower"})
	test(t, []string{"upper"})
	test(t, []string{"lower", "upper", "spec"})
	test(t, []string{"off"})
	test(t, []string{""})

	// Remove settings for other tests
	testComplextity([]string{"off"})
}

func testComplextity(values []string) {
	// Cleanup previous values
	validChars = ""
	requiredList = make([]complexity, 0, len(values))
	setupComplexity(values)
}
