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

package json

import (
	"bytes"
	"encoding/json" //nolint:depguard // this package wraps it
	"io"
)

type jsonV1 struct{}

var _ Interface = jsonV1{}

func (jsonV1) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonV1) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (jsonV1) NewEncoder(writer io.Writer) Encoder {
	return json.NewEncoder(writer)
}

func (jsonV1) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

func (jsonV1) Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	return json.Indent(dst, src, prefix, indent)
}
