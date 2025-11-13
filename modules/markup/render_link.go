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
	"context"
	"strings"

	"github.com/kumose/kmup/modules/httplib"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

func resolveLinkRelative(ctx context.Context, base, cur, link string, absolute bool) (finalLink string) {
	if IsFullURLString(link) {
		return link
	}
	if strings.HasPrefix(link, "/") {
		if strings.HasPrefix(link, base) && strings.Count(base, "/") >= 4 {
			// a trick to tolerate that some users were using absolute paths (the old kmup's behavior)
			finalLink = link
		} else {
			finalLink = util.URLJoin(base, "./", link)
		}
	} else {
		finalLink = util.URLJoin(base, "./", cur, link)
	}
	finalLink = strings.TrimSuffix(finalLink, "/")
	if absolute {
		finalLink = httplib.MakeAbsoluteURL(ctx, finalLink)
	}
	return finalLink
}

func (ctx *RenderContext) ResolveLinkRelative(base, cur, link string) string {
	if strings.HasPrefix(link, "/:") {
		setting.PanicInDevOrTesting("invalid link %q, forgot to cut?", link)
	}
	return resolveLinkRelative(ctx, base, cur, link, ctx.RenderOptions.UseAbsoluteLink)
}

func (ctx *RenderContext) ResolveLinkRoot(link string) string {
	return ctx.ResolveLinkRelative(setting.AppSubURL+"/", "", link)
}

func ParseRenderedLink(s, preferLinkType string) (linkType, link string) {
	if strings.HasPrefix(s, "/:") {
		p := strings.IndexByte(s[1:], '/')
		if p == -1 {
			return s, ""
		}
		return s[:p+1], s[p+2:]
	}
	return preferLinkType, s
}
