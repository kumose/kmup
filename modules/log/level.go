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
	"bytes"
	"strings"

	"github.com/kumose/kmup/modules/json"
)

// Level is the level of the logger
type Level int

const (
	UNDEFINED Level = iota
	TRACE
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	NONE
)

const CRITICAL = ERROR // most logger frameworks doesn't support CRITICAL, and it doesn't seem useful

var toString = map[Level]string{
	UNDEFINED: "undefined",

	TRACE: "trace",
	DEBUG: "debug",
	INFO:  "info",
	WARN:  "warn",
	ERROR: "error",

	FATAL: "fatal",
	NONE:  "none",
}

var toLevel = map[string]Level{
	"undefined": UNDEFINED,

	"trace":   TRACE,
	"debug":   DEBUG,
	"info":    INFO,
	"warn":    WARN,
	"warning": WARN,
	"error":   ERROR,

	"fatal": FATAL,
	"none":  NONE,
}

var levelToColor = map[Level][]ColorAttribute{
	TRACE: {Bold, FgCyan},
	DEBUG: {Bold, FgBlue},
	INFO:  {Bold, FgGreen},
	WARN:  {Bold, FgYellow},
	ERROR: {Bold, FgRed},
	FATAL: {Bold, BgRed},
	NONE:  {Reset},
}

func (l Level) String() string {
	s, ok := toString[l]
	if ok {
		return s
	}
	return "info"
}

func (l Level) ColorAttributes() []ColorAttribute {
	color, ok := levelToColor[l]
	if ok {
		return color
	}
	none := levelToColor[NONE]
	return none
}

// MarshalJSON takes a Level and turns it into text
func (l Level) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[l])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON takes text and turns it into a Level
func (l *Level) UnmarshalJSON(b []byte) error {
	var tmp any
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	switch v := tmp.(type) {
	case string:
		*l = LevelFromString(v)
	case int:
		*l = LevelFromString(Level(v).String())
	default:
		*l = INFO
	}
	return nil
}

// LevelFromString takes a level string and returns a Level
func LevelFromString(level string) Level {
	if l, ok := toLevel[strings.ToLower(level)]; ok {
		return l
	}
	return INFO
}
