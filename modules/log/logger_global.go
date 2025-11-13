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
	"fmt"
	"os"
)

// FallbackErrorf is the last chance to show an error if the logger has internal errors
func FallbackErrorf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func GetLevel() Level {
	return GetLogger(DEFAULT).GetLevel()
}

func Log(skip int, level Level, format string, v ...any) {
	// codeql[disable-next-line=go/clear-text-logging]
	GetLogger(DEFAULT).Log(skip+1, &Event{Level: level}, format, v...)
}

func Trace(format string, v ...any) {
	Log(1, TRACE, format, v...)
}

func IsTrace() bool {
	return GetLevel() <= TRACE
}

func Debug(format string, v ...any) {
	Log(1, DEBUG, format, v...)
}

func IsDebug() bool {
	return GetLevel() <= DEBUG
}

func Info(format string, v ...any) {
	Log(1, INFO, format, v...)
}

func Warn(format string, v ...any) {
	Log(1, WARN, format, v...)
}

func Error(format string, v ...any) {
	Log(1, ERROR, format, v...)
}

func ErrorWithSkip(skip int, format string, v ...any) {
	Log(skip+1, ERROR, format, v...)
}

func Critical(format string, v ...any) {
	Log(1, ERROR, format, v...)
}

var OsExiter = os.Exit

// Fatal records fatal log and exit process
func Fatal(format string, v ...any) {
	Log(1, FATAL, format, v...)
	GetManager().Close()
	OsExiter(1)
}

func GetLogger(name string) Logger {
	return GetManager().GetLogger(name)
}

func IsLoggerEnabled(name string) bool {
	return GetManager().GetLogger(name).IsEnabled()
}

func SetConsoleLogger(loggerName, writerName string, level Level) {
	writer := NewEventWriterConsole(writerName, WriterMode{
		Level:        level,
		Flags:        FlagsFromBits(LstdFlags),
		Colorize:     CanColorStdout,
		WriterOption: WriterConsoleOption{},
	})
	GetManager().GetLogger(loggerName).ReplaceAllWriters(writer)
}
