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
	"slices"
	"strings"

	"github.com/kumose/kmup/modules/assetfs"
	"github.com/kumose/kmup/modules/setting"
)

func AssetFS() *assetfs.LayeredFS {
	return assetfs.Layered(CustomAssets(), BuiltinAssets())
}

func CustomAssets() *assetfs.Layer {
	return assetfs.Local("custom", setting.CustomPath, "templates")
}

func ListWebTemplateAssetNames(assets *assetfs.LayeredFS) ([]string, error) {
	files, err := assets.ListAllFiles(".", true)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(files, func(file string) bool {
		return strings.HasPrefix(file, "mail/") || !strings.HasSuffix(file, ".tmpl")
	}), nil
}

func ListMailTemplateAssetNames(assets *assetfs.LayeredFS) ([]string, error) {
	files, err := assets.ListAllFiles(".", true)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(files, func(file string) bool {
		return !strings.HasPrefix(file, "mail/") || !strings.HasSuffix(file, ".tmpl")
	}), nil
}
