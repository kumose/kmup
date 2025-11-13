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

package oauth2

import (
	"html/template"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/svg"
)

// BaseProvider represents a common base for Provider
type BaseProvider struct {
	name        string
	displayName string

	// TODO: maybe some providers also support SSH public keys, then they can set this to true
	supportSSHPublicKey bool
}

func (b *BaseProvider) SupportSSHPublicKey() bool {
	return b.supportSSHPublicKey
}

// Name provides the technical name for this provider
func (b *BaseProvider) Name() string {
	return b.name
}

// DisplayName returns the friendly name for this provider
func (b *BaseProvider) DisplayName() string {
	return b.displayName
}

// IconHTML returns icon HTML for this provider
func (b *BaseProvider) IconHTML(size int) template.HTML {
	svgName := "kmup-" + b.name
	switch b.name {
	case "gplus":
		svgName = "kmup-google"
	case "github":
		svgName = "octicon-mark-github"
	}
	svgHTML := svg.RenderHTML(svgName, size, "tw-mr-2")
	if svgHTML == "" {
		log.Error("No SVG icon for oauth2 provider %q", b.name)
		svgHTML = svg.RenderHTML("kmup-openid", size, "tw-mr-2")
	}
	return svgHTML
}

// CustomURLSettings returns the custom url settings for this provider
func (b *BaseProvider) CustomURLSettings() *CustomURLSettings {
	return nil
}

var _ Provider = &BaseProvider{}
