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

package web

import (
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/routers/web/repo"
	"github.com/kumose/kmup/services/context"
)

func addOwnerRepoGitHTTPRouters(m *web.Router) {
	m.Group("/{username}/{reponame}", func() {
		m.Methods("POST,OPTIONS", "/git-upload-pack", repo.ServiceUploadPack)
		m.Methods("POST,OPTIONS", "/git-receive-pack", repo.ServiceReceivePack)
		m.Methods("GET,OPTIONS", "/info/refs", repo.GetInfoRefs)
		m.Methods("GET,OPTIONS", "/HEAD", repo.GetTextFile("HEAD"))
		m.Methods("GET,OPTIONS", "/objects/info/alternates", repo.GetTextFile("objects/info/alternates"))
		m.Methods("GET,OPTIONS", "/objects/info/http-alternates", repo.GetTextFile("objects/info/http-alternates"))
		m.Methods("GET,OPTIONS", "/objects/info/packs", repo.GetInfoPacks)
		m.Methods("GET,OPTIONS", "/objects/info/{file:[^/]*}", repo.GetTextFile(""))
		m.Methods("GET,OPTIONS", "/objects/{head:[0-9a-f]{2}}/{hash:[0-9a-f]{38,62}}", repo.GetLooseObject)
		m.Methods("GET,OPTIONS", "/objects/pack/pack-{file:[0-9a-f]{40,64}}.pack", repo.GetPackFile)
		m.Methods("GET,OPTIONS", "/objects/pack/pack-{file:[0-9a-f]{40,64}}.idx", repo.GetIdxFile)
	}, optSignInIgnoreCsrf, repo.HTTPGitEnabledHandler, repo.CorsHandler(), context.UserAssignmentWeb())
}
