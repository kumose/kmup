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
	"errors"
	"strings"

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"

	chi_cache "github.com/kumose-go/chi/cache" //nolint:depguard // we wrap this package here
)

type GetJSONError struct {
	err         error
	cachedError string // Golang error can't be stored in cache, only the string message could be stored
}

func (e *GetJSONError) ToError() error {
	if e.err != nil {
		return e.err
	}
	return errors.New("cached error: " + e.cachedError)
}

type StringCache interface {
	Ping() error

	Get(key string) (string, bool)
	Put(key, value string, ttl int64) error
	Delete(key string) error
	IsExist(key string) bool

	PutJSON(key string, v any, ttl int64) error
	GetJSON(key string, ptr any) (exist bool, err *GetJSONError)

	ChiCache() chi_cache.Cache
}

type stringCache struct {
	chiCache chi_cache.Cache
}

func NewStringCache(cacheConfig setting.Cache) (StringCache, error) {
	adapter := util.IfZero(cacheConfig.Adapter, "memory")
	interval := util.IfZero(cacheConfig.Interval, 60)
	cc, err := chi_cache.NewCacher(chi_cache.Options{
		Adapter:       adapter,
		AdapterConfig: cacheConfig.Conn,
		Interval:      interval,
	})
	if err != nil {
		return nil, err
	}
	return &stringCache{chiCache: cc}, nil
}

func (sc *stringCache) Ping() error {
	return sc.chiCache.Ping()
}

func (sc *stringCache) Get(key string) (string, bool) {
	v := sc.chiCache.Get(key)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func (sc *stringCache) Put(key, value string, ttl int64) error {
	return sc.chiCache.Put(key, value, ttl)
}

func (sc *stringCache) Delete(key string) error {
	return sc.chiCache.Delete(key)
}

func (sc *stringCache) IsExist(key string) bool {
	return sc.chiCache.IsExist(key)
}

const cachedErrorPrefix = "<CACHED-ERROR>:"

func (sc *stringCache) PutJSON(key string, v any, ttl int64) error {
	var s string
	switch v := v.(type) {
	case error:
		s = cachedErrorPrefix + v.Error()
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		s = util.UnsafeBytesToString(b)
	}
	return sc.chiCache.Put(key, s, ttl)
}

func (sc *stringCache) GetJSON(key string, ptr any) (exist bool, getErr *GetJSONError) {
	s, ok := sc.Get(key)
	if !ok || s == "" {
		return false, nil
	}
	s, isCachedError := strings.CutPrefix(s, cachedErrorPrefix)
	if isCachedError {
		return true, &GetJSONError{cachedError: s}
	}
	if err := json.Unmarshal(util.UnsafeStringToBytes(s), ptr); err != nil {
		return false, &GetJSONError{err: err}
	}
	return true, nil
}

func (sc *stringCache) ChiCache() chi_cache.Cache {
	return sc.chiCache
}
