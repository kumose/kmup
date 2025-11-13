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
	"os"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/tempdir"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gitHomePath, cleanup, err := tempdir.OsTempDir("kmup-test").MkdirTempRandom("git-home")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to create temp dir: %v", err)
		os.Exit(1)
	}
	defer cleanup()

	setting.Git.HomePath = gitHomePath
	os.Exit(m.Run())
}

func TestRunWithContextStd(t *testing.T) {
	{
		cmd := NewCommand("--version")
		stdout, stderr, err := cmd.RunStdString(t.Context())
		assert.NoError(t, err)
		assert.Empty(t, stderr)
		assert.Contains(t, stdout, "git version")
	}

	{
		cmd := NewCommand("ls-tree", "no-such")
		stdout, stderr, err := cmd.RunStdString(t.Context())
		if assert.Error(t, err) {
			assert.Equal(t, stderr, err.Stderr())
			assert.Equal(t, "fatal: Not a valid object name no-such\n", err.Stderr())
			// FIXME: GIT-CMD-STDERR: it is a bad design, the stderr should not be put in the error message
			assert.Equal(t, "exit status 128 - fatal: Not a valid object name no-such\n", err.Error())
			assert.Empty(t, stdout)
		}
	}

	{
		cmd := NewCommand("ls-tree", "no-such")
		stdout, stderr, err := cmd.RunStdBytes(t.Context())
		if assert.Error(t, err) {
			assert.Equal(t, string(stderr), err.Stderr())
			assert.Equal(t, "fatal: Not a valid object name no-such\n", err.Stderr())
			// FIXME: GIT-CMD-STDERR: it is a bad design, the stderr should not be put in the error message
			assert.Equal(t, "exit status 128 - fatal: Not a valid object name no-such\n", err.Error())
			assert.Empty(t, stdout)
		}
	}

	{
		cmd := NewCommand()
		cmd.AddDynamicArguments("-test")
		assert.ErrorIs(t, cmd.Run(t.Context()), ErrBrokenCommand)

		cmd = NewCommand()
		cmd.AddDynamicArguments("--test")
		assert.ErrorIs(t, cmd.Run(t.Context()), ErrBrokenCommand)
	}

	{
		subCmd := "version"
		cmd := NewCommand().AddDynamicArguments(subCmd) // for test purpose only, the sub-command should never be dynamic for production
		stdout, stderr, err := cmd.RunStdString(t.Context())
		assert.NoError(t, err)
		assert.Empty(t, stderr)
		assert.Contains(t, stdout, "git version")
	}
}

func TestGitArgument(t *testing.T) {
	assert.True(t, isValidArgumentOption("-x"))
	assert.True(t, isValidArgumentOption("--xx"))
	assert.False(t, isValidArgumentOption(""))
	assert.False(t, isValidArgumentOption("x"))

	assert.True(t, isSafeArgumentValue(""))
	assert.True(t, isSafeArgumentValue("x"))
	assert.False(t, isSafeArgumentValue("-x"))
}

func TestCommandString(t *testing.T) {
	cmd := NewCommand("a", "-m msg", "it's a test", `say "hello"`)
	assert.Equal(t, cmd.prog+` a "-m msg" "it's a test" "say \"hello\""`, cmd.LogString())

	cmd = NewCommand("url: https://a:b@c/", "/root/dir-a/dir-b")
	assert.Equal(t, cmd.prog+` "url: https://sanitized-credential@c/" .../dir-a/dir-b`, cmd.LogString())
}
