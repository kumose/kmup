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
	"sync"

	"github.com/kumose/kmup/modules/setting"
)

var (
	defaultLocker Locker
	initOnce      sync.Once
	initFunc      = func() {
		switch setting.GlobalLock.ServiceType {
		case "redis":
			defaultLocker = NewRedisLocker(setting.GlobalLock.ServiceConnStr)
		case "memory":
			fallthrough
		default:
			defaultLocker = NewMemoryLocker()
		}
	} // define initFunc as a variable to make it possible to change it in tests
)

// DefaultLocker returns the default locker.
func DefaultLocker() Locker {
	initOnce.Do(func() {
		initFunc()
	})
	return defaultLocker
}

// Lock tries to acquire a lock for the given key, it uses the default locker.
// Read the documentation of Locker.Lock for more information about the behavior.
func Lock(ctx context.Context, key string) (ReleaseFunc, error) {
	return DefaultLocker().Lock(ctx, key)
}

// TryLock tries to acquire a lock for the given key, it uses the default locker.
// Read the documentation of Locker.TryLock for more information about the behavior.
func TryLock(ctx context.Context, key string) (bool, ReleaseFunc, error) {
	return DefaultLocker().TryLock(ctx, key)
}

// LockAndDo tries to acquire a lock for the given key and then calls the given function.
// It uses the default locker, and it will return an error if failed to acquire the lock.
func LockAndDo(ctx context.Context, key string, f func(context.Context) error) error {
	release, err := Lock(ctx, key)
	if err != nil {
		return err
	}
	defer release()

	return f(ctx)
}

// TryLockAndDo tries to acquire a lock for the given key and then calls the given function.
// It uses the default locker, and it will return false if failed to acquire the lock.
func TryLockAndDo(ctx context.Context, key string, f func(context.Context) error) (bool, error) {
	ok, release, err := TryLock(ctx, key)
	if err != nil {
		return false, err
	}
	defer release()

	if !ok {
		return false, nil
	}

	return true, f(ctx)
}
