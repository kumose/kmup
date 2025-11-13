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

package optional

import (
	"github.com/kumose/kmup/modules/json"

	"gopkg.in/yaml.v3"
)

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	var v *T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = FromPtr(v)
	return nil
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.Has() {
		return []byte("null"), nil
	}

	return json.Marshal(o.Value())
}

func (o *Option[T]) UnmarshalYAML(value *yaml.Node) error {
	var v *T
	if err := value.Decode(&v); err != nil {
		return err
	}
	*o = FromPtr(v)
	return nil
}

func (o Option[T]) MarshalYAML() (any, error) {
	if !o.Has() {
		return nil, nil
	}

	value := new(yaml.Node)
	err := value.Encode(o.Value())
	return value, err
}
