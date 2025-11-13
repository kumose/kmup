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

package queue

import (
	"fmt"
	"sync"
)

// testStateRecorder is used to record state changes for testing, to help debug async behaviors
type testStateRecorder struct {
	records []string
	mu      sync.Mutex
}

var testRecorder = &testStateRecorder{}

func (t *testStateRecorder) Record(format string, args ...any) {
	t.mu.Lock()
	t.records = append(t.records, fmt.Sprintf(format, args...))
	if len(t.records) > 1000 {
		t.records = t.records[len(t.records)-1000:]
	}
	t.mu.Unlock()
}

func (t *testStateRecorder) Records() []string {
	t.mu.Lock()
	r := make([]string, len(t.records))
	copy(r, t.records)
	t.mu.Unlock()
	return r
}

func (t *testStateRecorder) Reset() {
	t.mu.Lock()
	t.records = nil
	t.mu.Unlock()
}
