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

package nosql

import "net/url"

// ToLevelDBURI converts old style connections to a LevelDBURI
//
// A LevelDBURI matches the pattern:
//
// leveldb://path[?[option=value]*]
//
// We have previously just provided the path but this prevent other options
func ToLevelDBURI(connection string) *url.URL {
	uri, err := url.Parse(connection)
	if err == nil && uri.Scheme == "leveldb" {
		return uri
	}
	uri, _ = url.Parse("leveldb://common")
	uri.Host = ""
	uri.Path = connection
	return uri
}
