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

package foreachref

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

var (
	nullChar     = []byte("\x00")
	dualNullChar = []byte("\x00\x00")
)

// Format supports specifying and parsing an output format for 'git
// for-each-ref'. See See git-for-each-ref(1) for available fields.
type Format struct {
	// fieldNames hold %(fieldname)s to be passed to the '--format' flag of
	// for-each-ref. See git-for-each-ref(1) for available fields.
	fieldNames []string

	// fieldDelim is the character sequence that is used to separate fields
	// for each reference. fieldDelim and refDelim should be selected to not
	// interfere with each other and to not be present in field values.
	fieldDelim []byte
	// fieldDelimStr is a string representation of fieldDelim. Used to save
	// us from repetitive reallocation whenever we need the delimiter as a
	// string.
	fieldDelimStr string
	// refDelim is the character sequence used to separate reference from
	// each other in the output. fieldDelim and refDelim should be selected
	// to not interfere with each other and to not be present in field
	// values.
	refDelim []byte
}

// NewFormat creates a forEachRefFormat using the specified fieldNames. See
// git-for-each-ref(1) for available fields.
func NewFormat(fieldNames ...string) Format {
	return Format{
		fieldNames:    fieldNames,
		fieldDelim:    nullChar,
		fieldDelimStr: string(nullChar),
		refDelim:      dualNullChar,
	}
}

// Flag returns a for-each-ref --format flag value that captures the fieldNames.
func (f Format) Flag() string {
	var formatFlag strings.Builder
	for i, field := range f.fieldNames {
		// field key and field value
		formatFlag.WriteString(fmt.Sprintf("%s %%(%s)", field, field))

		if i < len(f.fieldNames)-1 {
			// note: escape delimiters to allow control characters as
			// delimiters. For example, '%00' for null character or '%0a'
			// for newline.
			formatFlag.WriteString(f.hexEscaped(f.fieldDelim))
		}
	}
	formatFlag.WriteString(f.hexEscaped(f.refDelim))
	return formatFlag.String()
}

// Parser returns a Parser capable of parsing 'git for-each-ref' output produced
// with this Format.
func (f Format) Parser(r io.Reader) *Parser {
	return NewParser(r, f)
}

// hexEscaped produces hex-escpaed characters from a string. For example, "\n\0"
// would turn into "%0a%00".
func (f Format) hexEscaped(delim []byte) string {
	escaped := ""
	for i := range delim {
		escaped += "%" + hex.EncodeToString([]byte{delim[i]})
	}
	return escaped
}
