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

package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/util"
)

// RedirectURL returns the redirect URL of a http response.
// It also works for JSONRedirect: `{"redirect": "..."}`
// FIXME: it should separate the logic of checking from header and JSON body
func RedirectURL(resp http.ResponseWriter) string {
	loc := resp.Header().Get("Location")
	if loc != "" {
		return loc
	}
	if r, ok := resp.(*httptest.ResponseRecorder); ok {
		m := map[string]any{}
		err := json.Unmarshal(r.Body.Bytes(), &m)
		if err == nil {
			if loc, ok := m["redirect"].(string); ok {
				return loc
			}
		}
	}
	return ""
}

func ParseJSONError(buf []byte) (ret struct {
	ErrorMessage string `json:"errorMessage"`
	RenderFormat string `json:"renderFormat"`
},
) {
	_ = json.Unmarshal(buf, &ret)
	return ret
}

func IsNormalPageCompleted(s string) bool {
	return strings.Contains(s, `<footer class="page-footer"`) && strings.Contains(s, `</html>`)
}

func MockVariableValue[T any](p *T, v ...T) (reset func()) {
	old := *p
	if len(v) > 0 {
		*p = v[0]
	}
	return func() { *p = old }
}

// SetupKmupRoot Sets KMUP_ROOT if it is not already set and returns the value
func SetupKmupRoot() string {
	kmupRoot := os.Getenv("KMUP_ROOT")
	if kmupRoot != "" {
		return kmupRoot
	}
	_, filename, _, _ := runtime.Caller(0)
	kmupRoot = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	fixturesDir := filepath.Join(kmupRoot, "models", "fixtures")
	if exist, _ := util.IsDir(fixturesDir); !exist {
		panic("fixtures directory not found: " + fixturesDir)
	}
	_ = os.Setenv("KMUP_ROOT", kmupRoot)
	return kmupRoot
}
