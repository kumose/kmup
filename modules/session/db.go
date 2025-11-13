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
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/timeutil"

	"github.com/kumose-go/chi/session"
)

// DBStore represents a session store implementation based on the DB.
type DBStore struct {
	sid  string
	lock sync.RWMutex
	data map[any]any
}

func dbContext() context.Context {
	return context.Background()
}

// NewDBStore creates and returns a DB session store.
func NewDBStore(sid string, kv map[any]any) *DBStore {
	return &DBStore{
		sid:  sid,
		data: kv,
	}
}

// Set sets value to given key in session.
func (s *DBStore) Set(key, val any) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = val
	return nil
}

// Get gets value by given key in session.
func (s *DBStore) Get(key any) any {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data[key]
}

// Delete delete a key from session.
func (s *DBStore) Delete(key any) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.data, key)
	return nil
}

// ID returns current session ID.
func (s *DBStore) ID() string {
	return s.sid
}

// Release releases resource and save data to provider.
func (s *DBStore) Release() error {
	// Skip encoding if the data is empty
	if len(s.data) == 0 {
		return nil
	}

	data, err := session.EncodeGob(s.data)
	if err != nil {
		return err
	}

	return auth.UpdateSession(dbContext(), s.sid, data)
}

// Flush deletes all session data.
func (s *DBStore) Flush() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data = make(map[any]any)
	return nil
}

// DBProvider represents a DB session provider implementation.
type DBProvider struct {
	maxLifetime int64
}

// Init initializes DB session provider.
// connStr: username:password@protocol(address)/dbname?param=value
func (p *DBProvider) Init(maxLifetime int64, connStr string) error {
	p.maxLifetime = maxLifetime
	return nil
}

// Read returns raw session store by session ID.
func (p *DBProvider) Read(sid string) (session.RawStore, error) {
	s, err := auth.ReadSession(dbContext(), sid)
	if err != nil {
		return nil, err
	}

	var kv map[any]any
	if len(s.Data) == 0 || s.Expiry.Add(p.maxLifetime) <= timeutil.TimeStampNow() {
		kv = make(map[any]any)
	} else {
		kv, err = session.DecodeGob(s.Data)
		if err != nil {
			return nil, err
		}
	}

	return NewDBStore(sid, kv), nil
}

// Exist returns true if session with given ID exists.
func (p *DBProvider) Exist(sid string) (bool, error) {
	has, err := auth.ExistSession(dbContext(), sid)
	if err != nil {
		return false, fmt.Errorf("session/DB: error checking existence: %w", err)
	}
	return has, nil
}

// Destroy deletes a session by session ID.
func (p *DBProvider) Destroy(sid string) error {
	return auth.DestroySession(dbContext(), sid)
}

// Regenerate regenerates a session store from old session ID to new one.
func (p *DBProvider) Regenerate(oldsid, sid string) (_ session.RawStore, err error) {
	s, err := auth.RegenerateSession(dbContext(), oldsid, sid)
	if err != nil {
		return nil, err
	}

	var kv map[any]any
	if len(s.Data) == 0 || s.Expiry.Add(p.maxLifetime) <= timeutil.TimeStampNow() {
		kv = make(map[any]any)
	} else {
		kv, err = session.DecodeGob(s.Data)
		if err != nil {
			return nil, err
		}
	}

	return NewDBStore(sid, kv), nil
}

// Count counts and returns number of sessions.
func (p *DBProvider) Count() (int, error) {
	total, err := auth.CountSessions(dbContext())
	if err != nil {
		return 0, fmt.Errorf("session/DB: error counting records: %w", err)
	}
	return int(total), nil
}

// GC calls GC to clean expired sessions.
func (p *DBProvider) GC() {
	if err := auth.CleanupSessions(dbContext(), p.maxLifetime); err != nil {
		log.Printf("session/DB: error garbage collecting: %v", err)
	}
}

func init() {
	session.Register("db", &DBProvider{})
}
