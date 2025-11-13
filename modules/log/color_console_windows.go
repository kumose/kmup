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

package log

import (
	"os"

	"github.com/mattn/go-isatty"
	"golang.org/x/sys/windows"
)

func enableVTMode(console windows.Handle) bool {
	mode := uint32(0)
	err := windows.GetConsoleMode(console, &mode)
	if err != nil {
		return false
	}

	// EnableVirtualTerminalProcessing is the console mode to allow ANSI code
	// interpretation on the console. See:
	// https://docs.microsoft.com/en-us/windows/console/setconsolemode
	// It only works on Windows 10. Earlier terminals will fail with an err which we will
	// handle to say don't color
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	err = windows.SetConsoleMode(console, mode)
	return err == nil
}

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		CanColorStdout = enableVTMode(windows.Stdout)
	} else {
		CanColorStdout = isatty.IsCygwinTerminal(os.Stderr.Fd())
	}

	if isatty.IsTerminal(os.Stderr.Fd()) {
		CanColorStderr = enableVTMode(windows.Stderr)
	} else {
		CanColorStderr = isatty.IsCygwinTerminal(os.Stderr.Fd())
	}
}
