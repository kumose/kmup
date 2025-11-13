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

package admin

import (
	"runtime/pprof"
	"sync"
	"time"

	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/queue"
	"github.com/kumose/kmup/modules/setting"
)

var testQueueOnce sync.Once

// initTestQueueOnce initializes the test queue for dev mode
// the test queue will also be shown in the queue list
// developers could see the queue length / worker number / items number on the admin page and try to remove the items
func initTestQueueOnce() {
	testQueueOnce.Do(func() {
		ctx, _, finished := process.GetManager().AddTypedContext(graceful.GetManager().ShutdownContext(), "TestQueue", process.SystemProcessType, false)
		qs := setting.QueueSettings{
			Name:        "test-queue",
			Type:        "channel",
			Length:      20,
			BatchLength: 2,
			MaxWorkers:  3,
		}
		testQueue, err := queue.NewWorkerPoolQueueWithContext(ctx, "test-queue", qs, func(t ...int64) (unhandled []int64) {
			for range t {
				select {
				case <-graceful.GetManager().ShutdownContext().Done():
				case <-time.After(5 * time.Second):
				}
			}
			return nil
		}, true)
		if err != nil {
			log.Error("unable to create test queue: %v", err)
			return
		}

		queue.GetManager().AddManagedQueue(testQueue)
		testQueue.SetWorkerMaxNumber(5)
		go graceful.GetManager().RunWithCancel(testQueue)
		go func() {
			pprof.SetGoroutineLabels(ctx)
			defer finished()

			cnt := int64(0)
			adding := true
			for {
				select {
				case <-ctx.Done():
				case <-time.After(500 * time.Millisecond):
					if adding {
						if testQueue.GetQueueItemNumber() == qs.Length {
							adding = false
						}
					} else {
						if testQueue.GetQueueItemNumber() == 0 {
							adding = true
						}
					}
					if adding {
						_ = testQueue.Push(cnt)
						cnt++
					}
				}
			}
		}()
	})
}
