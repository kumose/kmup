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

package forms

import (
	"net/http"

	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/web/middleware"
	"github.com/kumose/kmup/services/context"

	"github.com/kumose-go/chi/binding"
)

type CommitCommonForm struct {
	TreePath      string `binding:"MaxSize(500)"`
	CommitSummary string `binding:"MaxSize(100)"`
	CommitMessage string
	CommitChoice  string `binding:"Required;MaxSize(50)"`
	NewBranchName string `binding:"GitRefName;MaxSize(100)"`
	LastCommit    string
	Signoff       bool
	CommitEmail   string
}

func (f *CommitCommonForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

type CommitCommonFormInterface interface {
	GetCommitCommonForm() *CommitCommonForm
}

func (f *CommitCommonForm) GetCommitCommonForm() *CommitCommonForm {
	return f
}

type EditRepoFileForm struct {
	CommitCommonForm
	Content optional.Option[string]
}

type DeleteRepoFileForm struct {
	CommitCommonForm
}

type UploadRepoFileForm struct {
	CommitCommonForm
	Files []string
}

type CherryPickForm struct {
	CommitCommonForm
	Revert bool
}
