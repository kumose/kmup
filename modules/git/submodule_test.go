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

package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kumose/kmup/modules/git/gitcmd"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTemplateSubmoduleCommits(t *testing.T) {
	testRepoPath := filepath.Join(testReposDir, "repo4_submodules")
	submodules, err := GetTemplateSubmoduleCommits(t.Context(), testRepoPath)
	require.NoError(t, err)

	assert.Len(t, submodules, 2)

	assert.Equal(t, "<°)))><", submodules[0].Path)
	assert.Equal(t, "d2932de67963f23d43e1c7ecf20173e92ee6c43c", submodules[0].Commit)

	assert.Equal(t, "libtest", submodules[1].Path)
	assert.Equal(t, "1234567890123456789012345678901234567890", submodules[1].Commit)
}

func TestAddTemplateSubmoduleIndexes(t *testing.T) {
	ctx := t.Context()
	tmpDir := t.TempDir()
	var err error
	_, _, err = gitcmd.NewCommand("init").WithDir(tmpDir).RunStdString(ctx)
	require.NoError(t, err)
	_ = os.Mkdir(filepath.Join(tmpDir, "new-dir"), 0o755)
	err = AddTemplateSubmoduleIndexes(ctx, tmpDir, []TemplateSubmoduleCommit{{Path: "new-dir", Commit: "1234567890123456789012345678901234567890"}})
	require.NoError(t, err)
	_, _, err = gitcmd.NewCommand("add", "--all").WithDir(tmpDir).RunStdString(ctx)
	require.NoError(t, err)
	_, _, err = gitcmd.NewCommand("-c", "user.name=a", "-c", "user.email=b", "commit", "-m=test").WithDir(tmpDir).RunStdString(ctx)
	require.NoError(t, err)
	submodules, err := GetTemplateSubmoduleCommits(t.Context(), tmpDir)
	require.NoError(t, err)
	assert.Len(t, submodules, 1)
	assert.Equal(t, "new-dir", submodules[0].Path)
	assert.Equal(t, "1234567890123456789012345678901234567890", submodules[0].Commit)
}
