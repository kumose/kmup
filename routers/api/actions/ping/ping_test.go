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

package ping

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	pingv1 "github.com/kumose/actions-proto-go/ping/v1"
	"github.com/kumose/actions-proto-go/ping/v1/pingv1connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(pingv1connect.NewPingServiceHandler(
		&Service{},
	))
	MainServiceTest(t, mux)
}

func MainServiceTest(t *testing.T, h http.Handler) {
	t.Parallel()
	server := httptest.NewUnstartedServer(h)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
	)

	grpcClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)

	grpcWebClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPCWeb(),
	)

	clients := []pingv1connect.PingServiceClient{connectClient, grpcClient, grpcWebClient}
	t.Run("ping request", func(t *testing.T) {
		for _, client := range clients {
			result, err := client.Ping(t.Context(), connect.NewRequest(&pingv1.PingRequest{
				Data: "foobar",
			}))
			require.NoError(t, err)
			assert.Equal(t, "Hello, foobar!", result.Msg.Data)
		}
	})
}
