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

package process

import (
	"os/exec"
	"syscall"
)

// SetSysProcAttribute sets the common SysProcAttrs for commands
func SetSysProcAttribute(cmd *exec.Cmd) {
	// When Kmup runs SubProcessA -> SubProcessB and SubProcessA gets killed by context timeout, use setpgid to make sure the sub processes can be reaped instead of leaving defunct(zombie) processes.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
