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

package log

import (
	"fmt"
	"testing"

	"github.com/kumose/kmup/modules/json"

	"github.com/stretchr/testify/assert"
)

type testLevel struct {
	Level Level `json:"level"`
}

func TestLevelMarshalUnmarshalJSON(t *testing.T) {
	levelBytes, err := json.Marshal(testLevel{
		Level: INFO,
	})
	assert.NoError(t, err)
	assert.Equal(t, string(makeTestLevelBytes(INFO.String())), string(levelBytes))

	var testLevel testLevel
	err = json.Unmarshal(levelBytes, &testLevel)
	assert.NoError(t, err)
	assert.Equal(t, INFO, testLevel.Level)

	err = json.Unmarshal(makeTestLevelBytes(`FOFOO`), &testLevel)
	assert.NoError(t, err)
	assert.Equal(t, INFO, testLevel.Level)

	err = json.Unmarshal(fmt.Appendf(nil, `{"level":%d}`, 2), &testLevel)
	assert.NoError(t, err)
	assert.Equal(t, INFO, testLevel.Level)

	err = json.Unmarshal(fmt.Appendf(nil, `{"level":%d}`, 10012), &testLevel)
	assert.NoError(t, err)
	assert.Equal(t, INFO, testLevel.Level)

	err = json.Unmarshal([]byte(`{"level":{}}`), &testLevel)
	assert.NoError(t, err)
	assert.Equal(t, INFO, testLevel.Level)

	assert.Equal(t, INFO.String(), Level(1001).String())

	err = json.Unmarshal([]byte(`{"level":{}`), &testLevel.Level)
	assert.Error(t, err)
}

func makeTestLevelBytes(level string) []byte {
	return fmt.Appendf(nil, `{"level":"%s"}`, level)
}
