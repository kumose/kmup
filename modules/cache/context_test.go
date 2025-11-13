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

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestWithCacheContext(t *testing.T) {
	ctx := WithCacheContext(t.Context())
	c := GetContextCache(ctx)
	v, _ := c.Get("empty_field", "my_config1")
	assert.Nil(t, v)

	const field = "system_setting"
	v, _ = c.Get(field, "my_config1")
	assert.Nil(t, v)
	c.Put(field, "my_config1", 1)
	v, _ = c.Get(field, "my_config1")
	assert.NotNil(t, v)
	assert.Equal(t, 1, v.(int))

	c.Delete(field, "my_config1")
	c.Delete(field, "my_config2") // remove a non-exist key

	v, _ = c.Get(field, "my_config1")
	assert.Nil(t, v)

	vInt, err := GetWithContextCache(ctx, field, "my_config1", func(context.Context, string) (int, error) {
		return 1, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, vInt)

	v, _ = c.Get(field, "my_config1")
	assert.EqualValues(t, 1, v)

	defer test.MockVariableValue(&timeNow, func() time.Time {
		return time.Now().Add(5 * time.Minute)
	})()
	v, _ = c.Get(field, "my_config1")
	assert.Nil(t, v)
}
