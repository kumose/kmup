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

package util

import (
	"bytes"
	"encoding/gob"
)

// PackData uses gob to encode the given data in sequence
func PackData(data ...any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	for _, datum := range data {
		if err := enc.Encode(datum); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// UnpackData uses gob to decode the given data in sequence
func UnpackData(buf []byte, data ...any) error {
	r := bytes.NewReader(buf)
	enc := gob.NewDecoder(r)
	for _, datum := range data {
		if err := enc.Decode(datum); err != nil {
			return err
		}
	}
	return nil
}
