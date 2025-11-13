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

package languagestats

import (
	"fmt"
	"os"
	"testing"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

func testRun(m *testing.M) error {
	gitHomePath, err := os.MkdirTemp(os.TempDir(), "git-home")
	if err != nil {
		return fmt.Errorf("unable to create temp dir: %w", err)
	}
	defer util.RemoveAll(gitHomePath)
	setting.Git.HomePath = gitHomePath

	if err = git.InitFull(); err != nil {
		return fmt.Errorf("failed to call Init: %w", err)
	}

	exitCode := m.Run()
	if exitCode != 0 {
		return fmt.Errorf("run test failed, ExitCode=%d", exitCode)
	}
	return nil
}

func TestMain(m *testing.M) {
	if err := testRun(m); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Test failed: %v", err)
		os.Exit(1)
	}
}
