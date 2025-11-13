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

package hostmatcher

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"syscall"
	"time"
)

// NewDialContext returns a DialContext for Transport, the DialContext will do allow/block list check
func NewDialContext(usage string, allowList, blockList *HostMatchList, proxy *url.URL) func(ctx context.Context, network, addr string) (net.Conn, error) {
	// How Go HTTP Client works with redirection:
	//   transport.RoundTrip URL=http://domain.com, Host=domain.com
	//   transport.DialContext addrOrHost=domain.com:80
	//   dialer.Control tcp4:11.22.33.44:80
	//   transport.RoundTrip URL=http://www.domain.com/, Host=(empty here, in the direction, HTTP client doesn't fill the Host field)
	//   transport.DialContext addrOrHost=domain.com:80
	//   dialer.Control tcp4:11.22.33.44:80
	return func(ctx context.Context, network, addrOrHost string) (net.Conn, error) {
		dialer := net.Dialer{
			// default values comes from http.DefaultTransport
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,

			Control: func(network, ipAddr string, c syscall.RawConn) error {
				host, port, err := net.SplitHostPort(addrOrHost)
				if err != nil {
					return err
				}
				if proxy != nil {
					// Always allow the host of the proxy, but only on the specified port.
					if host == proxy.Hostname() && port == proxy.Port() {
						return nil
					}
				}

				// in Control func, the addr was already resolved to IP:PORT format, there is no cost to do ResolveTCPAddr here
				tcpAddr, err := net.ResolveTCPAddr(network, ipAddr)
				if err != nil {
					return fmt.Errorf("%s can only call HTTP servers via TCP, deny '%s(%s:%s)', err=%w", usage, host, network, ipAddr, err)
				}

				var blockedError error
				if blockList.MatchHostOrIP(host, tcpAddr.IP) {
					blockedError = fmt.Errorf("%s can not call blocked HTTP servers (check your %s setting), deny '%s(%s)'", usage, blockList.SettingKeyHint, host, ipAddr)
				}

				// if we have an allow-list, check the allow-list first
				if !allowList.IsEmpty() {
					if !allowList.MatchHostOrIP(host, tcpAddr.IP) {
						return fmt.Errorf("%s can only call allowed HTTP servers (check your %s setting), deny '%s(%s)'", usage, allowList.SettingKeyHint, host, ipAddr)
					}
				}
				// otherwise, we always follow the blocked list
				return blockedError
			},
		}
		return dialer.DialContext(ctx, network, addrOrHost)
	}
}
