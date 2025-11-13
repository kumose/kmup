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

//go:build !windows

package util

import (
	"os"

	"golang.org/x/sys/unix"
)

var defaultUmask int

func init() {
	// at the moment, the umask could only be gotten by calling unix.Umask(newUmask)
	// use 0o077 as temp new umask to reduce the risks if this umask is used anywhere else before the correct umask is recovered
	tempUmask := 0o077
	defaultUmask = unix.Umask(tempUmask)
	unix.Umask(defaultUmask)
}

func ApplyUmask(f string, newMode os.FileMode) error {
	mod := newMode & ^os.FileMode(defaultUmask)
	return os.Chmod(f, mod)
}
