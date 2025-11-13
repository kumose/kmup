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

	"github.com/kumose/kmup/modules/httplib"
)

// FetchRedirectDelegate helps the "fetch" requests to redirect to the correct location
func FetchRedirectDelegate(resp http.ResponseWriter, req *http.Request) {
	// When use "fetch" to post requests and the response is a redirect, browser's "location.href = uri" has limitations.
	// 1. change "location" from old "/foo" to new "/foo#hash", the browser will not reload the page.
	// 2. when use "window.reload()", the hash is not respected, the newly loaded page won't scroll to the hash target.
	// The typical page is "issue comment" page. The backend responds "/owner/repo/issues/1#comment-2",
	// then frontend needs this delegate to redirect to the new location with hash correctly.
	redirect := req.PostFormValue("redirect")
	if !httplib.IsCurrentKmupSiteURL(req.Context(), redirect) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Add("Location", redirect)
	resp.WriteHeader(http.StatusSeeOther)
}
