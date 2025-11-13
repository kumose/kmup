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
	"unicode"

	"github.com/kumose/kmup/modules/emoji"
	"github.com/kumose/kmup/modules/setting"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func createEmoji(ctx *RenderContext, content, name string) *html.Node {
	span := &html.Node{
		Type: html.ElementNode,
		Data: atom.Span.String(),
		Attr: []html.Attribute{},
	}
	span.Attr = append(span.Attr, ctx.RenderInternal.NodeSafeAttr("class", "emoji"))
	if name != "" {
		span.Attr = append(span.Attr, html.Attribute{Key: "aria-label", Val: name})
	}

	text := &html.Node{
		Type: html.TextNode,
		Data: content,
	}

	span.AppendChild(text)
	return span
}

func createCustomEmoji(ctx *RenderContext, alias string) *html.Node {
	span := &html.Node{
		Type: html.ElementNode,
		Data: atom.Span.String(),
		Attr: []html.Attribute{},
	}
	span.Attr = append(span.Attr, ctx.RenderInternal.NodeSafeAttr("class", "emoji"))
	span.Attr = append(span.Attr, html.Attribute{Key: "aria-label", Val: alias})

	img := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Img,
		Data:     "img",
		Attr:     []html.Attribute{},
	}
	img.Attr = append(img.Attr, html.Attribute{Key: "alt", Val: ":" + alias + ":"})
	img.Attr = append(img.Attr, html.Attribute{Key: "src", Val: setting.StaticURLPrefix + "/assets/img/emoji/" + alias + ".png"})

	span.AppendChild(img)
	return span
}

// emojiShortCodeProcessor for rendering text like :smile: into emoji
func emojiShortCodeProcessor(ctx *RenderContext, node *html.Node) {
	start := 0
	next := node.NextSibling
	for node != nil && node != next && start < len(node.Data) {
		m := globalVars().emojiShortCodeRegex.FindStringSubmatchIndex(node.Data[start:])
		if m == nil {
			return
		}
		m[0] += start
		m[1] += start
		start = m[1]

		alias := node.Data[m[0]:m[1]]

		var nextChar byte
		if m[1] < len(node.Data) {
			nextChar = node.Data[m[1]]
		}
		if nextChar == ':' || unicode.IsLetter(rune(nextChar)) || unicode.IsDigit(rune(nextChar)) {
			continue
		}

		alias = strings.Trim(alias, ":")
		converted := emoji.FromAlias(alias)
		if converted != nil {
			// standard emoji
			replaceContent(node, m[0], m[1], createEmoji(ctx, converted.Emoji, converted.Description))
			node = node.NextSibling.NextSibling
			start = 0 // restart searching start since node has changed
		} else if _, exist := setting.UI.CustomEmojisMap[alias]; exist {
			// custom reaction
			replaceContent(node, m[0], m[1], createCustomEmoji(ctx, alias))
			node = node.NextSibling.NextSibling
			start = 0 // restart searching start since node has changed
		}
	}
}

// emoji processor to match emoji and add emoji class
func emojiProcessor(ctx *RenderContext, node *html.Node) {
	start := 0
	next := node.NextSibling
	for node != nil && node != next && start < len(node.Data) {
		m := emoji.FindEmojiSubmatchIndex(node.Data[start:])
		if m == nil {
			return
		}
		m[0] += start
		m[1] += start

		codepoint := node.Data[m[0]:m[1]]
		start = m[1]
		val := emoji.FromCode(codepoint)
		if val != nil {
			replaceContent(node, m[0], m[1], createEmoji(ctx, codepoint, val.Description))
			node = node.NextSibling.NextSibling
			start = 0
		}
	}
}
