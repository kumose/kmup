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

package helper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	packages_model "github.com/kumose/kmup/models/packages"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/context"
)

// ProcessErrorForUser logs the error and returns a user-error message for the end user.
// If the status is http.StatusInternalServerError, the message is stripped for non-admin users in production.
func ProcessErrorForUser(ctx *context.Context, status int, errObj any) string {
	var message string
	if err, ok := errObj.(error); ok {
		message = err.Error()
	} else if errObj != nil {
		message = fmt.Sprint(errObj)
	}

	if status == http.StatusInternalServerError {
		log.Log(2, log.ERROR, "Package registry API internal error: %d %s", status, message)
		if setting.IsProd && (ctx.Doer == nil || !ctx.Doer.IsAdmin) {
			message = "internal server error"
		}
		return message
	}

	log.Log(2, log.DEBUG, "Package registry API user error: %d %s", status, message)
	return message
}

// ServePackageFile the content of the package file
// If the url is set it will redirect the request, otherwise the content is copied to the response.
func ServePackageFile(ctx *context.Context, s io.ReadSeekCloser, u *url.URL, pf *packages_model.PackageFile, forceOpts ...*context.ServeHeaderOptions) {
	if u != nil {
		ctx.Redirect(u.String())
		return
	}

	defer s.Close()

	var opts *context.ServeHeaderOptions
	if len(forceOpts) > 0 {
		opts = forceOpts[0]
	} else {
		opts = &context.ServeHeaderOptions{
			Filename:     pf.Name,
			LastModified: pf.CreatedUnix.AsLocalTime(),
		}
	}

	ctx.ServeContent(s, opts)
}
