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

package storage

import (
	"fmt"
	"io"
	"net/url"
	"os"
)

var uninitializedStorage = discardStorage("uninitialized storage")

type discardStorage string

func (s discardStorage) Open(_ string) (Object, error) {
	return nil, fmt.Errorf("%s", s)
}

func (s discardStorage) Save(_ string, _ io.Reader, _ int64) (int64, error) {
	return 0, fmt.Errorf("%s", s)
}

func (s discardStorage) Stat(_ string) (os.FileInfo, error) {
	return nil, fmt.Errorf("%s", s)
}

func (s discardStorage) Delete(_ string) error {
	return fmt.Errorf("%s", s)
}

func (s discardStorage) URL(_, _, _ string, _ url.Values) (*url.URL, error) {
	return nil, fmt.Errorf("%s", s)
}

func (s discardStorage) IterateObjects(_ string, _ func(string, Object) error) error {
	return fmt.Errorf("%s", s)
}
