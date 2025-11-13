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

package smtp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"
)

//   _________   __________________________
//  /   _____/  /     \__    ___/\______   \
//  \_____  \  /  \ /  \|    |    |     ___/
//  /        \/    Y    \    |    |    |
// /_______  /\____|__  /____|    |____|
//         \/         \/

type loginAuthenticator struct {
	username, password string
}

func (auth *loginAuthenticator) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(auth.username), nil
}

func (auth *loginAuthenticator) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(auth.username), nil
		case "Password:":
			return []byte(auth.password), nil
		}
	}
	return nil, nil
}

// SMTP authentication type names.
const (
	PlainAuthentication   = "PLAIN"
	LoginAuthentication   = "LOGIN"
	CRAMMD5Authentication = "CRAM-MD5"
)

// Authenticators contains available SMTP authentication type names.
var Authenticators = []string{PlainAuthentication, LoginAuthentication, CRAMMD5Authentication}

// ErrUnsupportedLoginType login source is unknown error
var ErrUnsupportedLoginType = errors.New("Login source is unknown")

// Authenticate performs an SMTP authentication.
func Authenticate(a smtp.Auth, source *Source) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: source.SkipVerify,
		ServerName:         source.Host,
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(source.Host, strconv.Itoa(source.Port)))
	if err != nil {
		return err
	}
	defer conn.Close()

	if source.UseTLS() {
		conn = tls.Client(conn, tlsConfig)
	}

	client, err := smtp.NewClient(conn, source.Host)
	if err != nil {
		return fmt.Errorf("failed to create NewClient: %w", err)
	}
	defer client.Close()

	if !source.DisableHelo {
		hostname := source.HeloHostname
		if len(hostname) == 0 {
			hostname, err = os.Hostname()
			if err != nil {
				return fmt.Errorf("failed to find Hostname: %w", err)
			}
		}

		if err = client.Hello(hostname); err != nil {
			return fmt.Errorf("failed to send Helo: %w", err)
		}
	}

	// If not using SMTPS, always use STARTTLS if available
	hasStartTLS, _ := client.Extension("STARTTLS")
	if !source.UseTLS() && hasStartTLS {
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start StartTLS: %w", err)
		}
	}

	if ok, _ := client.Extension("AUTH"); ok {
		return client.Auth(a)
	}

	return ErrUnsupportedLoginType
}
