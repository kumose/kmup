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

package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsePostgreSQLHostPort(t *testing.T) {
	tests := map[string]struct {
		HostPort string
		Host     string
		Port     string
	}{
		"host-port": {
			HostPort: "127.0.0.1:1234",
			Host:     "127.0.0.1",
			Port:     "1234",
		},
		"no-port": {
			HostPort: "127.0.0.1",
			Host:     "127.0.0.1",
			Port:     "5432",
		},
		"ipv6-port": {
			HostPort: "[::1]:1234",
			Host:     "::1",
			Port:     "1234",
		},
		"ipv6-no-port": {
			HostPort: "[::1]",
			Host:     "::1",
			Port:     "5432",
		},
		"unix-socket": {
			HostPort: "/tmp/pg.sock:1234",
			Host:     "/tmp/pg.sock",
			Port:     "1234",
		},
		"unix-socket-no-port": {
			HostPort: "/tmp/pg.sock",
			Host:     "/tmp/pg.sock",
			Port:     "5432",
		},
	}
	for k, test := range tests {
		t.Run(k, func(t *testing.T) {
			t.Log(test.HostPort)
			host, port := parsePostgreSQLHostPort(test.HostPort)
			assert.Equal(t, test.Host, host)
			assert.Equal(t, test.Port, port)
		})
	}
}

func Test_getPostgreSQLConnectionString(t *testing.T) {
	tests := []struct {
		Host    string
		User    string
		Passwd  string
		Name    string
		SSLMode string
		Output  string
	}{
		{
			Host:   "", // empty means default
			Output: "postgres://:@127.0.0.1:5432?sslmode=",
		},
		{
			Host:    "/tmp/pg.sock",
			User:    "testuser",
			Passwd:  "space space !#$%^^%^```-=?=",
			Name:    "kmup",
			SSLMode: "false",
			Output:  "postgres://testuser:space%20space%20%21%23$%25%5E%5E%25%5E%60%60%60-=%3F=@:5432/kmup?host=%2Ftmp%2Fpg.sock&sslmode=false",
		},
		{
			Host:    "/tmp/pg.sock:6432",
			User:    "testuser",
			Passwd:  "pass",
			Name:    "kmup",
			SSLMode: "false",
			Output:  "postgres://testuser:pass@:6432/kmup?host=%2Ftmp%2Fpg.sock&sslmode=false",
		},
		{
			Host:    "localhost",
			User:    "pgsqlusername",
			Passwd:  "I love Kmup!",
			Name:    "kmup",
			SSLMode: "true",
			Output:  "postgres://pgsqlusername:I%20love%20Kmup%21@localhost:5432/kmup?sslmode=true",
		},
		{
			Host:   "localhost:1234",
			User:   "user",
			Passwd: "pass",
			Name:   "kmup?param=1",
			Output: "postgres://user:pass@localhost:1234/kmup?param=1&sslmode=",
		},
	}

	for _, test := range tests {
		connStr := getPostgreSQLConnectionString(test.Host, test.User, test.Passwd, test.Name, test.SSLMode)
		assert.Equal(t, test.Output, connStr)
	}
}
