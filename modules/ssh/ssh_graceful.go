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

package ssh

import (
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	"github.com/gliderlabs/ssh"
)

func listen(server *ssh.Server) {
	gracefulServer := graceful.NewServer("tcp", server.Addr, "SSH")
	gracefulServer.PerWriteTimeout = setting.SSH.PerWriteTimeout
	gracefulServer.PerWritePerKbTimeout = setting.SSH.PerWritePerKbTimeout

	err := gracefulServer.ListenAndServe(server.Serve, setting.SSH.UseProxyProtocol)
	if err != nil {
		select {
		case <-graceful.GetManager().IsShutdown():
			log.Critical("Failed to start SSH server: %v", err)
		default:
			log.Fatal("Failed to start SSH server: %v", err)
		}
	}
	log.Info("SSH Listener: %s Closed", server.Addr)
}

// builtinUnused informs our cleanup routine that we will not be using a ssh port
func builtinUnused() {
	graceful.GetManager().InformCleanup()
}
