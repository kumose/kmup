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

import "sync"

// Messenger is a per uid message store
type Messenger struct {
	mutex    sync.Mutex
	uid      int64
	channels []chan *Event
}

// NewMessenger creates a messenger for a particular uid
func NewMessenger(uid int64) *Messenger {
	return &Messenger{
		uid:      uid,
		channels: [](chan *Event){},
	}
}

// Register returns a new chan []byte
func (m *Messenger) Register() <-chan *Event {
	m.mutex.Lock()
	// TODO: Limit the number of messengers per uid
	channel := make(chan *Event, 1)
	m.channels = append(m.channels, channel)
	m.mutex.Unlock()
	return channel
}

// Unregister removes the provider chan []byte
func (m *Messenger) Unregister(channel <-chan *Event) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for i, toRemove := range m.channels {
		if channel == toRemove {
			m.channels = append(m.channels[:i], m.channels[i+1:]...)
			close(toRemove)
			break
		}
	}
	return len(m.channels) == 0
}

// UnregisterAll removes all chan []byte
func (m *Messenger) UnregisterAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, channel := range m.channels {
		close(channel)
	}
	m.channels = nil
}

// SendMessage sends the message to all registered channels
func (m *Messenger) SendMessage(message *Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for i := range m.channels {
		channel := m.channels[i]
		select {
		case channel <- message:
		default:
		}
	}
}

// SendMessageBlocking sends the message to all registered channels and ensures it gets sent
func (m *Messenger) SendMessageBlocking(message *Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for i := range m.channels {
		m.channels[i] <- message
	}
}
