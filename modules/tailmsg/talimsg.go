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

package tailmsg

import (
	"sync"
	"time"
)

type MsgRecord struct {
	Time    time.Time
	Content string
}

type MsgRecorder interface {
	Record(content string)
	GetRecords() []*MsgRecord
}

type memoryMsgRecorder struct {
	mu    sync.RWMutex
	msgs  []*MsgRecord
	limit int
}

// TODO: use redis for a clustered environment

func (m *memoryMsgRecorder) Record(content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.msgs = append(m.msgs, &MsgRecord{
		Time:    time.Now(),
		Content: content,
	})
	if len(m.msgs) > m.limit {
		m.msgs = m.msgs[len(m.msgs)-m.limit:]
	}
}

func (m *memoryMsgRecorder) GetRecords() []*MsgRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := make([]*MsgRecord, len(m.msgs))
	copy(ret, m.msgs)
	return ret
}

func NewMsgRecorder(limit int) MsgRecorder {
	return &memoryMsgRecorder{
		limit: limit,
	}
}

type Manager struct {
	traceRecorder MsgRecorder
	logRecorder   MsgRecorder
}

func (m *Manager) GetTraceRecorder() MsgRecorder {
	return m.traceRecorder
}

func (m *Manager) GetLogRecorder() MsgRecorder {
	return m.logRecorder
}

var GetManager = sync.OnceValue(func() *Manager {
	return &Manager{
		traceRecorder: NewMsgRecorder(100),
		logRecorder:   NewMsgRecorder(1000),
	}
})
