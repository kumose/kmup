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

package upload

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	testContent := []byte(`This is a plain text file.`)
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(testContent)
	w.Close()

	kases := []struct {
		data         []byte
		fileName     string
		allowedTypes string
		err          error
	}{
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "dir/test.txt",
			allowedTypes: "",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "../../../test.txt",
			allowedTypes: "",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: ",",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "|",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "*/*",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "*/*,",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "*/*|",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "text/plain",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "dir/test.txt",
			allowedTypes: "text/plain",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "/dir.txt/test.js",
			allowedTypes: ".js",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: " text/plain ",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: ".txt",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: " .txt,.js",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: " .txt|.js",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "../../test.txt",
			allowedTypes: " .txt|.js",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: " .txt ,.js ",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "text/plain, .txt",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "text/*",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "text/*,.js",
			err:          nil,
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "text/**",
			err:          ErrFileTypeForbidden{"text/plain; charset=utf-8"},
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: "application/x-gzip",
			err:          ErrFileTypeForbidden{"text/plain; charset=utf-8"},
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: ".zip",
			err:          ErrFileTypeForbidden{"text/plain; charset=utf-8"},
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: ".zip,.txtx",
			err:          ErrFileTypeForbidden{"text/plain; charset=utf-8"},
		},
		{
			data:         testContent,
			fileName:     "test.txt",
			allowedTypes: ".zip|.txtx",
			err:          ErrFileTypeForbidden{"text/plain; charset=utf-8"},
		},
		{
			data:         b.Bytes(),
			fileName:     "test.txt",
			allowedTypes: "application/x-gzip",
			err:          nil,
		},
	}

	for _, kase := range kases {
		assert.Equal(t, kase.err, Verify(kase.data, kase.fileName, kase.allowedTypes))
	}
}
