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

package e2e

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func onKmupRunTB(t testing.TB, callback func(testing.TB, *url.URL), prepare ...bool) {
	if len(prepare) == 0 || prepare[0] {
		defer tests.PrepareTestEnv(t, 1)()
	}
	s := http.Server{
		Handler: testE2eWebRoutes,
	}

	u, err := url.Parse(setting.AppURL)
	assert.NoError(t, err)
	listener, err := net.Listen("tcp", u.Host)
	i := 0
	for err != nil && i <= 10 {
		time.Sleep(100 * time.Millisecond)
		listener, err = net.Listen("tcp", u.Host)
		i++
	}
	assert.NoError(t, err)
	u.Host = listener.Addr().String()

	defer func() {
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
		s.Shutdown(ctx)
		cancel()
	}()

	go s.Serve(listener)
	// Started by config go ssh.Listen(setting.SSH.ListenHost, setting.SSH.ListenPort, setting.SSH.ServerCiphers, setting.SSH.ServerKeyExchanges, setting.SSH.ServerMACs)

	callback(t, u)
}

func onKmupRun(t *testing.T, callback func(*testing.T, *url.URL), prepare ...bool) {
	onKmupRunTB(t, func(t testing.TB, u *url.URL) {
		callback(t.(*testing.T), u)
	}, prepare...)
}
