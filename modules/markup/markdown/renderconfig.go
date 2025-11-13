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

package markdown

import (
	"fmt"
	"strings"

	"github.com/kumose/kmup/modules/markup"

	"github.com/yuin/goldmark/ast"
	"gopkg.in/yaml.v3"
)

// RenderConfig represents rendering configuration for this file
type RenderConfig struct {
	Meta     markup.RenderMetaMode
	TOC      string // "false": hide,  "side"/empty: in sidebar,  "main"/"true": in main view
	Lang     string
	yamlNode *yaml.Node

	// Used internally.  Cannot be controlled by frontmatter.
	metaLength int
}

func renderMetaModeFromString(s string) markup.RenderMetaMode {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "none":
		return markup.RenderMetaAsNone
	case "table":
		return markup.RenderMetaAsTable
	default: // "details"
		return markup.RenderMetaAsDetails
	}
}

// UnmarshalYAML implement yaml.v3 UnmarshalYAML
func (rc *RenderConfig) UnmarshalYAML(value *yaml.Node) error {
	if rc == nil {
		return nil
	}

	rc.yamlNode = value

	type commonRenderConfig struct {
		TOC  string `yaml:"include_toc"`
		Lang string `yaml:"lang"`
	}
	var basic commonRenderConfig
	if err := value.Decode(&basic); err != nil {
		return fmt.Errorf("unable to decode into commonRenderConfig %w", err)
	}

	if basic.Lang != "" {
		rc.Lang = basic.Lang
	}

	rc.TOC = basic.TOC

	type controlStringRenderConfig struct {
		Kmup string `yaml:"kmup"`
	}

	var stringBasic controlStringRenderConfig

	if err := value.Decode(&stringBasic); err == nil {
		if stringBasic.Kmup != "" {
			rc.Meta = renderMetaModeFromString(stringBasic.Kmup)
		}
		return nil
	}

	type yamlRenderConfig struct {
		Meta *string `yaml:"meta"`
		Icon *string `yaml:"details_icon"` // deprecated, because there is no font icon, so no custom icon
		TOC  *string `yaml:"include_toc"`
		Lang *string `yaml:"lang"`
	}

	type yamlRenderConfigWrapper struct {
		Kmup *yamlRenderConfig `yaml:"kmup"`
	}

	var cfg yamlRenderConfigWrapper
	if err := value.Decode(&cfg); err != nil {
		return fmt.Errorf("unable to decode into yamlRenderConfigWrapper %w", err)
	}

	if cfg.Kmup == nil {
		return nil
	}

	if cfg.Kmup.Meta != nil {
		rc.Meta = renderMetaModeFromString(*cfg.Kmup.Meta)
	}

	if cfg.Kmup.Lang != nil && *cfg.Kmup.Lang != "" {
		rc.Lang = *cfg.Kmup.Lang
	}

	if cfg.Kmup.TOC != nil {
		rc.TOC = *cfg.Kmup.TOC
	}

	return nil
}

func (rc *RenderConfig) toMetaNode(g *ASTTransformer) ast.Node {
	if rc.yamlNode == nil {
		return nil
	}
	switch rc.Meta {
	case markup.RenderMetaAsTable:
		return nodeToTable(rc.yamlNode)
	case markup.RenderMetaAsDetails:
		return nodeToDetails(g, rc.yamlNode)
	default:
		return nil
	}
}
