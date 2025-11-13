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

package migration

import "fmt"

// ErrNotSupported represents status if a downloader do not supported something.
type ErrNotSupported struct {
	Entity string
}

// IsErrNotSupported checks if an error is an ErrNotSupported
func IsErrNotSupported(err error) bool {
	_, ok := err.(ErrNotSupported)
	return ok
}

// Error return error message
func (err ErrNotSupported) Error() string {
	if len(err.Entity) != 0 {
		return fmt.Sprintf("'%s' not supported", err.Entity)
	}
	return "not supported"
}
