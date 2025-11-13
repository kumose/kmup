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
	"context"
	"fmt"
	"net/http"

	"github.com/kumose/kmup/modules/log"

	"connectrpc.com/connect"
	pingv1 "github.com/kumose/actions-proto-go/ping/v1"
	"github.com/kumose/actions-proto-go/ping/v1/pingv1connect"
)

func NewPingServiceHandler() (string, http.Handler) {
	return pingv1connect.NewPingServiceHandler(&Service{})
}

var _ pingv1connect.PingServiceHandler = (*Service)(nil)

type Service struct{}

func (s *Service) Ping(
	ctx context.Context,
	req *connect.Request[pingv1.PingRequest],
) (*connect.Response[pingv1.PingResponse], error) {
	log.Trace("Content-Type: %s", req.Header().Get("Content-Type"))
	log.Trace("User-Agent: %s", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&pingv1.PingResponse{
		Data: fmt.Sprintf("Hello, %s!", req.Msg.Data),
	})
	return res, nil
}
