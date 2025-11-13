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
	"net/http"

	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/repository/files"
)

// ApplyDiffPatch handles API call for applying a patch
func ApplyDiffPatch(ctx *context.APIContext) {
	// swagger:operation POST /repos/{owner}/{repo}/diffpatch repository repoApplyDiffPatch
	// ---
	// summary: Apply diff patch to repository
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ApplyDiffPatchFileOptions"
	// responses:
	//   "200":
	//     "$ref": "#/responses/FileResponse"
	//   "404":
	//     "$ref": "#/responses/notFound"
	//   "423":
	//     "$ref": "#/responses/repoArchivedError"
	apiOpts, changeRepoFileOpts := getAPIChangeRepoFileOptions[*api.ApplyDiffPatchFileOptions](ctx)
	opts := &files.ApplyDiffPatchOptions{
		Content: apiOpts.Content,
		Message: util.IfZero(apiOpts.Message, "apply-patch"),

		OldBranch: changeRepoFileOpts.OldBranch,
		NewBranch: changeRepoFileOpts.NewBranch,
		Committer: changeRepoFileOpts.Committer,
		Author:    changeRepoFileOpts.Author,
		Dates:     changeRepoFileOpts.Dates,
		Signoff:   changeRepoFileOpts.Signoff,
	}

	fileResponse, err := files.ApplyDiffPatch(ctx, ctx.Repo.Repository, ctx.Doer, opts)
	if err != nil {
		handleChangeRepoFilesError(ctx, err)
	} else {
		ctx.JSON(http.StatusCreated, fileResponse)
	}
}
