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

package setting

import "strings"

// MimeTypeMap defines custom mime type mapping settings
var MimeTypeMap = struct {
	Enabled bool
	Map     map[string]string
}{
	Enabled: false,
	Map:     map[string]string{},
}

func loadMimeTypeMapFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("repository.mimetype_mapping")
	keys := sec.Keys()
	m := make(map[string]string, len(keys))
	for _, key := range keys {
		m[strings.ToLower(key.Name())] = key.Value()
	}
	MimeTypeMap.Map = m
	if len(keys) > 0 {
		MimeTypeMap.Enabled = true
	}
}
