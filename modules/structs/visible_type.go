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

package structs

// VisibleType defines the visibility of user and org
type VisibleType int

const (
	// VisibleTypePublic Visible for everyone
	VisibleTypePublic VisibleType = iota

	// VisibleTypeLimited Visible for every connected user
	VisibleTypeLimited

	// VisibleTypePrivate Visible only for self or admin user
	VisibleTypePrivate
)

// VisibilityModes is a map of Visibility types
var VisibilityModes = map[string]VisibleType{
	"public":  VisibleTypePublic,
	"limited": VisibleTypeLimited,
	"private": VisibleTypePrivate,
}

// IsPublic returns true if VisibleType is public
func (vt VisibleType) IsPublic() bool {
	return vt == VisibleTypePublic
}

// IsLimited returns true if VisibleType is limited
func (vt VisibleType) IsLimited() bool {
	return vt == VisibleTypeLimited
}

// IsPrivate returns true if VisibleType is private
func (vt VisibleType) IsPrivate() bool {
	return vt == VisibleTypePrivate
}

// VisibilityString provides the mode string of the visibility type (public, limited, private)
func (vt VisibleType) String() string {
	for k, v := range VisibilityModes {
		if vt == v {
			return k
		}
	}
	return ""
}

// ExtractKeysFromMapString provides a slice of keys from map
func ExtractKeysFromMapString(in map[string]VisibleType) (keys []string) {
	for k := range in {
		keys = append(keys, k)
	}
	return keys
}
