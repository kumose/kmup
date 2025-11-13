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

package storage

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_discardStorage(t *testing.T) {
	tests := []discardStorage{
		uninitializedStorage,
		discardStorage("empty"),
	}
	for _, tt := range tests {
		t.Run(string(tt), func(t *testing.T) {
			{
				got, err := tt.Open("path")
				assert.Nil(t, got)
				assert.Error(t, err, string(tt))
			}
			{
				got, err := tt.Save("path", bytes.NewReader([]byte{0}), 1)
				assert.Equal(t, int64(0), got)
				assert.Error(t, err, string(tt))
			}
			{
				got, err := tt.Stat("path")
				assert.Nil(t, got)
				assert.Error(t, err, string(tt))
			}
			{
				err := tt.Delete("path")
				assert.Error(t, err, string(tt))
			}
			{
				got, err := tt.URL("path", "name", "GET", nil)
				assert.Nil(t, got)
				assert.Errorf(t, err, string(tt))
			}
			{
				err := tt.IterateObjects("", func(_ string, _ Object) error { return nil })
				assert.Error(t, err, string(tt))
			}
		})
	}
}
