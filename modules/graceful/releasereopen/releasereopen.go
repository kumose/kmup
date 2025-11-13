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

package releasereopen

import (
	"errors"
	"sync"
)

type ReleaseReopener interface {
	ReleaseReopen() error
}

type Manager struct {
	mu      sync.Mutex
	counter int64

	releaseReopeners map[int64]ReleaseReopener
}

func (r *Manager) Register(rr ReleaseReopener) (cancel func()) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	currentCounter := r.counter
	r.releaseReopeners[r.counter] = rr

	return func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		delete(r.releaseReopeners, currentCounter)
	}
}

func (r *Manager) ReleaseReopen() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var errs []error
	for _, rr := range r.releaseReopeners {
		if err := rr.ReleaseReopen(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func GetManager() *Manager {
	return manager
}

func NewManager() *Manager {
	return &Manager{
		releaseReopeners: make(map[int64]ReleaseReopener),
	}
}

var manager = NewManager()
