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

// Package log provides logging capabilities for Kmup.
// Concepts:
//
// * Logger: a Logger provides logging functions and dispatches log events to all its writers
//
// * EventWriter: written log Event to a destination (eg: file, console)
//   - EventWriterBase: the base struct of a writer, it contains common fields and functions for all writers
//   - WriterType: the type name of a writer, eg: console, file
//   - WriterName: aka Mode Name in document, the name of a writer instance, it's usually defined by the config file.
//     It is called "mode name" because old code use MODE as config key, to keep compatibility, keep this concept.
//
// * WriterMode: the common options for all writers, eg: log level.
//   - WriterConsoleOption and others: the specified options for a writer, eg: file path, remote address.
//
// Call graph:
// -> log.Info()
// -> LoggerImpl.Log()
// -> LoggerImpl.SendLogEvent, then the event goes into writer's goroutines
// -> EventWriter.Run() handles the events
package log

// BaseLogger provides the basic logging functions
type BaseLogger interface {
	Log(skip int, event *Event, format string, v ...any)
	GetLevel() Level
}

// LevelLogger provides level-related logging functions
type LevelLogger interface {
	LevelEnabled(level Level) bool

	Trace(format string, v ...any)
	Debug(format string, v ...any)
	Info(format string, v ...any)
	Warn(format string, v ...any)
	Error(format string, v ...any)
	Critical(format string, v ...any)
}

type Logger interface {
	BaseLogger
	LevelLogger
}

type LogStringer interface { //nolint:revive // export stutter
	LogString() string
}
