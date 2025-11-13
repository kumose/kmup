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
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

type testAccessLoggerMock struct {
	logs []string
}

func (t *testAccessLoggerMock) Log(skip int, event *log.Event, format string, v ...any) {
	t.logs = append(t.logs, fmt.Sprintf(format, v...))
}

func (t *testAccessLoggerMock) GetLevel() log.Level {
	return log.INFO
}

type testAccessLoggerResponseWriterMock struct{}

func (t testAccessLoggerResponseWriterMock) Header() http.Header {
	return nil
}

func (t testAccessLoggerResponseWriterMock) Before(f func(ResponseWriter)) {}

func (t testAccessLoggerResponseWriterMock) WriteHeader(statusCode int) {}

func (t testAccessLoggerResponseWriterMock) Write(bytes []byte) (int, error) {
	return 0, nil
}

func (t testAccessLoggerResponseWriterMock) Flush() {}

func (t testAccessLoggerResponseWriterMock) WrittenStatus() int {
	return http.StatusOK
}

func (t testAccessLoggerResponseWriterMock) WrittenSize() int {
	return 123123
}

func TestAccessLogger(t *testing.T) {
	setting.Log.AccessLogTemplate = `{{.Ctx.RemoteHost}} - {{.Identity}} {{.Start.Format "[02/Jan/2006:15:04:05 -0700]" }} "{{.Ctx.Req.Method}} {{.Ctx.Req.URL.RequestURI}} {{.Ctx.Req.Proto}}" {{.ResponseWriter.Status}} {{.ResponseWriter.Size}} "{{.Ctx.Req.Referer}}" "{{.Ctx.Req.UserAgent}}"`
	recorder := newAccessLogRecorder()
	mockLogger := &testAccessLoggerMock{}
	recorder.logger = mockLogger
	req := &http.Request{
		RemoteAddr: "remote-addr",
		Method:     http.MethodGet,
		Proto:      "https",
		URL:        &url.URL{Path: "/path"},
	}
	req.Header = http.Header{}
	req.Header.Add("Referer", "referer")
	req.Header.Add("User-Agent", "user-agent")
	recorder.record(time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC), &testAccessLoggerResponseWriterMock{}, req)
	assert.Equal(t, []string{`remote-addr - - [02/Jan/2000:03:04:05 +0000] "GET /path https" 200 123123 "referer" "user-agent"`}, mockLogger.logs)
}

func TestAccessLoggerRequestID(t *testing.T) {
	assert.False(t, isSafeRequestID("\x00"))
	assert.True(t, isSafeRequestID("a b-c"))
}
