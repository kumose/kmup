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

package nosql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToRedisURI(t *testing.T) {
	tests := []struct {
		name       string
		connection string
		want       string
	}{
		{
			name:       "old_default",
			connection: "addrs=127.0.0.1:6379 db=0",
			want:       "redis://127.0.0.1:6379/0",
		},
		{
			name:       "old_macaron_session_default",
			connection: "network=tcp,addr=127.0.0.1:6379,password=macaron,db=0,pool_size=100,idle_timeout=180",
			want:       "redis://:macaron@127.0.0.1:6379/0?idle_timeout=180s&pool_size=100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToRedisURI(tt.connection)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got.String())
		})
	}
}
