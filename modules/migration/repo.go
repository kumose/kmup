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

package migration

// Repository defines a standard repository information
type Repository struct {
	Name          string
	Owner         string
	IsPrivate     bool `yaml:"is_private"`
	IsMirror      bool `yaml:"is_mirror"`
	Description   string
	CloneURL      string `yaml:"clone_url"` // SECURITY: This must be checked to ensure that is safe to be used
	OriginalURL   string `yaml:"original_url"`
	DefaultBranch string
}
