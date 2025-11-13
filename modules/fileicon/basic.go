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

package fileicon

import (
	"html/template"

	"github.com/kumose/kmup/modules/svg"
	"github.com/kumose/kmup/modules/util"
)

func BasicEntryIconName(entry *EntryInfo) string {
	svgName := "octicon-file"
	switch {
	case entry.EntryMode.IsLink():
		svgName = "octicon-file-symlink-file"
		if entry.SymlinkToMode.IsDir() {
			svgName = "octicon-file-directory-symlink"
		}
	case entry.EntryMode.IsDir():
		svgName = util.Iif(entry.IsOpen, "octicon-file-directory-open-fill", "octicon-file-directory-fill")
	case entry.EntryMode.IsSubModule():
		svgName = "octicon-file-submodule"
	}
	return svgName
}

func BasicEntryIconHTML(entry *EntryInfo) template.HTML {
	return svg.RenderHTML(BasicEntryIconName(entry))
}
