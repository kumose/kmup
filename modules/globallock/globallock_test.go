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

package globallock

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLockAndDo(t *testing.T) {
	t.Run("redis", func(t *testing.T) {
		url := "redis://127.0.0.1:6379/0"
		if os.Getenv("CI") == "" {
			// Make it possible to run tests against a local redis instance
			url = os.Getenv("TEST_REDIS_URL")
			if url == "" {
				t.Skip("TEST_REDIS_URL not set and not running in CI")
				return
			}
		}

		oldDefaultLocker := defaultLocker
		oldInitFunc := initFunc
		defer func() {
			defaultLocker = oldDefaultLocker
			initFunc = oldInitFunc
			if defaultLocker == nil {
				initOnce = sync.Once{}
			}
		}()

		initOnce = sync.Once{}
		initFunc = func() {
			defaultLocker = NewRedisLocker(url)
		}

		testLockAndDo(t)
		require.NoError(t, defaultLocker.(*redisLocker).Close())
	})
	t.Run("memory", func(t *testing.T) {
		oldDefaultLocker := defaultLocker
		oldInitFunc := initFunc
		defer func() {
			defaultLocker = oldDefaultLocker
			initFunc = oldInitFunc
			if defaultLocker == nil {
				initOnce = sync.Once{}
			}
		}()

		initOnce = sync.Once{}
		initFunc = func() {
			defaultLocker = NewMemoryLocker()
		}

		testLockAndDo(t)
	})
}

func testLockAndDo(t *testing.T) {
	const concurrency = 50

	ctx := t.Context()
	count := 0
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for range concurrency {
		go func() {
			defer wg.Done()
			err := LockAndDo(ctx, "test", func(ctx context.Context) error {
				count++

				// It's impossible to acquire the lock inner the function
				ok, err := TryLockAndDo(ctx, "test", func(ctx context.Context) error {
					assert.Fail(t, "should not acquire the lock")
					return nil
				})
				assert.False(t, ok)
				assert.NoError(t, err)

				return nil
			})
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	assert.Equal(t, concurrency, count)
}
