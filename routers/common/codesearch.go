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

package common

import (
	"github.com/kumose/kmup/modules/indexer"
	code_indexer "github.com/kumose/kmup/modules/indexer/code"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/context"
)

func PrepareCodeSearch(ctx *context.Context) (ret struct {
	Keyword    string
	Language   string
	SearchMode indexer.SearchModeType
},
) {
	ret.Language = ctx.FormTrim("l")
	ret.Keyword = ctx.FormTrim("q")
	ret.SearchMode = indexer.SearchModeType(ctx.FormTrim("search_mode"))

	ctx.Data["Keyword"] = ret.Keyword
	ctx.Data["Language"] = ret.Language
	ctx.Data["SelectedSearchMode"] = string(ret.SearchMode)
	if setting.Indexer.RepoIndexerEnabled {
		ctx.Data["SearchModes"] = code_indexer.SupportedSearchModes()
	} else {
		ctx.Data["SearchModes"] = indexer.GitGrepSupportedSearchModes()
	}
	ctx.Data["IsRepoIndexerEnabled"] = setting.Indexer.RepoIndexerEnabled
	return ret
}
