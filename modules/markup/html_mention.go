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

	"github.com/kumose/kmup/modules/references"
	"github.com/kumose/kmup/modules/util"

	"golang.org/x/net/html"
)

func mentionProcessor(ctx *RenderContext, node *html.Node) {
	start := 0
	nodeStop := node.NextSibling
	for node != nodeStop {
		found, loc := references.FindFirstMentionBytes(util.UnsafeStringToBytes(node.Data[start:]))
		if !found {
			node = node.NextSibling
			start = 0
			continue
		}
		loc.Start += start
		loc.End += start
		mention := node.Data[loc.Start:loc.End]
		teams, ok := ctx.RenderOptions.Metas["teams"]
		// FIXME: util.URLJoin may not be necessary here:
		// - setting.AppURL is defined to have a terminal '/' so unless mention[1:]
		// is an AppSubURL link we can probably fallback to concatenation.
		// team mention should follow @orgName/teamName style
		if ok && strings.Contains(mention, "/") {
			mentionOrgAndTeam := strings.Split(mention, "/")
			if mentionOrgAndTeam[0][1:] == ctx.RenderOptions.Metas["org"] && strings.Contains(teams, ","+strings.ToLower(mentionOrgAndTeam[1])+",") {
				link := "/:root/" + util.URLJoin("org", ctx.RenderOptions.Metas["org"], "teams", mentionOrgAndTeam[1])
				replaceContent(node, loc.Start, loc.End, createLink(ctx, link, mention, "" /*mention*/))
				node = node.NextSibling.NextSibling
				start = 0
				continue
			}
			start = loc.End
			continue
		}
		mentionedUsername := mention[1:]

		if DefaultRenderHelperFuncs != nil && DefaultRenderHelperFuncs.IsUsernameMentionable(ctx, mentionedUsername) {
			link := "/:root/" + mentionedUsername
			replaceContent(node, loc.Start, loc.End, createLink(ctx, link, mention, "" /*mention*/))
			node = node.NextSibling.NextSibling
			start = 0
		} else {
			start = loc.End
		}
	}
}
