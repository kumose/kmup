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

package renderhelper

import (
	"context"

	"github.com/kumose/kmup/modules/markup"
)

type SimpleDocument struct {
	*markup.SimpleRenderHelper
	ctx      *markup.RenderContext
	baseLink string
}

func (r *SimpleDocument) ResolveLink(link, preferLinkType string) string {
	linkType, link := markup.ParseRenderedLink(link, preferLinkType)
	switch linkType {
	case markup.LinkTypeRoot:
		return r.ctx.ResolveLinkRoot(link)
	default:
		return r.ctx.ResolveLinkRelative(r.baseLink, "", link)
	}
}

var _ markup.RenderHelper = (*SimpleDocument)(nil)

func NewRenderContextSimpleDocument(ctx context.Context, baseLink string) *markup.RenderContext {
	helper := &SimpleDocument{baseLink: baseLink}
	rctx := markup.NewRenderContext(ctx).WithHelper(helper).WithMetas(markup.ComposeSimpleDocumentMetas())
	helper.ctx = rctx
	return rctx
}
