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

//go:build !bindata

package migration

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func openSchema(s string) (io.ReadCloser, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	basename := path.Base(u.Path)
	filename := basename
	//
	// Schema reference each other within the schemas directory but
	// the tests run in the parent directory.
	//
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = filepath.Join("schemas", basename)
		//
		// Integration tests run from the git root directory, not the
		// directory in which the test source is located.
		//
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			filename = filepath.Join("modules/migration/schemas", basename)
		}
	}
	return os.Open(filename)
}
