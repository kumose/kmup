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
	"html/template"

	"github.com/kumose/kmup/modules/setting"
)

const (
	LinkTypeDefault = ""
	LinkTypeRoot    = "/:root"  // the link is relative to the AppSubURL(ROOT_URL)
	LinkTypeMedia   = "/:media" // the link should be used to access media files (images, videos)
	LinkTypeRaw     = "/:raw"   // not really useful, mainly for environment KMUP_PREFIX_RAW for external renders
)

type RenderHelper interface {
	CleanUp()

	// TODO: such dependency is not ideal. We should decouple the processors step by step.
	// It should make the render choose different processors for different purposes,
	// but not make processors to guess "is it rendering a comment or a wiki?" or "does it need to check commit ID?"

	IsCommitIDExisting(commitID string) bool
	ResolveLink(link, preferLinkType string) string
}

// RenderHelperFuncs is used to decouple cycle-import
// At the moment there are different packages:
// modules/markup: basic markup rendering
// models/renderhelper: need to access models and git repo, and models/issues needs it
// services/markup: some real helper functions could only be provided here because it needs to access various services & templates
type RenderHelperFuncs struct {
	IsUsernameMentionable     func(ctx context.Context, username string) bool
	RenderRepoFileCodePreview func(ctx context.Context, options RenderCodePreviewOptions) (template.HTML, error)
	RenderRepoIssueIconTitle  func(ctx context.Context, options RenderIssueIconTitleOptions) (template.HTML, error)
}

var DefaultRenderHelperFuncs *RenderHelperFuncs

type SimpleRenderHelper struct{}

func (r *SimpleRenderHelper) CleanUp() {}

func (r *SimpleRenderHelper) IsCommitIDExisting(commitID string) bool {
	return false
}

func (r *SimpleRenderHelper) ResolveLink(link, preferLinkType string) string {
	_, link = ParseRenderedLink(link, preferLinkType)
	return resolveLinkRelative(context.Background(), setting.AppSubURL+"/", "", link, false)
}

var _ RenderHelper = (*SimpleRenderHelper)(nil)
