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

package repo

import (
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/services/context"
)

func HandleGitError(ctx *context.Context, msg string, err error) {
	if git.IsErrNotExist(err) {
		ctx.Data["NotFoundPrompt"] = ctx.Locale.Tr("repo.tree_path_not_found", ctx.Repo.TreePath, ctx.Repo.RefTypeNameSubURL())
		ctx.Data["NotFoundGoBackURL"] = ctx.Repo.RepoLink + "/src/" + ctx.Repo.RefTypeNameSubURL()
		ctx.NotFound(err)
	} else {
		ctx.ServerError(msg, err)
	}
}
