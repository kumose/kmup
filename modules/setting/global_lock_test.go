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

	"github.com/stretchr/testify/assert"
)

func TestLoadGlobalLockConfig(t *testing.T) {
	t.Run("DefaultGlobalLockConfig", func(t *testing.T) {
		iniStr := ``
		cfg, err := NewConfigProviderFromData(iniStr)
		assert.NoError(t, err)

		loadGlobalLockFrom(cfg)
		assert.Equal(t, "memory", GlobalLock.ServiceType)
	})

	t.Run("RedisGlobalLockConfig", func(t *testing.T) {
		iniStr := `
[global_lock]
SERVICE_TYPE = redis
SERVICE_CONN_STR = addrs=127.0.0.1:6379 db=0
`
		cfg, err := NewConfigProviderFromData(iniStr)
		assert.NoError(t, err)

		loadGlobalLockFrom(cfg)
		assert.Equal(t, "redis", GlobalLock.ServiceType)
		assert.Equal(t, "addrs=127.0.0.1:6379 db=0", GlobalLock.ServiceConnStr)
	})
}
