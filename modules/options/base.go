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

package options

import (
	"github.com/kumose/kmup/modules/assetfs"
	"github.com/kumose/kmup/modules/setting"
)

func CustomAssets() *assetfs.Layer {
	return assetfs.Local("custom", setting.CustomPath, "options")
}

func AssetFS() *assetfs.LayeredFS {
	return assetfs.Layered(CustomAssets(), BuiltinAssets())
}

// Locale reads the content of a specific locale from static/bindata or custom path.
func Locale(name string) ([]byte, error) {
	return AssetFS().ReadFile("locale", name)
}

// Readme reads the content of a specific readme from static/bindata or custom path.
func Readme(name string) ([]byte, error) {
	return AssetFS().ReadFile("readme", name)
}

// Gitignore reads the content of a gitignore locale from static/bindata or custom path.
func Gitignore(name string) ([]byte, error) {
	return AssetFS().ReadFile("gitignore", name)
}

// License reads the content of a specific license from static/bindata or custom path.
func License(name string) ([]byte, error) {
	return AssetFS().ReadFile("license", name)
}

// Labels reads the content of a specific labels from static/bindata or custom path.
func Labels(name string) ([]byte, error) {
	return AssetFS().ReadFile("label", name)
}
