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

package setting

import (
	"strings"
	"time"

	"github.com/kumose/kmup/modules/log"
)

// Cache represents cache settings
type Cache struct {
	Adapter  string
	Interval int
	Conn     string
	TTL      time.Duration `ini:"ITEM_TTL"`
}

// CacheService the global cache
var CacheService = struct {
	Cache `ini:"cache"`

	LastCommit struct {
		TTL          time.Duration `ini:"ITEM_TTL"`
		CommitsCount int64
	} `ini:"cache.last_commit"`
}{
	Cache: Cache{
		Adapter:  "memory",
		Interval: 60,
		TTL:      16 * time.Hour,
	},
	LastCommit: struct {
		TTL          time.Duration `ini:"ITEM_TTL"`
		CommitsCount int64
	}{
		TTL:          8760 * time.Hour,
		CommitsCount: 1000,
	},
}

// MemcacheMaxTTL represents the maximum memcache TTL
const MemcacheMaxTTL = 30 * 24 * time.Hour

func loadCacheFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("cache")
	if err := sec.MapTo(&CacheService); err != nil {
		log.Fatal("Failed to map Cache settings: %v", err)
	}

	CacheService.Adapter = sec.Key("ADAPTER").In("memory", []string{"memory", "redis", "memcache", "twoqueue"})
	switch CacheService.Adapter {
	case "memory":
	case "redis", "memcache":
		CacheService.Conn = strings.Trim(sec.Key("HOST").String(), "\" ")
	case "twoqueue":
		CacheService.Conn = strings.TrimSpace(sec.Key("HOST").String())
		if CacheService.Conn == "" {
			CacheService.Conn = "50000"
		}
	default:
		log.Fatal("Unknown cache adapter: %s", CacheService.Adapter)
	}

	sec = rootCfg.Section("cache.last_commit")
	CacheService.LastCommit.CommitsCount = sec.Key("COMMITS_COUNT").MustInt64(1000)
}

// TTLSeconds returns the TTLSeconds or unix timestamp for memcache
func (c Cache) TTLSeconds() int64 {
	if c.Adapter == "memcache" && c.TTL > MemcacheMaxTTL {
		return time.Now().Add(c.TTL).Unix()
	}
	return int64(c.TTL.Seconds())
}

// LastCommitCacheTTLSeconds returns the TTLSeconds or unix timestamp for memcache
func LastCommitCacheTTLSeconds() int64 {
	if CacheService.Adapter == "memcache" && CacheService.LastCommit.TTL > MemcacheMaxTTL {
		return time.Now().Add(CacheService.LastCommit.TTL).Unix()
	}
	return int64(CacheService.LastCommit.TTL.Seconds())
}
