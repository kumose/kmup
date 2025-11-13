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

package test

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kumose/kmup/modules/log"
)

type LogChecker struct {
	*log.EventWriterBaseImpl

	filterMessages []string
	filtered       []bool

	stopMark string
	stopped  bool

	mu sync.Mutex
}

func (lc *LogChecker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-lc.Queue:
			if !ok {
				return
			}
			lc.checkLogEvent(event)
		}
	}
}

func (lc *LogChecker) checkLogEvent(event *log.EventFormatted) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	for i, msg := range lc.filterMessages {
		if strings.Contains(event.Origin.MsgSimpleText, msg) {
			lc.filtered[i] = true
		}
	}
	if strings.Contains(event.Origin.MsgSimpleText, lc.stopMark) {
		lc.stopped = true
	}
}

var checkerIndex int64

func NewLogChecker(namePrefix string) (logChecker *LogChecker, cancel func()) {
	logger := log.GetManager().GetLogger(namePrefix)
	newCheckerIndex := atomic.AddInt64(&checkerIndex, 1)
	writerName := namePrefix + "-" + strconv.FormatInt(newCheckerIndex, 10)

	lc := &LogChecker{}
	lc.EventWriterBaseImpl = log.NewEventWriterBase(writerName, "test-log-checker", log.WriterMode{})
	logger.AddWriters(lc)
	return lc, func() { _ = logger.RemoveWriter(writerName) }
}

// Filter will make the `Check` function to check if these logs are outputted.
func (lc *LogChecker) Filter(msgs ...string) *LogChecker {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.filterMessages = make([]string, len(msgs))
	copy(lc.filterMessages, msgs)
	lc.filtered = make([]bool, len(lc.filterMessages))
	return lc
}

func (lc *LogChecker) StopMark(msg string) *LogChecker {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.stopMark = msg
	lc.stopped = false
	return lc
}

// Check returns the filtered slice and whether the stop mark is reached.
func (lc *LogChecker) Check(d time.Duration) (filtered []bool, stopped bool) {
	stop := time.Now().Add(d)

	for {
		lc.mu.Lock()
		stopped = lc.stopped
		lc.mu.Unlock()

		if time.Now().After(stop) || stopped {
			lc.mu.Lock()
			f := make([]bool, len(lc.filtered))
			copy(f, lc.filtered)
			lc.mu.Unlock()
			return f, stopped
		}
		time.Sleep(10 * time.Millisecond)
	}
}
