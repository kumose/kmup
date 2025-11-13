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

package context

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/kumose/kmup/modules/paginator"
)

// Pagination provides a pagination via paginator.Paginator and additional configurations for the link params used in rendering
type Pagination struct {
	Paginater *paginator.Paginator
	urlParams []string
}

// NewPagination creates a new instance of the Pagination struct.
// "pagingNum" is "page size" or "limit", "current" is "page"
// total=-1 means only showing prev/next
func NewPagination(total, pagingNum, current, numPages int) *Pagination {
	p := &Pagination{}
	p.Paginater = paginator.New(total, pagingNum, current, numPages)
	return p
}

func (p *Pagination) WithCurRows(n int) *Pagination {
	p.Paginater.SetCurRows(n)
	return p
}

func (p *Pagination) AddParamFromQuery(q url.Values) {
	for key, values := range q {
		if key == "page" || len(values) == 0 || (len(values) == 1 && values[0] == "") {
			continue
		}
		for _, value := range values {
			urlParam := fmt.Sprintf("%s=%v", url.QueryEscape(key), url.QueryEscape(value))
			p.urlParams = append(p.urlParams, urlParam)
		}
	}
}

func (p *Pagination) AddParamFromRequest(req *http.Request) {
	p.AddParamFromQuery(req.URL.Query())
}

// GetParams returns the configured URL params
func (p *Pagination) GetParams() template.URL {
	return template.URL(strings.Join(p.urlParams, "&"))
}
