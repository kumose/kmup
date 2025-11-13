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
	"bytes"
	"encoding/gob"
	"net/http"

	"github.com/kumose-go/chi/session"
)

type mockMemRawStore struct {
	s *session.MemStore
}

var _ session.RawStore = (*mockMemRawStore)(nil)

func (m *mockMemRawStore) Set(k, v any) error {
	// We need to use gob to encode the value, to make it have the same behavior as other stores and catch abuses.
	// Because gob needs to "Register" the type before it can encode it, and it's unable to decode a struct to "any" so use a map to help to decode the value.
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(map[string]any{"v": v}); err != nil {
		return err
	}
	return m.s.Set(k, buf.Bytes())
}

func (m *mockMemRawStore) Get(k any) (ret any) {
	v, ok := m.s.Get(k).([]byte)
	if !ok {
		return nil
	}
	var w map[string]any
	_ = gob.NewDecoder(bytes.NewBuffer(v)).Decode(&w)
	return w["v"]
}

func (m *mockMemRawStore) Delete(k any) error {
	return m.s.Delete(k)
}

func (m *mockMemRawStore) ID() string {
	return m.s.ID()
}

func (m *mockMemRawStore) Release() error {
	return m.s.Release()
}

func (m *mockMemRawStore) Flush() error {
	return m.s.Flush()
}

type mockMemStore struct {
	*mockMemRawStore
}

var _ Store = (*mockMemStore)(nil)

func (m mockMemStore) Destroy(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func NewMockMemStore(sid string) Store {
	return &mockMemStore{&mockMemRawStore{session.NewMemStore(sid)}}
}
