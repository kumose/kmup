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
	"strings"
	"testing"

	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func gitConfigContains(sub string) bool {
	if b, err := os.ReadFile(gitcmd.HomeDir() + "/.gitconfig"); err == nil {
		return strings.Contains(string(b), sub)
	}
	return false
}

func TestGitConfig(t *testing.T) {
	ctx := t.Context()
	assert.False(t, gitConfigContains("key-a"))

	assert.NoError(t, configSetNonExist(ctx, "test.key-a", "val-a"))
	assert.True(t, gitConfigContains("key-a = val-a"))

	assert.NoError(t, configSetNonExist(ctx, "test.key-a", "val-a-changed"))
	assert.False(t, gitConfigContains("key-a = val-a-changed"))

	assert.NoError(t, configSet(ctx, "test.key-a", "val-a-changed"))
	assert.True(t, gitConfigContains("key-a = val-a-changed"))

	assert.NoError(t, configAddNonExist(ctx, "test.key-b", "val-b"))
	assert.True(t, gitConfigContains("key-b = val-b"))

	assert.NoError(t, configAddNonExist(ctx, "test.key-b", "val-2b"))
	assert.True(t, gitConfigContains("key-b = val-b"))
	assert.True(t, gitConfigContains("key-b = val-2b"))

	assert.NoError(t, configUnsetAll(ctx, "test.key-b", "val-b"))
	assert.False(t, gitConfigContains("key-b = val-b"))
	assert.True(t, gitConfigContains("key-b = val-2b"))

	assert.NoError(t, configUnsetAll(ctx, "test.key-b", "val-2b"))
	assert.False(t, gitConfigContains("key-b = val-2b"))

	assert.NoError(t, configSet(ctx, "test.key-x", "*"))
	assert.True(t, gitConfigContains("key-x = *"))
	assert.NoError(t, configSetNonExist(ctx, "test.key-x", "*"))
	assert.NoError(t, configUnsetAll(ctx, "test.key-x", "*"))
	assert.False(t, gitConfigContains("key-x = *"))
}

func TestSyncConfig(t *testing.T) {
	oldGitConfig := setting.GitConfig
	defer func() {
		setting.GitConfig = oldGitConfig
	}()

	setting.GitConfig.Options["sync-test.cfg-key-a"] = "CfgValA"
	assert.NoError(t, syncGitConfig(t.Context()))
	assert.True(t, gitConfigContains("[sync-test]"))
	assert.True(t, gitConfigContains("cfg-key-a = CfgValA"))
}
