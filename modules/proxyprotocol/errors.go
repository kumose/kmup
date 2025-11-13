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

package proxyprotocol

import "fmt"

// ErrBadHeader is an error demonstrating a bad proxy header
type ErrBadHeader struct {
	Header []byte
}

func (e *ErrBadHeader) Error() string {
	return fmt.Sprintf("Unexpected proxy header: %v", e.Header)
}

// ErrBadAddressType is an error demonstrating a bad proxy header with bad Address type
type ErrBadAddressType struct {
	Address string
}

func (e *ErrBadAddressType) Error() string {
	return "Unexpected proxy header address type: " + e.Address
}

// ErrBadRemote is an error demonstrating a bad proxy header with bad Remote
type ErrBadRemote struct {
	IP   string
	Port string
}

func (e *ErrBadRemote) Error() string {
	return fmt.Sprintf("Unexpected proxy header remote IP and port: %s %s", e.IP, e.Port)
}

// ErrBadLocal is an error demonstrating a bad proxy header with bad Local
type ErrBadLocal struct {
	IP   string
	Port string
}

func (e *ErrBadLocal) Error() string {
	return fmt.Sprintf("Unexpected proxy header local IP and port: %s %s", e.IP, e.Port)
}
