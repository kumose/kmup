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

package gitcmd

import (
	"fmt"
	"os/exec"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
)

var GitExecutable = "git" // the command name of git, will be updated to an absolute path during initialization

// SetExecutablePath changes the path of git executable and checks the file permission and version.
func SetExecutablePath(path string) error {
	// If path is empty, we use the default value of GitExecutable "git" to search for the location of git.
	if path != "" {
		GitExecutable = path
	}
	absPath, err := exec.LookPath(GitExecutable)
	if err != nil {
		return fmt.Errorf("git not found: %w", err)
	}
	GitExecutable = absPath
	return nil
}

// HomeDir is the home dir for git to store the global config file used by Kmup internally
func HomeDir() string {
	if setting.Git.HomePath == "" {
		// strict check, make sure the git module is initialized correctly.
		// attention: when the git module is called in kmup sub-command (serv/hook), the log module might not obviously show messages to users/developers.
		// for example: if there is kmup git hook code calling NewCommand before git.InitXxx, the integration test won't show the real failure reasons.
		log.Fatal("Unable to init Git's HomeDir, incorrect initialization of the setting and git modules")
		return ""
	}
	return setting.Git.HomePath
}
