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

package nosql

import (
	"net/url"
	"testing"
)

func TestRedisUsernameOpt(t *testing.T) {
	uri, _ := url.Parse("redis://redis:password@myredis/0")
	opts := getRedisOptions(uri)

	if opts.Username != "redis" {
		t.Fail()
	}
}

func TestRedisPasswordOpt(t *testing.T) {
	uri, _ := url.Parse("redis://redis:password@myredis/0")
	opts := getRedisOptions(uri)

	if opts.Password != "password" {
		t.Fail()
	}
}

func TestSkipVerifyOpt(t *testing.T) {
	uri, _ := url.Parse("rediss://myredis/0?skipverify=true")
	tlsConfig := getRedisTLSOptions(uri)

	if !tlsConfig.InsecureSkipVerify {
		t.Fail()
	}
}

func TestInsecureSkipVerifyOpt(t *testing.T) {
	uri, _ := url.Parse("rediss://myredis/0?insecureskipverify=true")
	tlsConfig := getRedisTLSOptions(uri)

	if !tlsConfig.InsecureSkipVerify {
		t.Fail()
	}
}

func TestRedisSentinelUsernameOpt(t *testing.T) {
	uri, _ := url.Parse("redis+sentinel://redis:password@myredis/0?sentinelusername=suser&sentinelpassword=spass")
	opts := getRedisOptions(uri).Failover()

	if opts.SentinelUsername != "suser" {
		t.Fail()
	}
}

func TestRedisSentinelPasswordOpt(t *testing.T) {
	uri, _ := url.Parse("redis+sentinel://redis:password@myredis/0?sentinelusername=suser&sentinelpassword=spass")
	opts := getRedisOptions(uri).Failover()

	if opts.SentinelPassword != "spass" {
		t.Fail()
	}
}

func TestRedisDatabaseIndexTcp(t *testing.T) {
	uri, _ := url.Parse("redis://redis:password@myredis/12")
	opts := getRedisOptions(uri)

	if opts.DB != 12 {
		t.Fail()
	}
}

func TestRedisDatabaseIndexUnix(t *testing.T) {
	uri, _ := url.Parse("redis+socket:///var/run/redis.sock?database=12")
	opts := getRedisOptions(uri)

	if opts.DB != 12 {
		t.Fail()
	}
}
