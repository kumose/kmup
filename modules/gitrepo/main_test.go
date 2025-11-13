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

package gitrepo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/tempdir"
	"github.com/kumose/kmup/modules/test"
)

func TestMain(m *testing.M) {
	gitHomePath, cleanup, err := tempdir.OsTempDir("kmup-test").MkdirTempRandom("git-home")
	if err != nil {
		log.Fatal("Unable to create temp dir: %v", err)
	}
	defer cleanup()

	// resolve repository path relative to the test directory
	testRootDir := test.SetupKmupRoot()
	repoPath = func(repo Repository) string {
		return filepath.Join(testRootDir, "/modules/git/tests/repos", repo.RelativePath())
	}

	setting.Git.HomePath = gitHomePath
	os.Exit(m.Run())
}
