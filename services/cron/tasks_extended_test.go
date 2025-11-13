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

package cron

import (
	"testing"
	"time"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func Test_GCLFSConfig(t *testing.T) {
	cfg, err := setting.NewConfigProviderFromData(`
[cron.gc_lfs]
ENABLED = true
RUN_AT_START = true
SCHEDULE = "@every 2h"
OLDER_THAN = "1h"
LAST_UPDATED_MORE_THAN_AGO = "7h"
NUMBER_TO_CHECK_PER_REPO = 10
PROPORTION_TO_CHECK_PER_REPO = 0.1
`)
	assert.NoError(t, err)
	defer test.MockVariableValue(&setting.CfgProvider, cfg)()

	config := &GCLFSConfig{
		BaseConfig: BaseConfig{
			Enabled:    false,
			RunAtStart: false,
			Schedule:   "@every 24h",
		},
		OlderThan:                24 * time.Hour * 7,
		LastUpdatedMoreThanAgo:   24 * time.Hour * 3,
		NumberToCheckPerRepo:     100,
		ProportionToCheckPerRepo: 0.6,
	}

	_, err = setting.GetCronSettings("gc_lfs", config)
	assert.NoError(t, err)
	assert.True(t, config.Enabled)
	assert.True(t, config.RunAtStart)
	assert.Equal(t, "@every 2h", config.Schedule)
	assert.Equal(t, 1*time.Hour, config.OlderThan)
	assert.Equal(t, 7*time.Hour, config.LastUpdatedMoreThanAgo)
	assert.Equal(t, int64(10), config.NumberToCheckPerRepo)
	assert.InDelta(t, 0.1, config.ProportionToCheckPerRepo, 0.001)
}
