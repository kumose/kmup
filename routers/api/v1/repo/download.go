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

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/services/context"
	archiver_service "github.com/kumose/kmup/services/repository/archiver"
)

func serveRepoArchive(ctx *context.APIContext, reqFileName string) {
	aReq, err := archiver_service.NewRequest(ctx.Repo.Repository, ctx.Repo.GitRepo, reqFileName)
	if err != nil {
		if errors.Is(err, archiver_service.ErrUnknownArchiveFormat{}) {
			ctx.APIError(http.StatusBadRequest, err)
		} else if errors.Is(err, archiver_service.RepoRefNotFoundError{}) {
			ctx.APIError(http.StatusNotFound, err)
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}
	archiver_service.ServeRepoArchive(ctx.Base, aReq)
}

func DownloadArchive(ctx *context.APIContext) {
	var tp repo_model.ArchiveType
	switch ballType := ctx.PathParam("ball_type"); ballType {
	case "tarball":
		tp = repo_model.ArchiveTarGz
	case "zipball":
		tp = repo_model.ArchiveZip
	case "bundle":
		tp = repo_model.ArchiveBundle
	default:
		ctx.APIError(http.StatusBadRequest, "Unknown archive type: "+ballType)
		return
	}
	serveRepoArchive(ctx, ctx.PathParam("*")+"."+tp.String())
}
