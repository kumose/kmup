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
	"fmt"

	asymkey_model "github.com/kumose/kmup/models/asymkey"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/perm"
)

// AddPrincipalKey adds new principal to database and authorized_principals file.
func AddPrincipalKey(ctx context.Context, ownerID int64, content string, authSourceID int64) (*asymkey_model.PublicKey, error) {
	key := &asymkey_model.PublicKey{
		OwnerID:       ownerID,
		Name:          content,
		Content:       content,
		Mode:          perm.AccessModeWrite,
		Type:          asymkey_model.KeyTypePrincipal,
		LoginSourceID: authSourceID,
	}

	if err := db.WithTx(ctx, func(ctx context.Context) error {
		// Principals cannot be duplicated.
		has, err := db.GetEngine(ctx).
			Where("content = ? AND type = ?", content, asymkey_model.KeyTypePrincipal).
			Get(new(asymkey_model.PublicKey))
		if err != nil {
			return err
		} else if has {
			return asymkey_model.ErrKeyAlreadyExist{
				Content: content,
			}
		}

		if err = db.Insert(ctx, key); err != nil {
			return fmt.Errorf("addKey: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return key, RewriteAllPrincipalKeys(ctx)
}
