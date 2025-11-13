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

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/context/upload"
	files_service "github.com/kumose/kmup/services/repository/files"
)

// UploadFileToServer upload file to server file dir not git
func UploadFileToServer(ctx *context.Context) {
	file, header, err := ctx.Req.FormFile("file")
	if err != nil {
		ctx.ServerError("FormFile", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, _ := util.ReadAtMost(file, buf)
	if n > 0 {
		buf = buf[:n]
	}

	err = upload.Verify(buf, header.Filename, setting.Repository.Upload.AllowedTypes)
	if err != nil {
		ctx.HTTPError(http.StatusBadRequest, err.Error())
		return
	}

	name := files_service.CleanGitTreePath(header.Filename)
	if len(name) == 0 {
		ctx.HTTPError(http.StatusBadRequest, "Upload file name is invalid")
		return
	}

	// FIXME: need to check the file size according to setting.Repository.Upload.FileMaxSize

	uploaded, err := repo_model.NewUpload(ctx, name, buf, file)
	if err != nil {
		ctx.ServerError("NewUpload", err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"uuid": uploaded.UUID})
}

// RemoveUploadFileFromServer remove file from server file dir
func RemoveUploadFileFromServer(ctx *context.Context) {
	fileUUID := ctx.FormString("file")
	if err := repo_model.DeleteUploadByUUID(ctx, fileUUID); err != nil {
		ctx.ServerError("DeleteUploadByUUID", err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
