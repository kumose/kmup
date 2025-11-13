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
	"context"
	"time"

	activities_model "github.com/kumose/kmup/models/activities"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/services/convert"
)

// Init starts this eventsource
func (m *Manager) Init() {
	if setting.UI.Notification.EventSourceUpdateTime <= 0 {
		return
	}
	go graceful.GetManager().RunWithShutdownContext(m.Run)
}

// Run runs the manager within a provided context
func (m *Manager) Run(ctx context.Context) {
	ctx, _, finished := process.GetManager().AddTypedContext(ctx, "Service: EventSource", process.SystemProcessType, true)
	defer finished()

	then := timeutil.TimeStampNow().Add(-2)
	timer := time.NewTicker(setting.UI.Notification.EventSourceUpdateTime)
loop:
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			break loop
		case <-timer.C:
			m.mutex.Lock()
			connectionCount := len(m.messengers)
			if connectionCount == 0 {
				log.Trace("Event source has no listeners")
				// empty the connection channel
				select {
				case <-m.connection:
				default:
				}
			}
			m.mutex.Unlock()
			if connectionCount == 0 {
				// No listeners so the source can be paused
				log.Trace("Pausing the eventsource")
				select {
				case <-ctx.Done():
					break loop
				case <-m.connection:
					log.Trace("Connection detected - restarting the eventsource")
					// OK we're back so lets reset the timer and start again
					// We won't change the "then" time because there could be concurrency issues
					select {
					case <-timer.C:
					default:
					}
					continue
				}
			}

			now := timeutil.TimeStampNow().Add(-2)

			uidCounts, err := activities_model.GetUIDsAndNotificationCounts(ctx, then, now)
			if err != nil {
				log.Error("Unable to get UIDcounts: %v", err)
			}
			for _, uidCount := range uidCounts {
				m.SendMessage(uidCount.UserID, &Event{
					Name: "notification-count",
					Data: uidCount,
				})
			}
			then = now

			if setting.Service.EnableTimetracking {
				usersStopwatches, err := issues_model.GetUIDsAndStopwatch(ctx)
				if err != nil {
					log.Error("Unable to get GetUIDsAndStopwatch: %v", err)
					return
				}

				for _, userStopwatches := range usersStopwatches {
					apiSWs, err := convert.ToStopWatches(ctx, userStopwatches.StopWatches)
					if err != nil {
						if !issues_model.IsErrIssueNotExist(err) {
							log.Error("Unable to APIFormat stopwatches: %v", err)
						}
						continue
					}
					dataBs, err := json.Marshal(apiSWs)
					if err != nil {
						log.Error("Unable to marshal stopwatches: %v", err)
						continue
					}
					m.SendMessage(userStopwatches.UserID, &Event{
						Name: "stopwatches",
						Data: string(dataBs),
					})
				}
			}
		}
	}
	m.UnregisterAll()
}
