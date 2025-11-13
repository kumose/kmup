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

package templates

import (
	"context"
	"html"
	"html/template"
	"strconv"

	activities_model "github.com/kumose/kmup/models/activities"
	"github.com/kumose/kmup/models/avatars"
	"github.com/kumose/kmup/models/organization"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	kmup_html "github.com/kumose/kmup/modules/htmlutil"
	"github.com/kumose/kmup/modules/setting"
)

type AvatarUtils struct {
	ctx context.Context
}

func NewAvatarUtils(ctx context.Context) *AvatarUtils {
	return &AvatarUtils{ctx: ctx}
}

// AvatarHTML creates the HTML for an avatar
func AvatarHTML(src string, size int, class, name string) template.HTML {
	sizeStr := strconv.Itoa(size)

	if name == "" {
		name = "avatar"
	}

	// use empty alt, otherwise if the image fails to load, the width will follow the "alt" text's width
	return template.HTML(`<img loading="lazy" alt class="` + class + `" src="` + src + `" title="` + html.EscapeString(name) + `" width="` + sizeStr + `" height="` + sizeStr + `"/>`)
}

// Avatar renders user avatars. args: user, size (int), class (string)
func (au *AvatarUtils) Avatar(item any, others ...any) template.HTML {
	size, class := kmup_html.ParseSizeAndClass(avatars.DefaultAvatarPixelSize, avatars.DefaultAvatarClass, others...)

	switch t := item.(type) {
	case *user_model.User:
		src := t.AvatarLinkWithSize(au.ctx, size*setting.Avatar.RenderedSizeFactor)
		if src != "" {
			return AvatarHTML(src, size, class, t.DisplayName())
		}
	case *repo_model.Collaborator:
		src := t.AvatarLinkWithSize(au.ctx, size*setting.Avatar.RenderedSizeFactor)
		if src != "" {
			return AvatarHTML(src, size, class, t.DisplayName())
		}
	case *organization.Organization:
		src := t.AsUser().AvatarLinkWithSize(au.ctx, size*setting.Avatar.RenderedSizeFactor)
		if src != "" {
			return AvatarHTML(src, size, class, t.AsUser().DisplayName())
		}
	}

	return AvatarHTML(avatars.DefaultAvatarLink(), size, class, "")
}

// AvatarByAction renders user avatars from action. args: action, size (int), class (string)
func (au *AvatarUtils) AvatarByAction(action *activities_model.Action, others ...any) template.HTML {
	action.LoadActUser(au.ctx)
	return au.Avatar(action.ActUser, others...)
}

// AvatarByEmail renders avatars by email address. args: email, name, size (int), class (string)
func (au *AvatarUtils) AvatarByEmail(email, name string, others ...any) template.HTML {
	size, class := kmup_html.ParseSizeAndClass(avatars.DefaultAvatarPixelSize, avatars.DefaultAvatarClass, others...)
	src := avatars.GenerateEmailAvatarFastLink(au.ctx, email, size*setting.Avatar.RenderedSizeFactor)

	if src != "" {
		return AvatarHTML(src, size, class, name)
	}

	return ""
}
