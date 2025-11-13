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

package actions

import (
	"net/http"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/storage"
	"github.com/kumose/kmup/services/context"
)

// Artifacts using the v4 backend are stored as a single combined zip file per artifact on the backend
// The v4 backend ensures ContentEncoding is set to "application/zip", which is not the case for the old backend
func IsArtifactV4(art *actions_model.ActionArtifact) bool {
	return art.ArtifactName+".zip" == art.ArtifactPath && art.ContentEncoding == "application/zip"
}

func DownloadArtifactV4ServeDirectOnly(ctx *context.Base, art *actions_model.ActionArtifact) (bool, error) {
	if setting.Actions.ArtifactStorage.ServeDirect() {
		u, err := storage.ActionsArtifacts.URL(art.StoragePath, art.ArtifactPath, ctx.Req.Method, nil)
		if u != nil && err == nil {
			ctx.Redirect(u.String(), http.StatusFound)
			return true, nil
		}
	}
	return false, nil
}

func DownloadArtifactV4Fallback(ctx *context.Base, art *actions_model.ActionArtifact) error {
	f, err := storage.ActionsArtifacts.Open(art.StoragePath)
	if err != nil {
		return err
	}
	defer f.Close()
	http.ServeContent(ctx.Resp, ctx.Req, art.ArtifactName+".zip", art.CreatedUnix.AsLocalTime(), f)
	return nil
}

func DownloadArtifactV4(ctx *context.Base, art *actions_model.ActionArtifact) error {
	ok, err := DownloadArtifactV4ServeDirectOnly(ctx, art)
	if ok || err != nil {
		return err
	}
	return DownloadArtifactV4Fallback(ctx, art)
}
