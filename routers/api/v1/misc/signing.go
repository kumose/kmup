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

package misc

import (
	"github.com/kumose/kmup/modules/git"
	asymkey_service "github.com/kumose/kmup/services/asymkey"
	"github.com/kumose/kmup/services/context"
)

func getSigningKey(ctx *context.APIContext, expectedFormat string) {
	// if the handler is in the repo's route group, get the repo's signing key
	// otherwise, get the global signing key
	path := ""
	if ctx.Repo != nil && ctx.Repo.Repository != nil {
		path = ctx.Repo.Repository.RepoPath()
	}
	content, format, err := asymkey_service.PublicSigningKey(ctx, path)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	if format == "" {
		ctx.APIErrorNotFound("no signing key")
		return
	} else if format != expectedFormat {
		ctx.APIErrorNotFound("signing key format is " + format)
		return
	}
	_, _ = ctx.Write([]byte(content))
}

// SigningKeyGPG returns the public key of the default signing key if it exists
func SigningKeyGPG(ctx *context.APIContext) {
	// swagger:operation GET /signing-key.gpg miscellaneous getSigningKey
	// ---
	// summary: Get default signing-key.gpg
	// produces:
	//     - text/plain
	// responses:
	//   "200":
	//     description: "GPG armored public key"
	//     schema:
	//       type: string

	// swagger:operation GET /repos/{owner}/{repo}/signing-key.gpg repository repoSigningKey
	// ---
	// summary: Get signing-key.gpg for given repository
	// produces:
	//     - text/plain
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
	// responses:
	//   "200":
	//     description: "GPG armored public key"
	//     schema:
	//       type: string
	getSigningKey(ctx, git.SigningKeyFormatOpenPGP)
}

// SigningKeySSH returns the public key of the default signing key if it exists
func SigningKeySSH(ctx *context.APIContext) {
	// swagger:operation GET /signing-key.pub miscellaneous getSigningKeySSH
	// ---
	// summary: Get default signing-key.pub
	// produces:
	//     - text/plain
	// responses:
	//   "200":
	//     description: "ssh public key"
	//     schema:
	//       type: string

	// swagger:operation GET /repos/{owner}/{repo}/signing-key.pub repository repoSigningKeySSH
	// ---
	// summary: Get signing-key.pub for given repository
	// produces:
	//     - text/plain
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
	// responses:
	//   "200":
	//     description: "ssh public key"
	//     schema:
	//       type: string
	getSigningKey(ctx, git.SigningKeyFormatSSH)
}
