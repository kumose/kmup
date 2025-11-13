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

package charset

import (
	"bytes"
	"io"
)

// BreakWriter wraps an io.Writer to always write '\n' as '<br>'
type BreakWriter struct {
	io.Writer
}

// Write writes the provided byte slice transparently replacing '\n' with '<br>'
func (b *BreakWriter) Write(bs []byte) (n int, err error) {
	pos := 0
	for pos < len(bs) {
		idx := bytes.IndexByte(bs[pos:], '\n')
		if idx < 0 {
			wn, err := b.Writer.Write(bs[pos:])
			return n + wn, err
		}

		if idx > 0 {
			wn, err := b.Writer.Write(bs[pos : pos+idx])
			n += wn
			if err != nil {
				return n, err
			}
		}

		if _, err = b.Writer.Write([]byte("<br>")); err != nil {
			return n, err
		}
		pos += idx + 1

		n++
	}

	return n, err
}
