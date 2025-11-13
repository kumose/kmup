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

import "strconv"

// Option is a generic type that can hold a value of type T or be empty (None).
//
// It must use the slice type to work with "chi" form values binding:
// * non-existing value are represented as an empty slice (None)
// * existing value is represented as a slice with one element (Some)
// * multiple values are represented as a slice with multiple elements (Some), the Value is the first element (not well-defined in this case)
type Option[T any] []T

func None[T any]() Option[T] {
	return nil
}

func Some[T any](v T) Option[T] {
	return Option[T]{v}
}

func FromPtr[T any](v *T) Option[T] {
	if v == nil {
		return None[T]()
	}
	return Some(*v)
}

func FromMapLookup[K comparable, V any](m map[K]V, k K) Option[V] {
	if v, ok := m[k]; ok {
		return Some(v)
	}
	return None[V]()
}

func FromNonDefault[T comparable](v T) Option[T] {
	var zero T
	if v == zero {
		return None[T]()
	}
	return Some(v)
}

func (o Option[T]) Has() bool {
	return o != nil
}

func (o Option[T]) Value() T {
	var zero T
	return o.ValueOrDefault(zero)
}

func (o Option[T]) ValueOrDefault(v T) T {
	if o.Has() {
		return o[0]
	}
	return v
}

// ParseBool get the corresponding optional.Option[bool] of a string using strconv.ParseBool
func ParseBool(s string) Option[bool] {
	v, e := strconv.ParseBool(s)
	if e != nil {
		return None[bool]()
	}
	return Some(v)
}
