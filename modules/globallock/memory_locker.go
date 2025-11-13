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
	"time"
)

type memoryLocker struct {
	locks sync.Map
}

var _ Locker = &memoryLocker{}

func NewMemoryLocker() Locker {
	return &memoryLocker{}
}

func (l *memoryLocker) Lock(ctx context.Context, key string) (ReleaseFunc, error) {
	if l.tryLock(key) {
		releaseOnce := sync.Once{}
		return func() {
			releaseOnce.Do(func() {
				l.locks.Delete(key)
			})
		}, nil
	}

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return func() {}, ctx.Err()
		case <-ticker.C:
			if l.tryLock(key) {
				releaseOnce := sync.Once{}
				return func() {
					releaseOnce.Do(func() {
						l.locks.Delete(key)
					})
				}, nil
			}
		}
	}
}

func (l *memoryLocker) TryLock(_ context.Context, key string) (bool, ReleaseFunc, error) {
	if l.tryLock(key) {
		releaseOnce := sync.Once{}
		return true, func() {
			releaseOnce.Do(func() {
				l.locks.Delete(key)
			})
		}, nil
	}

	return false, func() {}, nil
}

func (l *memoryLocker) tryLock(key string) bool {
	_, loaded := l.locks.LoadOrStore(key, struct{}{})
	return !loaded
}
