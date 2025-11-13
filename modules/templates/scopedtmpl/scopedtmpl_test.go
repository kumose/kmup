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

package scopedtmpl

import (
	"bytes"
	"html/template"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScopedTemplateSetFuncMap(t *testing.T) {
	all := template.New("")

	all.Funcs(template.FuncMap{"CtxFunc": func(s string) string {
		return "default"
	}})

	_, err := all.New("base").Parse(`{{CtxFunc "base"}}`)
	assert.NoError(t, err)

	_, err = all.New("test").Parse(strings.TrimSpace(`
{{template "base"}}
{{CtxFunc "test"}}
{{template "base"}}
{{CtxFunc "test"}}
`))
	assert.NoError(t, err)

	ts, err := newScopedTemplateSet(all, "test")
	assert.NoError(t, err)

	// try to use different CtxFunc to render concurrently

	funcMap1 := template.FuncMap{
		"CtxFunc": func(s string) string {
			time.Sleep(100 * time.Millisecond)
			return s + "1"
		},
	}

	funcMap2 := template.FuncMap{
		"CtxFunc": func(s string) string {
			time.Sleep(100 * time.Millisecond)
			return s + "2"
		},
	}

	out1 := bytes.Buffer{}
	out2 := bytes.Buffer{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := ts.newExecutor(funcMap1).Execute(&out1, nil)
		assert.NoError(t, err)
		wg.Done()
	}()
	go func() {
		err := ts.newExecutor(funcMap2).Execute(&out2, nil)
		assert.NoError(t, err)
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, "base1\ntest1\nbase1\ntest1", out1.String())
	assert.Equal(t, "base2\ntest2\nbase2\ntest2", out2.String())
}

func TestScopedTemplateSetEscape(t *testing.T) {
	all := template.New("")
	_, err := all.New("base").Parse(`<a href="?q={{.param}}">{{.text}}</a>`)
	assert.NoError(t, err)

	_, err = all.New("test").Parse(`{{template "base" .}}<form action="?q={{.param}}">{{.text}}</form>`)
	assert.NoError(t, err)

	ts, err := newScopedTemplateSet(all, "test")
	assert.NoError(t, err)

	out := bytes.Buffer{}
	err = ts.newExecutor(nil).Execute(&out, map[string]string{"param": "/", "text": "<"})
	assert.NoError(t, err)

	assert.Equal(t, `<a href="?q=%2f">&lt;</a><form action="?q=%2f">&lt;</form>`, out.String())
}

func TestScopedTemplateSetUnsafe(t *testing.T) {
	all := template.New("")
	_, err := all.New("test").Parse(`<a href="{{if true}}?{{end}}a={{.param}}"></a>`)
	assert.NoError(t, err)

	_, err = newScopedTemplateSet(all, "test")
	assert.ErrorContains(t, err, "appears in an ambiguous context within a URL")
}
