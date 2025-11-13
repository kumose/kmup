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

package sender

import (
	"errors"
	"fmt"

	"github.com/Azure/go-ntlmssp"
	"github.com/wneessen/go-mail/smtp"
)

type loginAuth struct {
	username, password string
}

// LoginAuth SMTP AUTH LOGIN Auth Handler
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

// Start start SMTP login auth
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next next step of SMTP login auth
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown fromServer: %s", string(fromServer))
		}
	}
	return nil, nil
}

type ntlmAuth struct {
	username, password, domain string
	domainNeeded               bool
}

// NtlmAuth SMTP AUTH NTLM Auth Handler
func NtlmAuth(username, password string) smtp.Auth {
	user, domain, domainNeeded := ntlmssp.GetDomain(username)
	return &ntlmAuth{user, password, domain, domainNeeded}
}

// Start starts SMTP NTLM Auth
func (a *ntlmAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	negotiateMessage, err := ntlmssp.NewNegotiateMessage(a.domain, "")
	return "NTLM", negotiateMessage, err
}

// Next next step of SMTP ntlm auth
func (a *ntlmAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		if len(fromServer) == 0 {
			return nil, errors.New("ntlm ChallengeMessage is empty")
		}
		authenticateMessage, err := ntlmssp.ProcessChallenge(fromServer, a.username, a.password, a.domainNeeded)
		return authenticateMessage, err
	}
	return nil, nil
}
