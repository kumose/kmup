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
)

// EventWriter is the general interface for all event writers
// EventWriterBase is only used as its base interface
// A writer implementation could override the default EventWriterBase functions
// eg: a writer can override the Run to handle events in its own way with its own goroutine
type EventWriter interface {
	EventWriterBase
}

// WriterMode is the mode for creating a new EventWriter, it contains common options for all writers
// Its WriterOption field is the specified options for a writer, it should be passed by value but not by pointer
type WriterMode struct {
	BufferLen int

	Level    Level
	Prefix   string
	Colorize bool
	Flags    Flags

	Expression string

	StacktraceLevel Level

	WriterOption any
}

// EventWriterProvider is the function for creating a new EventWriter
type EventWriterProvider func(writerName string, writerMode WriterMode) EventWriter

var eventWriterProviders = map[string]EventWriterProvider{}

func RegisterEventWriter(writerType string, p EventWriterProvider) {
	eventWriterProviders[writerType] = p
}

func HasEventWriter(writerType string) bool {
	_, ok := eventWriterProviders[writerType]
	return ok
}

func NewEventWriter(name, writerType string, mode WriterMode) (EventWriter, error) {
	if p, ok := eventWriterProviders[writerType]; ok {
		return p(name, mode), nil
	}
	return nil, fmt.Errorf("unknown event writer type %q for writer %q", writerType, name)
}
