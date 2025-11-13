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

package log

import (
	"os"

	"github.com/mattn/go-isatty"
)

func init() {
	// when running kmup as a systemd unit with logging set to console, the output can not be colorized,
	// otherwise it spams the journal / syslog with escape sequences like "#033[0m#033[32mcmd/web.go:102:#033[32m"
	// this file covers non-windows platforms.
	CanColorStdout = isatty.IsTerminal(os.Stdout.Fd())
	CanColorStderr = isatty.IsTerminal(os.Stderr.Fd())
}
