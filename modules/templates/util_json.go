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

package templates

import (
	"bytes"

	"github.com/kumose/kmup/modules/json"
)

type JsonUtils struct{} //nolint:revive // variable naming triggers on Json, wants JSON

var jsonUtils = JsonUtils{}

func NewJsonUtils() *JsonUtils { //nolint:revive // variable naming triggers on Json, wants JSON
	return &jsonUtils
}

func (su *JsonUtils) EncodeToString(v any) string {
	out, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(out)
}

func (su *JsonUtils) PrettyIndent(s string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(s), "", "  ")
	if err != nil {
		return ""
	}
	return out.String()
}
