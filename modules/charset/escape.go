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

//go:generate go run invisible/generate.go -v -o ./invisible_gen.go

//go:generate go run ambiguous/generate.go -v -o ./ambiguous_gen.go ambiguous/ambiguous.json

package charset

import (
	"html/template"
	"io"
	"strings"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/translation"
)

// RuneNBSP is the codepoint for NBSP
const RuneNBSP = 0xa0

// EscapeControlHTML escapes the unicode control sequences in a provided html document
func EscapeControlHTML(html template.HTML, locale translation.Locale, allowed ...rune) (escaped *EscapeStatus, output template.HTML) {
	sb := &strings.Builder{}
	escaped, _ = EscapeControlReader(strings.NewReader(string(html)), sb, locale, allowed...) // err has been handled in EscapeControlReader
	return escaped, template.HTML(sb.String())
}

// EscapeControlReader escapes the unicode control sequences in a provided reader of HTML content and writer in a locale and returns the findings as an EscapeStatus
func EscapeControlReader(reader io.Reader, writer io.Writer, locale translation.Locale, allowed ...rune) (escaped *EscapeStatus, err error) {
	if !setting.UI.AmbiguousUnicodeDetection {
		_, err = io.Copy(writer, reader)
		return &EscapeStatus{}, err
	}
	outputStream := &HTMLStreamerWriter{Writer: writer}
	streamer := NewEscapeStreamer(locale, outputStream, allowed...).(*escapeStreamer)

	if err = StreamHTML(reader, streamer); err != nil {
		streamer.escaped.HasError = true
		log.Error("Error whilst escaping: %v", err)
	}
	return streamer.escaped, err
}
