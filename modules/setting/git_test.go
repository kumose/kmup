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

package setting

import (
	"testing"

	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestGitConfig(t *testing.T) {
	oldGit := Git
	oldGitConfig := GitConfig
	defer func() {
		Git = oldGit
		GitConfig = oldGitConfig
	}()

	cfg, err := NewConfigProviderFromData(`
[git.config]
a.b = 1
`)
	assert.NoError(t, err)
	loadGitFrom(cfg)
	assert.Equal(t, "1", GitConfig.Options["a.b"])
	assert.Equal(t, "histogram", GitConfig.Options["diff.algorithm"])

	cfg, err = NewConfigProviderFromData(`
[git.config]
diff.algorithm = other
`)
	assert.NoError(t, err)
	loadGitFrom(cfg)
	assert.Equal(t, "other", GitConfig.Options["diff.algorithm"])
}

func TestGitReflog(t *testing.T) {
	defer test.MockVariableValue(&Git)
	defer test.MockVariableValue(&GitConfig)

	// default reflog config without legacy options
	cfg, err := NewConfigProviderFromData(``)
	assert.NoError(t, err)
	loadGitFrom(cfg)

	assert.Equal(t, "true", GitConfig.GetOption("core.logAllRefUpdates"))
	assert.Equal(t, "90", GitConfig.GetOption("gc.reflogExpire"))

	// custom reflog config by legacy options
	cfg, err = NewConfigProviderFromData(`
[git.reflog]
ENABLED = false
EXPIRATION = 123
`)
	assert.NoError(t, err)
	loadGitFrom(cfg)

	assert.Equal(t, "false", GitConfig.GetOption("core.logAllRefUpdates"))
	assert.Equal(t, "123", GitConfig.GetOption("gc.reflogExpire"))
}
