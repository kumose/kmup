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

package npm

// TagProperty is the name of the property for tag management
const TagProperty = "npm.tag"

// Metadata represents the metadata of a npm package
type Metadata struct {
	Scope                   string            `json:"scope,omitempty"`
	Name                    string            `json:"name,omitempty"`
	Description             string            `json:"description,omitempty"`
	Author                  string            `json:"author,omitempty"`
	License                 string            `json:"license,omitempty"`
	ProjectURL              string            `json:"project_url,omitempty"`
	Keywords                []string          `json:"keywords,omitempty"`
	Dependencies            map[string]string `json:"dependencies,omitempty"`
	BundleDependencies      []string          `json:"bundleDependencies,omitempty"`
	DevelopmentDependencies map[string]string `json:"development_dependencies,omitempty"`
	PeerDependencies        map[string]string `json:"peer_dependencies,omitempty"`
	PeerDependenciesMeta    map[string]any    `json:"peer_dependencies_meta,omitempty"`
	OptionalDependencies    map[string]string `json:"optional_dependencies,omitempty"`
	Bin                     map[string]string `json:"bin,omitempty"`
	Readme                  string            `json:"readme,omitempty"`
	Repository              Repository        `json:"repository"`
}
