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
	"net/http"

	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/lfs"
)

const RouterMockPointCommonLFS = "common-lfs"

func AddOwnerRepoGitLFSRoutes(m *web.Router, middlewares ...any) {
	// shared by web and internal routers
	m.Group("/{username}/{reponame}/info/lfs", func() {
		m.Post("/objects/batch", lfs.CheckAcceptMediaType, lfs.BatchHandler)
		m.Put("/objects/{oid}/{size}", lfs.UploadHandler)
		m.Get("/objects/{oid}/{filename}", lfs.DownloadHandler)
		m.Get("/objects/{oid}", lfs.DownloadHandler)
		m.Post("/verify", lfs.CheckAcceptMediaType, lfs.VerifyHandler)
		m.Group("/locks", func() {
			m.Get("/", lfs.GetListLockHandler)
			m.Post("/", lfs.PostLockHandler)
			m.Post("/verify", lfs.VerifyLockHandler)
			m.Post("/{lid}/unlock", lfs.UnLockHandler)
		}, lfs.CheckAcceptMediaType)
		m.Any("/*", http.NotFound)
	}, append([]any{web.RouterMockPoint(RouterMockPointCommonLFS)}, middlewares...)...)
}
