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

package context

import (
	"io"
	"net/http"
	"strings"
)

// UploadStream returns the request body or the first form file
// Only form files need to get closed.
func (ctx *Context) UploadStream() (rd io.ReadCloser, needToClose bool, err error) {
	contentType := strings.ToLower(ctx.Req.Header.Get("Content-Type"))
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") || strings.HasPrefix(contentType, "multipart/form-data") {
		if err := ctx.Req.ParseMultipartForm(32 << 20); err != nil {
			return nil, false, err
		}
		if ctx.Req.MultipartForm.File == nil {
			return nil, false, http.ErrMissingFile
		}
		for _, files := range ctx.Req.MultipartForm.File {
			if len(files) > 0 {
				r, err := files[0].Open()
				return r, true, err
			}
		}
		return nil, false, http.ErrMissingFile
	}
	return ctx.Req.Body, false, nil
}
