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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallerFuncName(t *testing.T) {
	s := CallerFuncName()
	assert.Equal(t, "github.com/kumose/kmup/modules/util.TestCallerFuncName", s)
}

func BenchmarkCallerFuncName(b *testing.B) {
	// BenchmarkCaller/sprintf-12         	12744829	        95.49 ns/op
	b.Run("sprintf", func(b *testing.B) {
		for b.Loop() {
			_ = fmt.Sprintf("aaaaaaaaaaaaaaaa %s %s %s", "bbbbbbbbbbbbbbbbbbb", b.Name(), "ccccccccccccccccccccc")
		}
	})
	// BenchmarkCaller/caller-12          	10625133	       113.6 ns/op
	// It is almost as fast as fmt.Sprintf
	b.Run("caller", func(b *testing.B) {
		for b.Loop() {
			CallerFuncName()
		}
	})
}
