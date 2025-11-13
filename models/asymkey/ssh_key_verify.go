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

package asymkey

import (
	"context"
	"strings"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/log"

	"github.com/42wim/sshsig"
)

// VerifySSHKey marks a SSH key as verified
func VerifySSHKey(ctx context.Context, ownerID int64, fingerprint, token, signature string) (string, error) {
	return db.WithTx2(ctx, func(ctx context.Context) (string, error) {
		key := new(PublicKey)

		has, err := db.GetEngine(ctx).Where("owner_id = ? AND fingerprint = ?", ownerID, fingerprint).Get(key)
		if err != nil {
			return "", err
		} else if !has {
			return "", ErrKeyNotExist{}
		}

		err = sshsig.Verify(strings.NewReader(token), []byte(signature), []byte(key.Content), "kmup")
		if err != nil {
			// edge case for Windows based shells that will add CR LF if piped to ssh-keygen command
			// see https://github.com/PowerShell/PowerShell/issues/5974
			if sshsig.Verify(strings.NewReader(token+"\r\n"), []byte(signature), []byte(key.Content), "kmup") != nil {
				log.Debug("VerifySSHKey sshsig.Verify failed: %v", err)
				return "", ErrSSHInvalidTokenSignature{
					Fingerprint: key.Fingerprint,
				}
			}
		}

		key.Verified = true
		if _, err := db.GetEngine(ctx).ID(key.ID).Cols("verified").Update(key); err != nil {
			return "", err
		}

		return key.Fingerprint, nil
	})
}
