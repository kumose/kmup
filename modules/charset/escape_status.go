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

package charset

// EscapeStatus represents the findings of the unicode escaper
type EscapeStatus struct {
	Escaped      bool
	HasError     bool
	HasBadRunes  bool
	HasInvisible bool
	HasAmbiguous bool
}

// Or combines two EscapeStatus structs into one representing the conjunction of the two
func (status *EscapeStatus) Or(other *EscapeStatus) *EscapeStatus {
	st := status
	if status == nil {
		st = &EscapeStatus{}
	}
	st.Escaped = st.Escaped || other.Escaped
	st.HasError = st.HasError || other.HasError
	st.HasBadRunes = st.HasBadRunes || other.HasBadRunes
	st.HasAmbiguous = st.HasAmbiguous || other.HasAmbiguous
	st.HasInvisible = st.HasInvisible || other.HasInvisible
	return st
}
