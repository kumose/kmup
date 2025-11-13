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

	asymkey_model "github.com/kumose/kmup/models/asymkey"
	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
)

// DeletePublicKey deletes SSH key information both in database and authorized_keys file.
func DeletePublicKey(ctx context.Context, doer *user_model.User, id int64) (err error) {
	key, err := asymkey_model.GetPublicKeyByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user has access to delete this key.
	if !doer.IsAdmin && doer.ID != key.OwnerID {
		return asymkey_model.ErrKeyAccessDenied{
			UserID: doer.ID,
			KeyID:  key.ID,
			Note:   "public",
		}
	}

	if _, err = db.DeleteByID[asymkey_model.PublicKey](ctx, id); err != nil {
		return err
	}

	if key.Type == asymkey_model.KeyTypePrincipal {
		return RewriteAllPrincipalKeys(ctx)
	}

	return RewriteAllPublicKeys(ctx)
}
