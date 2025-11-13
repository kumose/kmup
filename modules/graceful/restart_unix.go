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

// This code is heavily inspired by the archived gofacebook/gracenet/net.go handler

//go:build !windows

package graceful

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var killParent sync.Once

// KillParent sends the kill signal to the parent process if we are a child
func KillParent() {
	killParent.Do(func() {
		if GetManager().IsChild() {
			ppid := syscall.Getppid()
			if ppid > 1 {
				_ = syscall.Kill(ppid, syscall.SIGTERM)
			}
		}
	})
}

// RestartProcess starts a new process passing it the active listeners. It
// doesn't fork, but starts a new process using the same environment and
// arguments as when it was originally started. This allows for a newly
// deployed binary to be started. It returns the pid of the newly started
// process when successful.
func RestartProcess() (int, error) {
	listeners := getActiveListeners()

	// Extract the fds from the listeners.
	files := make([]*os.File, len(listeners))
	for i, l := range listeners {
		var err error
		// Now, all our listeners actually have File() functions so instead of
		// individually casting we just use a hacky interface
		files[i], err = l.(filer).File()
		if err != nil {
			return 0, err
		}

		if unixListener, ok := l.(*net.UnixListener); ok {
			unixListener.SetUnlinkOnClose(false)
		}
		// Remember to close these at the end.
		defer func(i int) {
			_ = files[i].Close()
		}(i)
	}

	// Use the original binary location. This works with symlinks such that if
	// the file it points to has been changed we will use the updated symlink.
	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return 0, err
	}

	// Pass on the environment and replace the old count key with the new one.
	var env []string
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, listenFDsEnv+"=") {
			env = append(env, v)
		}
	}
	env = append(env, fmt.Sprintf("%s=%d", listenFDsEnv, len(listeners)))

	if notifySocketAddr != "" {
		env = append(env, fmt.Sprintf("%s=%s", notifySocketEnv, notifySocketAddr))
	}

	if watchdogTimeout != 0 {
		watchdogStr := strconv.FormatInt(int64(watchdogTimeout/time.Millisecond), 10)
		env = append(env, fmt.Sprintf("%s=%s", watchdogTimeoutEnv, watchdogStr))
	}

	sb := &strings.Builder{}
	for i, unlink := range getActiveListenersToUnlink() {
		if !unlink {
			continue
		}
		_, _ = sb.WriteString(strconv.Itoa(i))
		_, _ = sb.WriteString(",")
	}
	unlinkStr := sb.String()
	if len(unlinkStr) > 0 {
		unlinkStr = unlinkStr[:len(unlinkStr)-1]
		env = append(env, fmt.Sprintf("%s=%s", unlinkFDsEnv, unlinkStr))
	}

	allFiles := append([]*os.File{os.Stdin, os.Stdout, os.Stderr}, files...)
	process, err := os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   originalWD,
		Env:   env,
		Files: allFiles,
	})
	if err != nil {
		return 0, err
	}
	processPid := process.Pid
	_ = process.Release() // no wait, so release
	return processPid, nil
}
