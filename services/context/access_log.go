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

package context

import (
	"bytes"
	"net"
	"net/http"
	"strings"
	"text/template"
	"time"
	"unicode"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/web/middleware"
)

type accessLoggerTmplData struct {
	Identity       *string
	Start          *time.Time
	ResponseWriter struct {
		Status, Size int
	}
	Ctx       map[string]any
	RequestID *string
}

const keyOfRequestIDInTemplate = ".RequestID"

// According to:
// TraceId: A valid trace identifier is a 16-byte array with at least one non-zero byte
// MD5 output is 16 or 32 bytes: md5-bytes is 16, md5-hex is 32
// SHA1: similar, SHA1-bytes is 20, SHA1-hex is 40.
// UUID is 128-bit, 32 hex chars, 36 ASCII chars with 4 dashes
// So, we accept a Request ID with a maximum character length of 40
const maxRequestIDByteLength = 40

func isSafeRequestID(id string) bool {
	for _, r := range id {
		safe := unicode.IsPrint(r)
		if !safe {
			return false
		}
	}
	return true
}

func parseRequestIDFromRequestHeader(req *http.Request) string {
	requestID := "-"
	for _, key := range setting.Log.RequestIDHeaders {
		if req.Header.Get(key) != "" {
			requestID = req.Header.Get(key)
			break
		}
	}
	if !isSafeRequestID(requestID) {
		return "-"
	}
	if len(requestID) > maxRequestIDByteLength {
		requestID = requestID[:maxRequestIDByteLength] + "..."
	}
	return requestID
}

type accessLogRecorder struct {
	logger        log.BaseLogger
	logTemplate   *template.Template
	needRequestID bool
}

func (lr *accessLogRecorder) record(start time.Time, respWriter ResponseWriter, req *http.Request) {
	var requestID string
	if lr.needRequestID {
		requestID = parseRequestIDFromRequestHeader(req)
	}

	reqHost, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		reqHost = req.RemoteAddr
	}

	identity := "-"
	data := middleware.GetContextData(req.Context())
	if signedUser, ok := data[middleware.ContextDataKeySignedUser].(*user_model.User); ok {
		identity = signedUser.Name
	}
	buf := bytes.NewBuffer([]byte{})
	tmplData := accessLoggerTmplData{
		Identity: &identity,
		Start:    &start,
		Ctx: map[string]any{
			"RemoteAddr": req.RemoteAddr,
			"RemoteHost": reqHost,
			"Req":        req,
		},
		RequestID: &requestID,
	}
	tmplData.ResponseWriter.Status = respWriter.WrittenStatus()
	tmplData.ResponseWriter.Size = respWriter.WrittenSize()
	err = lr.logTemplate.Execute(buf, tmplData)
	if err != nil {
		log.Error("Could not execute access logger template: %v", err.Error())
	}

	lr.logger.Log(1, &log.Event{Level: log.INFO}, "%s", buf.String())
}

func newAccessLogRecorder() *accessLogRecorder {
	return &accessLogRecorder{
		logger:        log.GetLogger("access"),
		logTemplate:   template.Must(template.New("log").Parse(setting.Log.AccessLogTemplate)),
		needRequestID: len(setting.Log.RequestIDHeaders) > 0 && strings.Contains(setting.Log.AccessLogTemplate, keyOfRequestIDInTemplate),
	}
}

// AccessLogger returns a middleware to log access logger
func AccessLogger() func(http.Handler) http.Handler {
	recorder := newAccessLogRecorder()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, req)
			recorder.record(start, w.(ResponseWriter), req)
		})
	}
}
