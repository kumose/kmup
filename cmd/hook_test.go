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

package cmd

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPktLine(t *testing.T) {
	// test read
	ctx := t.Context()
	s := strings.NewReader("0000")
	r := bufio.NewReader(s)
	result, err := readPktLine(ctx, r, pktLineTypeFlush)
	assert.NoError(t, err)
	assert.Equal(t, pktLineTypeFlush, result.Type)

	s = strings.NewReader("0006a\n")
	r = bufio.NewReader(s)
	result, err = readPktLine(ctx, r, pktLineTypeData)
	assert.NoError(t, err)
	assert.Equal(t, pktLineTypeData, result.Type)
	assert.Equal(t, []byte("a\n"), result.Data)

	// test write
	w := bytes.NewBuffer([]byte{})
	err = writeFlushPktLine(ctx, w)
	assert.NoError(t, err)
	assert.Equal(t, []byte("0000"), w.Bytes())

	w.Reset()
	err = writeDataPktLine(ctx, w, []byte("a\nb"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("0007a\nb"), w.Bytes())
}
