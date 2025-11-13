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

package lfs

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// DownloadCallback gets called for every requested LFS object to process its content
type DownloadCallback func(p Pointer, content io.ReadCloser, objectError error) error

// UploadCallback gets called for every requested LFS object to provide its content
type UploadCallback func(p Pointer, objectError error) (io.ReadCloser, error)

// Client is used to communicate with a LFS source
type Client interface {
	BatchSize() int
	Download(ctx context.Context, objects []Pointer, callback DownloadCallback) error
	Upload(ctx context.Context, objects []Pointer, callback UploadCallback) error
}

// NewClient creates a LFS client
func NewClient(endpoint *url.URL, httpTransport *http.Transport) Client {
	if endpoint.Scheme == "file" {
		return newFilesystemClient(endpoint)
	}
	return newHTTPClient(endpoint, httpTransport)
}
