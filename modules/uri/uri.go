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

package uri

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ErrURISchemeNotSupported represents a scheme error
type ErrURISchemeNotSupported struct {
	Scheme string
}

func (e ErrURISchemeNotSupported) Error() string {
	return fmt.Sprintf("Unsupported scheme: %v", e.Scheme)
}

// Open open a local file or a remote file
func Open(uriStr string) (io.ReadCloser, error) {
	u, err := url.Parse(uriStr)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		f, err := http.Get(uriStr)
		if err != nil {
			return nil, err
		}
		return f.Body, nil
	case "file":
		return os.Open(u.Path)
	default:
		return nil, ErrURISchemeNotSupported{Scheme: u.Scheme}
	}
}
