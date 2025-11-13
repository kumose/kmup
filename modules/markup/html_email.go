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

package markup

import (
	"strings"

	"golang.org/x/net/html"
)

// emailAddressProcessor replaces raw email addresses with a mailto: link.
func emailAddressProcessor(ctx *RenderContext, node *html.Node) {
	next := node.NextSibling
	for node != nil && node != next {
		m := globalVars().emailRegex.FindStringSubmatchIndex(node.Data)
		if m == nil {
			return
		}

		var nextByte byte
		if len(node.Data) > m[3] {
			nextByte = node.Data[m[3]]
		}
		if strings.IndexByte(":/", nextByte) != -1 {
			// for cases: "git@kmup.com:owner/repo.git", "https://git@kmup.com/owner/repo.git"
			return
		}
		mail := node.Data[m[2]:m[3]]
		replaceContent(node, m[2], m[3], createLink(ctx, "mailto:"+mail, mail, "" /*mailto*/))
		node = node.NextSibling.NextSibling
	}
}
