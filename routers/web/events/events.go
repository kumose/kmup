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

package events

import (
	"net/http"
	"time"

	"github.com/kumose/kmup/modules/eventsource"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/routers/web/auth"
	"github.com/kumose/kmup/services/context"
)

// Events listens for events
func Events(ctx *context.Context) {
	// FIXME: Need to check if resp is actually a http.Flusher! - how though?

	// Set the headers related to event streaming.
	ctx.Resp.Header().Set("Content-Type", "text/event-stream")
	ctx.Resp.Header().Set("Cache-Control", "no-cache")
	ctx.Resp.Header().Set("Connection", "keep-alive")
	ctx.Resp.Header().Set("X-Accel-Buffering", "no")
	ctx.Resp.WriteHeader(http.StatusOK)

	if !ctx.IsSigned {
		// Return unauthorized status event
		event := &eventsource.Event{
			Name: "close",
			Data: "unauthorized",
		}
		_, _ = event.WriteTo(ctx)
		ctx.Resp.Flush()
		return
	}

	// Listen to connection close and un-register messageChan
	notify := ctx.Done()
	ctx.Resp.Flush()

	shutdownCtx := graceful.GetManager().ShutdownContext()

	uid := ctx.Doer.ID

	messageChan := eventsource.GetManager().Register(uid)

	unregister := func() {
		eventsource.GetManager().Unregister(uid, messageChan)
		// ensure the messageChan is closed
		for {
			_, ok := <-messageChan
			if !ok {
				break
			}
		}
	}

	if _, err := ctx.Resp.Write([]byte("\n")); err != nil {
		log.Error("Unable to write to EventStream: %v", err)
		unregister()
		return
	}

	timer := time.NewTicker(30 * time.Second)

loop:
	for {
		select {
		case <-timer.C:
			event := &eventsource.Event{
				Name: "ping",
			}
			_, err := event.WriteTo(ctx.Resp)
			if err != nil {
				log.Error("Unable to write to EventStream for user %s: %v", ctx.Doer.Name, err)
				go unregister()
				break loop
			}
			ctx.Resp.Flush()
		case <-notify:
			go unregister()
			break loop
		case <-shutdownCtx.Done():
			go unregister()
			break loop
		case event, ok := <-messageChan:
			if !ok {
				break loop
			}

			// Handle logout
			if event.Name == "logout" {
				if ctx.Session.ID() == event.Data {
					_, _ = (&eventsource.Event{
						Name: "logout",
						Data: "here",
					}).WriteTo(ctx.Resp)
					ctx.Resp.Flush()
					go unregister()
					auth.HandleSignOut(ctx)
					break loop
				}
				// Replace the event - we don't want to expose the session ID to the user
				event = &eventsource.Event{
					Name: "logout",
					Data: "elsewhere",
				}
			}

			_, err := event.WriteTo(ctx.Resp)
			if err != nil {
				log.Error("Unable to write to EventStream for user %s: %v", ctx.Doer.Name, err)
				go unregister()
				break loop
			}
			ctx.Resp.Flush()
		}
	}
	timer.Stop()
}
