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

package private

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kumose/kmup/modules/httplib"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/proxyprotocol"
	"github.com/kumose/kmup/modules/setting"
)

// Response is used for internal request response (for user message and error message)
type Response struct {
	Err     string `json:"err,omitempty"`      // server-side error log message, it won't be exposed to end users
	UserMsg string `json:"user_msg,omitempty"` // meaningful error message for end users, it will be shown in git client's output.
}

func getClientIP() string {
	sshConnEnv := strings.TrimSpace(os.Getenv("SSH_CONNECTION"))
	if len(sshConnEnv) == 0 {
		return "127.0.0.1"
	}
	return strings.Fields(sshConnEnv)[0]
}

func dialContextInternalAPI(ctx context.Context, network, address string) (conn net.Conn, err error) {
	d := net.Dialer{Timeout: 10 * time.Second}
	if setting.Protocol == setting.HTTPUnix {
		conn, err = d.DialContext(ctx, "unix", setting.HTTPAddr)
	} else {
		conn, err = d.DialContext(ctx, network, address)
	}
	if err != nil {
		return nil, err
	}
	if setting.LocalUseProxyProtocol {
		if err = proxyprotocol.WriteLocalHeader(conn); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}
	return conn, nil
}

var internalAPITransport = sync.OnceValue(func() http.RoundTripper {
	return &http.Transport{
		DialContext: dialContextInternalAPI,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         setting.Domain,
		},
	}
})

func NewInternalRequest(ctx context.Context, url, method string) *httplib.Request {
	if setting.InternalToken == "" {
		log.Fatal(`The INTERNAL_TOKEN setting is missing from the configuration file: %q.
Ensure you are running in the correct environment or set the correct configuration file with -c.`, setting.CustomConf)
	}

	if !strings.HasPrefix(url, setting.LocalURL) {
		log.Fatal("Invalid internal request URL: %q", url)
	}

	return httplib.NewRequest(url, method).
		SetContext(ctx).
		SetTransport(internalAPITransport()).
		Header("X-Real-IP", getClientIP()).
		Header("X-Kmup-Internal-Auth", "Bearer "+setting.InternalToken)
}

func newInternalRequestAPI(ctx context.Context, url, method string, body ...any) *httplib.Request {
	req := NewInternalRequest(ctx, url, method)
	if len(body) == 1 {
		req.Header("Content-Type", "application/json")
		jsonBytes, _ := json.Marshal(body[0])
		req.Body(jsonBytes)
	} else if len(body) > 1 {
		log.Fatal("Too many arguments for newInternalRequestAPI")
	}

	req.SetReadWriteTimeout(60 * time.Second)
	return req
}
