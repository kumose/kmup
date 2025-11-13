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

package ldap

// SecurityProtocol protocol type
type SecurityProtocol int

// Note: new type must be added at the end of list to maintain compatibility.
const (
	SecurityProtocolUnencrypted SecurityProtocol = iota
	SecurityProtocolLDAPS
	SecurityProtocolStartTLS
)

// String returns the name of the SecurityProtocol
func (s SecurityProtocol) String() string {
	return SecurityProtocolNames[s]
}

// Int returns the int value of the SecurityProtocol
func (s SecurityProtocol) Int() int {
	return int(s)
}

// SecurityProtocolNames contains the name of SecurityProtocol values.
var SecurityProtocolNames = map[SecurityProtocol]string{
	SecurityProtocolUnencrypted: "Unencrypted",
	SecurityProtocolLDAPS:       "LDAPS",
	SecurityProtocolStartTLS:    "StartTLS",
}
