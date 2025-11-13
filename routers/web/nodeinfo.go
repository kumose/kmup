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

package web

import (
	"net/http"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/context"
)

type nodeInfoLinks struct {
	Links []nodeInfoLink `json:"links"`
}

type nodeInfoLink struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

// NodeInfoLinks returns links to the node info endpoint
func NodeInfoLinks(ctx *context.Context) {
	nodeinfolinks := &nodeInfoLinks{
		Links: []nodeInfoLink{{
			setting.AppURL + "api/v1/nodeinfo",
			"http://nodeinfo.diaspora.software/ns/schema/2.1",
		}},
	}
	ctx.JSON(http.StatusOK, nodeinfolinks)
}
