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

package openid

import (
	"sync"
	"time"

	"github.com/yohcop/openid-go"
)

type timedDiscoveredInfo struct {
	info openid.DiscoveredInfo
	time time.Time
}

type timedDiscoveryCache struct {
	cache map[string]timedDiscoveredInfo
	ttl   time.Duration
	mutex *sync.Mutex
}

func newTimedDiscoveryCache(ttl time.Duration) *timedDiscoveryCache {
	return &timedDiscoveryCache{cache: map[string]timedDiscoveredInfo{}, ttl: ttl, mutex: &sync.Mutex{}}
}

func (s *timedDiscoveryCache) Put(id string, info openid.DiscoveredInfo) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache[id] = timedDiscoveredInfo{info: info, time: time.Now()}
}

// Delete timed-out cache entries
func (s *timedDiscoveryCache) cleanTimedOut() {
	now := time.Now()
	for k, e := range s.cache {
		diff := now.Sub(e.time)
		if diff > s.ttl {
			delete(s.cache, k)
		}
	}
}

func (s *timedDiscoveryCache) Get(id string) openid.DiscoveredInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Delete old cached while we are at it.
	s.cleanTimedOut()

	if info, has := s.cache[id]; has {
		return info.info
	}
	return nil
}
