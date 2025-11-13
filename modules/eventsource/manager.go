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

package eventsource

import (
	"sync"
)

// Manager manages the eventsource Messengers
type Manager struct {
	mutex sync.Mutex

	messengers map[int64]*Messenger
	connection chan struct{}
}

var manager *Manager

func init() {
	manager = &Manager{
		messengers: make(map[int64]*Messenger),
		connection: make(chan struct{}, 1),
	}
}

// GetManager returns a Manager and initializes one as singleton if there's none yet
func GetManager() *Manager {
	return manager
}

// Register message channel
func (m *Manager) Register(uid int64) <-chan *Event {
	m.mutex.Lock()
	messenger, ok := m.messengers[uid]
	if !ok {
		messenger = NewMessenger(uid)
		m.messengers[uid] = messenger
	}
	select {
	case m.connection <- struct{}{}:
	default:
	}
	m.mutex.Unlock()
	return messenger.Register()
}

// Unregister message channel
func (m *Manager) Unregister(uid int64, channel <-chan *Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	messenger, ok := m.messengers[uid]
	if !ok {
		return
	}
	if messenger.Unregister(channel) {
		delete(m.messengers, uid)
	}
}

// UnregisterAll message channels
func (m *Manager) UnregisterAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, messenger := range m.messengers {
		messenger.UnregisterAll()
	}
	m.messengers = map[int64]*Messenger{}
}

// SendMessage sends a message to a particular user
func (m *Manager) SendMessage(uid int64, message *Event) {
	m.mutex.Lock()
	messenger, ok := m.messengers[uid]
	m.mutex.Unlock()
	if ok {
		messenger.SendMessage(message)
	}
}

// SendMessageBlocking sends a message to a particular user
func (m *Manager) SendMessageBlocking(uid int64, message *Event) {
	m.mutex.Lock()
	messenger, ok := m.messengers[uid]
	m.mutex.Unlock()
	if ok {
		messenger.SendMessageBlocking(message)
	}
}
