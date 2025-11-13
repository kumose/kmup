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

package log

import (
	"io"
)

type baseToLogger struct {
	base BaseLogger
}

// BaseLoggerToGeneralLogger wraps a BaseLogger (which only has Log() function) to a Logger (which has Info() function)
func BaseLoggerToGeneralLogger(b BaseLogger) Logger {
	l := &baseToLogger{base: b}
	return l
}

var _ Logger = (*baseToLogger)(nil)

func (s *baseToLogger) Log(skip int, event *Event, format string, v ...any) {
	// codeql[disable-next-line=go/clear-text-logging]
	s.base.Log(skip+1, event, format, v...)
}

func (s *baseToLogger) GetLevel() Level {
	return s.base.GetLevel()
}

func (s *baseToLogger) LevelEnabled(level Level) bool {
	return s.base.GetLevel() <= level
}

func (s *baseToLogger) Trace(format string, v ...any) {
	s.base.Log(1, &Event{Level: TRACE}, format, v...)
}

func (s *baseToLogger) Debug(format string, v ...any) {
	s.base.Log(1, &Event{Level: DEBUG}, format, v...)
}

func (s *baseToLogger) Info(format string, v ...any) {
	s.base.Log(1, &Event{Level: INFO}, format, v...)
}

func (s *baseToLogger) Warn(format string, v ...any) {
	s.base.Log(1, &Event{Level: WARN}, format, v...)
}

func (s *baseToLogger) Error(format string, v ...any) {
	s.base.Log(1, &Event{Level: ERROR}, format, v...)
}

func (s *baseToLogger) Critical(format string, v ...any) {
	s.base.Log(1, &Event{Level: CRITICAL}, format, v...)
}

type PrintfLogger struct {
	Logf func(format string, args ...any)
}

func (p *PrintfLogger) Printf(format string, args ...any) {
	p.Logf(format, args...)
}

type loggerToWriter struct {
	logf func(format string, args ...any)
}

func (p *loggerToWriter) Write(bs []byte) (int, error) {
	p.logf("%s", string(bs))
	return len(bs), nil
}

// LoggerToWriter wraps a log function to an io.Writer
func LoggerToWriter(logf func(format string, args ...any)) io.Writer {
	return &loggerToWriter{logf: logf}
}
