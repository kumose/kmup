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

package eval

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tokens(s string) (a []any) {
	for v := range strings.FieldsSeq(s) {
		a = append(a, v)
	}
	return a
}

func TestEval(t *testing.T) {
	n, err := Expr(0, "/", 0.0)
	assert.NoError(t, err)
	assert.True(t, math.IsNaN(n.Value.(float64)))

	_, err = Expr(nil)
	assert.ErrorContains(t, err, "unsupported token type")
	_, err = Expr([]string{})
	assert.ErrorContains(t, err, "unsupported token type")
	_, err = Expr(struct{}{})
	assert.ErrorContains(t, err, "unsupported token type")

	cases := []struct {
		expr string
		want any
	}{
		{"-1", int64(-1)},
		{"1 + 2", int64(3)},
		{"3 - 2 + 4", int64(5)},
		{"1 + 2 * 3", int64(7)},
		{"1 + ( 2 * 3 )", int64(7)},
		{"( 1 + 2 ) * 3", int64(9)},
		{"( 1 + 2.0 ) / 3", float64(1)},
		{"sum( 1 , 2 , 3 , 4 )", int64(10)},
		{"100 + sum( 1 , 2 + 3 , 0.0 ) / 2", float64(103)},
		{"100 * 5 / ( 5 + 15 )", int64(25)},
		{"9 == 5", int64(0)},
		{"5 == 5", int64(1)},
		{"9 != 5", int64(1)},
		{"5 != 5", int64(0)},
		{"9 > 5", int64(1)},
		{"5 > 9", int64(0)},
		{"5 >= 9", int64(0)},
		{"9 >= 9", int64(1)},
		{"9 < 5", int64(0)},
		{"5 < 9", int64(1)},
		{"9 <= 5", int64(0)},
		{"5 <= 5", int64(1)},
		{"1 and 2", int64(1)}, // Golang template definition: non-zero values are all truth
		{"1 and 0", int64(0)},
		{"0 and 0", int64(0)},
		{"1 or 2", int64(1)},
		{"1 or 0", int64(1)},
		{"0 or 1", int64(1)},
		{"0 or 0", int64(0)},
		{"not 2 == 1", int64(1)},
		{"not not ( 9 < 5 )", int64(0)},
	}

	for _, c := range cases {
		n, err := Expr(tokens(c.expr)...)
		if assert.NoError(t, err, "expr: %s", c.expr) {
			assert.Equal(t, c.want, n.Value)
		}
	}

	bads := []struct {
		expr   string
		errMsg string
	}{
		{"0 / 0", "integer divide by zero"},
		{"1 +", "num stack is empty"},
		{"+ 1", "num stack is empty"},
		{"( 1", "incomplete sub-expression"},
		{"1 )", "op stack is empty"}, // can not find the corresponding open bracket after the stack becomes empty
		{"1 , 2", "expect 1 value as final result"},
		{"( 1 , 2 )", "too many values in one bracket"},
		{"1 a 2", "unknown operator"},
	}
	for _, c := range bads {
		_, err = Expr(tokens(c.expr)...)
		assert.ErrorContains(t, err, c.errMsg, "expr: %s", c.expr)
	}
}
