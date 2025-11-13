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

package svg

import (
	"bytes"
	"fmt"
	"regexp"
	"sync"
)

type globalVarsStruct struct {
	reXMLDoc,
	reComment,
	reAttrXMLNs,
	reAttrSize,
	reAttrClassPrefix *regexp.Regexp
}

var globalVars = sync.OnceValue(func() *globalVarsStruct {
	return &globalVarsStruct{
		reXMLDoc:  regexp.MustCompile(`(?s)<\?xml.*?>`),
		reComment: regexp.MustCompile(`(?s)<!--.*?-->`),

		reAttrXMLNs:       regexp.MustCompile(`(?s)\s+xmlns\s*=\s*"[^"]*"`),
		reAttrSize:        regexp.MustCompile(`(?s)\s+(width|height)\s*=\s*"[^"]+"`),
		reAttrClassPrefix: regexp.MustCompile(`(?s)\s+class\s*=\s*"`),
	}
})

// Normalize normalizes the SVG content: set default width/height, remove unnecessary tags/attributes
// It's designed to work with valid SVG content. For invalid SVG content, the returned content is not guaranteed.
func Normalize(data []byte, size int) []byte {
	vars := globalVars()
	data = vars.reXMLDoc.ReplaceAll(data, nil)
	data = vars.reComment.ReplaceAll(data, nil)

	data = bytes.TrimSpace(data)
	svgTag, svgRemaining, ok := bytes.Cut(data, []byte(">"))
	if !ok || !bytes.HasPrefix(svgTag, []byte(`<svg`)) {
		return data
	}
	normalized := bytes.Clone(svgTag)
	normalized = vars.reAttrXMLNs.ReplaceAll(normalized, nil)
	normalized = vars.reAttrSize.ReplaceAll(normalized, nil)
	normalized = vars.reAttrClassPrefix.ReplaceAll(normalized, []byte(` class="`))
	normalized = bytes.TrimSpace(normalized)
	normalized = fmt.Appendf(normalized, ` width="%d" height="%d"`, size, size)
	if !bytes.Contains(normalized, []byte(` class="`)) {
		normalized = append(normalized, ` class="svg"`...)
	}
	normalized = append(normalized, '>')
	normalized = append(normalized, svgRemaining...)
	return normalized
}
