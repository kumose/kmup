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
	"errors"
	"net/http"

	"github.com/kumose/kmup/modules/git"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// GetNote Get a note corresponding to a single commit from a repository
func GetNote(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/git/notes/{sha} repository repoGetNote
	// ---
	// summary: Get a note corresponding to a single commit from a repository
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
	// - name: sha
	//   in: path
	//   description: a git ref or commit sha
	//   type: string
	//   required: true
	// - name: verification
	//   in: query
	//   description: include verification for every commit (disable for speedup, default 'true')
	//   type: boolean
	// - name: files
	//   in: query
	//   description: include a list of affected files for every commit (disable for speedup, default 'true')
	//   type: boolean
	// responses:
	//   "200":
	//     "$ref": "#/responses/Note"
	//   "422":
	//     "$ref": "#/responses/validationError"
	//   "404":
	//     "$ref": "#/responses/notFound"

	sha := ctx.PathParam("sha")
	if !git.IsValidRefPattern(sha) {
		ctx.APIError(http.StatusUnprocessableEntity, "no valid ref or sha: "+sha)
		return
	}
	getNote(ctx, sha)
}

func getNote(ctx *context.APIContext, identifier string) {
	if ctx.Repo.GitRepo == nil {
		ctx.APIErrorInternal(errors.New("no open git repo"))
		return
	}

	commitID, err := ctx.Repo.GitRepo.ConvertToGitID(identifier)
	if err != nil {
		if git.IsErrNotExist(err) {
			ctx.APIErrorNotFound(err)
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}

	var note git.Note
	if err := git.GetNote(ctx, ctx.Repo.GitRepo, commitID.String(), &note); err != nil {
		if git.IsErrNotExist(err) {
			ctx.APIErrorNotFound("commit doesn't exist: " + identifier)
			return
		}
		ctx.APIErrorInternal(err)
		return
	}

	verification := ctx.FormString("verification") == "" || ctx.FormBool("verification")
	files := ctx.FormString("files") == "" || ctx.FormBool("files")

	cmt, err := convert.ToCommit(ctx, ctx.Repo.Repository, ctx.Repo.GitRepo, note.Commit, nil,
		convert.ToCommitOptions{
			Stat:         true,
			Verification: verification,
			Files:        files,
		})
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	apiNote := api.Note{Message: string(note.Message), Commit: cmt}
	ctx.JSON(http.StatusOK, apiNote)
}
