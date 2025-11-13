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

package internal

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/kumose/kmup/modules/htmlutil"

	"golang.org/x/net/html"
)

var reAttrClass = sync.OnceValue(func() *regexp.Regexp {
	// TODO: it isn't a problem at the moment because our HTML contents are always well constructed
	return regexp.MustCompile(`(<[^>]+)\s+class="([^"]+)"([^>]*>)`)
})

// RenderInternal also works without initialization
// If no initialization (no secureID), it will not protect any attributes and return the original name&value
type RenderInternal struct {
	secureID       string
	secureIDPrefix string
}

func (r *RenderInternal) Init(output io.Writer, extraHeadHTML template.HTML) io.WriteCloser {
	buf := make([]byte, 12)
	_, err := rand.Read(buf)
	if err != nil {
		panic("unable to generate secure id")
	}
	return r.init(base64.URLEncoding.EncodeToString(buf), output, extraHeadHTML)
}

func (r *RenderInternal) init(secID string, output io.Writer, extraHeadHTML template.HTML) io.WriteCloser {
	r.secureID = secID
	r.secureIDPrefix = r.secureID + ":"
	return &finalProcessor{renderInternal: r, output: output, extraHeadHTML: extraHeadHTML}
}

func (r *RenderInternal) RecoverProtectedValue(v string) (string, bool) {
	if !strings.HasPrefix(v, r.secureIDPrefix) {
		return "", false
	}
	return v[len(r.secureIDPrefix):], true
}

func (r *RenderInternal) SafeAttr(name string) string {
	if r.secureID == "" {
		return name
	}
	return "data-attr-" + name
}

func (r *RenderInternal) SafeValue(val string) string {
	if r.secureID == "" {
		return val
	}
	return r.secureID + ":" + val
}

func (r *RenderInternal) NodeSafeAttr(attr, val string) html.Attribute {
	return html.Attribute{Key: r.SafeAttr(attr), Val: r.SafeValue(val)}
}

func (r *RenderInternal) ProtectSafeAttrs(content template.HTML) template.HTML {
	if r.secureID == "" {
		return content
	}
	return template.HTML(reAttrClass().ReplaceAllString(string(content), `$1 data-attr-class="`+r.secureIDPrefix+`$2"$3`))
}

func (r *RenderInternal) FormatWithSafeAttrs(w io.Writer, fmt template.HTML, a ...any) error {
	_, err := w.Write([]byte(r.ProtectSafeAttrs(htmlutil.HTMLFormat(fmt, a...))))
	return err
}
