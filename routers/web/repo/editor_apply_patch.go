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
	"strings"

	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/forms"
	"github.com/kumose/kmup/services/repository/files"
)

func NewDiffPatch(ctx *context.Context) {
	prepareEditorCommitFormOptions(ctx, "_diffpatch")
	if ctx.Written() {
		return
	}

	ctx.Data["PageIsPatch"] = true
	ctx.HTML(http.StatusOK, tplPatchFile)
}

// NewDiffPatchPost response for sending patch page
func NewDiffPatchPost(ctx *context.Context) {
	parsed := prepareEditorCommitSubmittedForm[*forms.EditRepoFileForm](ctx)
	if ctx.Written() {
		return
	}

	defaultCommitMessage := ctx.Locale.TrString("repo.editor.patch")
	_, err := files.ApplyDiffPatch(ctx, ctx.Repo.Repository, ctx.Doer, &files.ApplyDiffPatchOptions{
		LastCommitID: parsed.form.LastCommit,
		OldBranch:    parsed.OldBranchName,
		NewBranch:    parsed.NewBranchName,
		Message:      parsed.GetCommitMessage(defaultCommitMessage),
		Content:      strings.ReplaceAll(parsed.form.Content.Value(), "\r\n", "\n"),
		Author:       parsed.GitCommitter,
		Committer:    parsed.GitCommitter,
	})
	if err != nil {
		err = util.ErrorWrapTranslatable(err, "repo.editor.fail_to_apply_patch")
	}
	if err != nil {
		editorHandleFileOperationError(ctx, parsed.NewBranchName, err)
		return
	}
	redirectForCommitChoice(ctx, parsed, parsed.form.TreePath)
}
