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
	"sync"
	"time"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/util"
)

// EphemeralCache is a cache that can be used to store data in a request level context
// This is useful for caching data that is expensive to calculate and is likely to be
// used multiple times in a request.
type EphemeralCache struct {
	data          map[any]map[any]any
	lock          sync.RWMutex
	created       time.Time
	checkLifeTime time.Duration
}

var timeNow = time.Now

func NewEphemeralCache(checkLifeTime ...time.Duration) *EphemeralCache {
	return &EphemeralCache{
		data:          make(map[any]map[any]any),
		created:       timeNow(),
		checkLifeTime: util.OptionalArg(checkLifeTime, 0),
	}
}

func (cc *EphemeralCache) checkExceededLifeTime(tp, key any) bool {
	if cc.checkLifeTime > 0 && timeNow().Sub(cc.created) > cc.checkLifeTime {
		log.Warn("EphemeralCache is expired, is highly likely to be abused for long-life tasks: %v, %v", tp, key)
		return true
	}
	return false
}

func (cc *EphemeralCache) Get(tp, key any) (any, bool) {
	if cc.checkExceededLifeTime(tp, key) {
		return nil, false
	}
	cc.lock.RLock()
	defer cc.lock.RUnlock()
	ret, ok := cc.data[tp][key]
	return ret, ok
}

func (cc *EphemeralCache) Put(tp, key, value any) {
	if cc.checkExceededLifeTime(tp, key) {
		return
	}

	cc.lock.Lock()
	defer cc.lock.Unlock()

	d := cc.data[tp]
	if d == nil {
		d = make(map[any]any)
		cc.data[tp] = d
	}
	d[key] = value
}

func (cc *EphemeralCache) Delete(tp, key any) {
	if cc.checkExceededLifeTime(tp, key) {
		return
	}

	cc.lock.Lock()
	defer cc.lock.Unlock()
	delete(cc.data[tp], key)
}

func GetWithEphemeralCache[T, K any](ctx context.Context, c *EphemeralCache, groupKey string, targetKey K, f func(context.Context, K) (T, error)) (T, error) {
	v, has := c.Get(groupKey, targetKey)
	if vv, ok := v.(T); has && ok {
		return vv, nil
	}
	t, err := f(ctx, targetKey)
	if err != nil {
		return t, err
	}
	c.Put(groupKey, targetKey, t)
	return t, nil
}
