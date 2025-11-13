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

package emoji

import (
	"testing"

	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	a := FromCode("\U0001f37a")
	b := FromCode("🍺")
	c := FromAlias(":beer:")
	d := FromAlias("beer")

	assert.Equal(t, a, b)
	assert.Equal(t, b, c)
	assert.Equal(t, c, d)

	m := FromCode("\U0001f44d")
	n := FromAlias(":thumbsup:")
	o := FromAlias("+1")

	assert.Equal(t, m, n)
	assert.Equal(t, m, o)

	defer test.MockVariableValue(&setting.UI.EnabledEmojisSet, container.SetOf("thumbsup"))()
	defer globalVarsStore.Store(nil)
	globalVarsStore.Store(nil)
	a = FromCode("\U0001f37a")
	c = FromAlias(":beer:")
	m = FromCode("\U0001f44d")
	n = FromAlias(":thumbsup:")
	o = FromAlias("+1")
	assert.Nil(t, a)
	assert.Nil(t, c)
	assert.NotNil(t, m)
	assert.NotNil(t, n)
	assert.Nil(t, o)
}

func TestReplacers(t *testing.T) {
	tests := []struct {
		f      func(string) string
		v, exp string
	}{
		{ReplaceCodes, ":thumbsup: +1 for \U0001f37a! 🍺 \U0001f44d", ":thumbsup: +1 for :beer:! :beer: :+1:"},
		{ReplaceAliases, ":thumbsup: +1 :+1: :beer:", "\U0001f44d +1 \U0001f44d \U0001f37a"},
	}

	for i, x := range tests {
		s := x.f(x.v)
		assert.Equalf(t, x.exp, s, "test %d `%s` expected `%s`, got: `%s`", i, x.v, x.exp, s)
	}
}

func TestFindEmojiSubmatchIndex(t *testing.T) {
	type testcase struct {
		teststring string
		expected   []int
	}

	testcases := []testcase{
		{
			"\U0001f44d",
			[]int{0, len("\U0001f44d")},
		},
		{
			"\U0001f44d +1 \U0001f44d \U0001f37a",
			[]int{0, 4},
		},
		{
			" \U0001f44d",
			[]int{1, 1 + len("\U0001f44d")},
		},
		{
			string([]byte{'\u0001'}) + "\U0001f44d",
			[]int{1, 1 + len("\U0001f44d")},
		},
	}

	for _, kase := range testcases {
		actual := FindEmojiSubmatchIndex(kase.teststring)
		assert.Equal(t, kase.expected, actual)
	}
}
