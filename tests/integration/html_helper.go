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

package integration

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

// HTMLDoc struct
type HTMLDoc struct {
	doc *goquery.Document
}

// NewHTMLParser parse html file
func NewHTMLParser(t testing.TB, body *bytes.Buffer) *HTMLDoc {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(body)
	assert.NoError(t, err)
	return &HTMLDoc{doc: doc}
}

// GetInputValueByName for get input value by name
func (doc *HTMLDoc) GetInputValueByName(name string) string {
	text, _ := doc.doc.Find(`input[name="` + name + `"]`).Attr("value")
	return text
}

// Find gets the descendants of each element in the current set of
// matched elements, filtered by a selector. It returns a new Selection
// object containing these matched elements.
func (doc *HTMLDoc) Find(selector string) *goquery.Selection {
	return doc.doc.Find(selector)
}

// GetCSRF for getting CSRF token value from input
func (doc *HTMLDoc) GetCSRF() string {
	return doc.GetInputValueByName("_csrf")
}

// AssertHTMLElement check if the element by selector exists or does not exist depending on checkExists
func AssertHTMLElement[T int | bool](t testing.TB, doc *HTMLDoc, selector string, checkExists T) {
	sel := doc.doc.Find(selector)
	switch v := any(checkExists).(type) {
	case bool:
		assert.Equal(t, v, sel.Length() > 0)
	case int:
		assert.Equal(t, v, sel.Length())
	}
}
