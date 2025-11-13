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

package conan

const (
	PropertyRecipeUser       = "conan.recipe.user"
	PropertyRecipeChannel    = "conan.recipe.channel"
	PropertyRecipeRevision   = "conan.recipe.revision"
	PropertyPackageReference = "conan.package.reference"
	PropertyPackageRevision  = "conan.package.revision"
	PropertyPackageInfo      = "conan.package.info"
)

// Metadata represents the metadata of a Conan package
type Metadata struct {
	Author        string   `json:"author,omitempty"`
	License       string   `json:"license,omitempty"`
	ProjectURL    string   `json:"project_url,omitempty"`
	RepositoryURL string   `json:"repository_url,omitempty"`
	Description   string   `json:"description,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
}
