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
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorTranslatable(t *testing.T) {
	var err error

	err = ErrorWrapTranslatable(io.EOF, "key", 1)
	assert.ErrorIs(t, err, io.EOF)
	assert.Equal(t, "EOF", err.Error())
	assert.Equal(t, "key", err.(*errorTranslatableWrapper).trKey)
	assert.Equal(t, []any{1}, err.(*errorTranslatableWrapper).trArgs)

	err = ErrorWrap(err, "new msg %d", 100)
	assert.ErrorIs(t, err, io.EOF)
	assert.Equal(t, "new msg 100", err.Error())

	errTr := ErrorAsTranslatable(err)
	assert.Equal(t, "EOF", errTr.Error())
	assert.Equal(t, "key", errTr.(*errorTranslatableWrapper).trKey)
}
