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

package private

import (
	"strconv"
	"strings"

	"github.com/kumose/kmup/modules/optional"
)

// GitPushOptions is a wrapper around a map[string]string
type GitPushOptions map[string]string

// GitPushOptions keys
const (
	GitPushOptionRepoPrivate  = "repo.private"
	GitPushOptionRepoTemplate = "repo.template"
	GitPushOptionForcePush    = "force-push"
)

// Bool checks for a key in the map and parses as a boolean
// An option without value is considered true, eg: "-o force-push" or "-o repo.private"
func (g GitPushOptions) Bool(key string) optional.Option[bool] {
	if val, ok := g[key]; ok {
		if val == "" {
			return optional.Some(true)
		}
		if b, err := strconv.ParseBool(val); err == nil {
			return optional.Some(b)
		}
	}
	return optional.None[bool]()
}

// AddFromKeyValue adds a key value pair to the map by "key=value" format or "key" for empty value
func (g GitPushOptions) AddFromKeyValue(line string) {
	kv := strings.SplitN(line, "=", 2)
	if len(kv) == 2 {
		g[kv[0]] = kv[1]
	} else {
		g[kv[0]] = ""
	}
}
