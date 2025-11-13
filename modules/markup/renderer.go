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
	"io"
	"path"
	"strings"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/typesniffer"
)

// Renderer defines an interface for rendering markup file to HTML
type Renderer interface {
	Name() string // markup format name
	Extensions() []string
	SanitizerRules() []setting.MarkupSanitizerRule
	Render(ctx *RenderContext, input io.Reader, output io.Writer) error
}

// PostProcessRenderer defines an interface for renderers who need post process
type PostProcessRenderer interface {
	NeedPostProcess() bool
}

type ExternalRendererOptions struct {
	SanitizerDisabled bool
	DisplayInIframe   bool
	ContentSandbox    string
}

// ExternalRenderer defines an interface for external renderers
type ExternalRenderer interface {
	GetExternalRendererOptions() ExternalRendererOptions
}

// RendererContentDetector detects if the content can be rendered
// by specified renderer
type RendererContentDetector interface {
	CanRender(filename string, sniffedType typesniffer.SniffedType, prefetchBuf []byte) bool
}

var (
	extRenderers = make(map[string]Renderer)
	renderers    = make(map[string]Renderer)
)

// RegisterRenderer registers a new markup file renderer
func RegisterRenderer(renderer Renderer) {
	renderers[renderer.Name()] = renderer
	for _, ext := range renderer.Extensions() {
		extRenderers[strings.ToLower(ext)] = renderer
	}
}

// GetRendererByFileName get renderer by filename
func GetRendererByFileName(filename string) Renderer {
	extension := strings.ToLower(path.Ext(filename))
	return extRenderers[extension]
}

// DetectRendererType detects the markup type of the content
func DetectRendererType(filename string, sniffedType typesniffer.SniffedType, prefetchBuf []byte) string {
	for _, renderer := range renderers {
		if detector, ok := renderer.(RendererContentDetector); ok && detector.CanRender(filename, sniffedType, prefetchBuf) {
			return renderer.Name()
		}
	}
	return ""
}

// DetectMarkupTypeByFileName returns the possible markup format type via the filename
func DetectMarkupTypeByFileName(filename string) string {
	if parser := GetRendererByFileName(filename); parser != nil {
		return parser.Name()
	}
	return ""
}

func PreviewableExtensions() []string {
	extensions := make([]string, 0, len(extRenderers))
	for extension := range extRenderers {
		extensions = append(extensions, extension)
	}
	return extensions
}
