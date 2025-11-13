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

// This code is heavily inspired by the archived gofacebook/gracenet/net.go handler

//go:build windows

package graceful

import "net"

// DefaultGetListener obtains a listener for the local network address.
// On windows this is basically just a shim around net.Listen.
func DefaultGetListener(network, address string) (net.Listener, error) {
	// Add a deferral to say that we've tried to grab a listener
	defer GetManager().InformCleanup()

	return net.Listen(network, address)
}
