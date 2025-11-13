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

package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/nosql"

	"github.com/kumose-go/chi/session"
	"github.com/redis/go-redis/v9"
)

// RedisStore represents a redis session store implementation.
type RedisStore struct {
	c           redis.UniversalClient
	prefix, sid string
	duration    time.Duration
	lock        sync.RWMutex
	data        map[any]any
}

// NewRedisStore creates and returns a redis session store.
func NewRedisStore(c redis.UniversalClient, prefix, sid string, dur time.Duration, kv map[any]any) *RedisStore {
	return &RedisStore{
		c:        c,
		prefix:   prefix,
		sid:      sid,
		duration: dur,
		data:     kv,
	}
}

// Set sets value to given key in session.
func (s *RedisStore) Set(key, val any) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = val
	return nil
}

// Get gets value by given key in session.
func (s *RedisStore) Get(key any) any {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data[key]
}

// Delete delete a key from session.
func (s *RedisStore) Delete(key any) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.data, key)
	return nil
}

// ID returns current session ID.
func (s *RedisStore) ID() string {
	return s.sid
}

// Release releases resource and save data to provider.
func (s *RedisStore) Release() error {
	// Skip encoding if the data is empty
	if len(s.data) == 0 {
		return nil
	}

	data, err := session.EncodeGob(s.data)
	if err != nil {
		return err
	}

	return s.c.Set(graceful.GetManager().HammerContext(), s.prefix+s.sid, string(data), s.duration).Err()
}

// Flush deletes all session data.
func (s *RedisStore) Flush() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data = make(map[any]any)
	return nil
}

// RedisProvider represents a redis session provider implementation.
type RedisProvider struct {
	c        redis.UniversalClient
	duration time.Duration
	prefix   string
}

// Init initializes redis session provider.
// configs: network=tcp,addr=:6379,password=macaron,db=0,pool_size=100,idle_timeout=180,prefix=session;
func (p *RedisProvider) Init(maxlifetime int64, configs string) (err error) {
	p.duration, err = time.ParseDuration(fmt.Sprintf("%ds", maxlifetime))
	if err != nil {
		return err
	}

	uri := nosql.ToRedisURI(configs)

	for k, v := range uri.Query() {
		switch k {
		case "prefix":
			p.prefix = v[0]
		}
	}

	p.c = nosql.GetManager().GetRedisClient(uri.String())
	return p.c.Ping(graceful.GetManager().ShutdownContext()).Err()
}

// Read returns raw session store by session ID.
func (p *RedisProvider) Read(sid string) (session.RawStore, error) {
	psid := p.prefix + sid
	if exist, err := p.Exist(sid); err == nil && !exist {
		if err := p.c.Set(graceful.GetManager().HammerContext(), psid, "", p.duration).Err(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var kv map[any]any
	kvs, err := p.c.Get(graceful.GetManager().HammerContext(), psid).Result()
	if err != nil {
		return nil, err
	}
	if len(kvs) == 0 {
		kv = make(map[any]any)
	} else {
		kv, err = session.DecodeGob([]byte(kvs))
		if err != nil {
			return nil, err
		}
	}

	return NewRedisStore(p.c, p.prefix, sid, p.duration, kv), nil
}

// Exist returns true if session with given ID exists.
func (p *RedisProvider) Exist(sid string) (bool, error) {
	v, err := p.c.Exists(graceful.GetManager().HammerContext(), p.prefix+sid).Result()
	return err == nil && v == 1, err
}

// Destroy deletes a session by session ID.
func (p *RedisProvider) Destroy(sid string) error {
	return p.c.Del(graceful.GetManager().HammerContext(), p.prefix+sid).Err()
}

// Regenerate regenerates a session store from old session ID to new one.
func (p *RedisProvider) Regenerate(oldsid, sid string) (_ session.RawStore, err error) {
	poldsid := p.prefix + oldsid
	psid := p.prefix + sid

	if exist, err := p.Exist(sid); err != nil {
		return nil, err
	} else if exist {
		return nil, fmt.Errorf("new sid '%s' already exists", sid)
	}
	if exist, err := p.Exist(oldsid); err == nil && !exist {
		// Make a fake old session.
		if err := p.c.Set(graceful.GetManager().HammerContext(), poldsid, "", p.duration).Err(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// do not use Rename here, because the old sid and new sid may be in different redis cluster slot.
	kvs, err := p.c.Get(graceful.GetManager().HammerContext(), poldsid).Result()
	if err != nil {
		return nil, err
	}

	if err = p.c.Del(graceful.GetManager().HammerContext(), poldsid).Err(); err != nil {
		return nil, err
	}

	if err = p.c.Set(graceful.GetManager().HammerContext(), psid, kvs, p.duration).Err(); err != nil {
		return nil, err
	}

	var kv map[any]any
	if len(kvs) == 0 {
		kv = make(map[any]any)
	} else {
		kv, err = session.DecodeGob([]byte(kvs))
		if err != nil {
			return nil, err
		}
	}

	return NewRedisStore(p.c, p.prefix, sid, p.duration, kv), nil
}

// Count counts and returns number of sessions.
func (p *RedisProvider) Count() (int, error) {
	size, err := p.c.DBSize(graceful.GetManager().HammerContext()).Result()
	return int(size), err
}

// GC calls GC to clean expired sessions.
func (*RedisProvider) GC() {}

func init() {
	session.Register("redis", &RedisProvider{})
}
