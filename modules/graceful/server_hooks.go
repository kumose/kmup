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

package graceful

import (
	"os"

	"github.com/kumose/kmup/modules/log"
)

// awaitShutdown waits for the shutdown signal from the Manager
func (srv *Server) awaitShutdown() {
	select {
	case <-GetManager().IsShutdown():
		// Shutdown
		srv.doShutdown()
	case <-GetManager().IsHammer():
		// Hammer
		srv.doShutdown()
		srv.doHammer()
	}
	<-GetManager().IsHammer()
	srv.doHammer()
}

// shutdown closes the listener so that no new connections are accepted
// and starts a goroutine that will hammer (stop all running requests) the server
// after setting.GracefulHammerTime.
func (srv *Server) doShutdown() {
	// only shutdown if we're running.
	if srv.getState() != stateRunning {
		return
	}

	srv.setState(stateShuttingDown)

	if srv.OnShutdown != nil {
		srv.OnShutdown()
	}
	err := srv.listener.Close()
	if err != nil {
		log.Error("PID: %d Listener.Close() error: %v", os.Getpid(), err)
	} else {
		log.Info("PID: %d Listener (%s) closed.", os.Getpid(), srv.listener.Addr())
	}
}

func (srv *Server) doHammer() {
	if srv.getState() != stateShuttingDown {
		return
	}
	srv.closeAllConnections()
}
