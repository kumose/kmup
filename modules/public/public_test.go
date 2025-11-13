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

package public

import (
	"testing"

	"github.com/kumose/kmup/modules/container"

	"github.com/stretchr/testify/assert"
)

func TestParseAcceptEncoding(t *testing.T) {
	kases := []struct {
		Header   string
		Expected container.Set[string]
	}{
		{
			Header:   "deflate, gzip;q=1.0, *;q=0.5",
			Expected: container.SetOf("deflate", "gzip"),
		},
		{
			Header:   " gzip, deflate, br",
			Expected: container.SetOf("deflate", "gzip", "br"),
		},
	}

	for _, kase := range kases {
		t.Run(kase.Header, func(t *testing.T) {
			assert.EqualValues(t, kase.Expected, parseAcceptEncoding(kase.Header))
		})
	}
}
