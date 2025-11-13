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

package actions

import (
	"testing"
	"time"

	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionScheduleSpec_Parse(t *testing.T) {
	// Mock the local timezone is not UTC
	tz, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)
	defer test.MockVariableValue(&time.Local, tz)()

	now, err := time.Parse(time.RFC3339, "2024-07-31T15:47:55+08:00")
	require.NoError(t, err)

	tests := []struct {
		name    string
		spec    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "regular",
			spec:    "0 10 * * *",
			want:    "2024-07-31T10:00:00Z",
			wantErr: assert.NoError,
		},
		{
			name:    "invalid",
			spec:    "0 10 * *",
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "with timezone",
			spec:    "TZ=America/New_York 0 10 * * *",
			want:    "2024-07-31T14:00:00Z",
			wantErr: assert.NoError,
		},
		{
			name:    "timezone irrelevant",
			spec:    "@every 5m",
			want:    "2024-07-31T07:52:55Z",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ActionScheduleSpec{
				Spec: tt.spec,
			}
			got, err := s.Parse()
			tt.wantErr(t, err)

			if err == nil {
				assert.Equal(t, tt.want, got.Next(now).UTC().Format(time.RFC3339))
			}
		})
	}
}
