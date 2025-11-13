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

package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/avatar"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/storage"
)

// UploadAvatar saves custom avatar for user.
func UploadAvatar(ctx context.Context, u *user_model.User, data []byte) error {
	avatarData, err := avatar.ProcessAvatarImage(data)
	if err != nil {
		return fmt.Errorf("UploadAvatar: failed to process user avatar image: %w", err)
	}

	return db.WithTx(ctx, func(ctx context.Context) error {
		u.UseCustomAvatar = true
		u.Avatar = avatar.HashAvatar(u.ID, data)
		if err = user_model.UpdateUserCols(ctx, u, "use_custom_avatar", "avatar"); err != nil {
			return fmt.Errorf("UploadAvatar: failed to update user avatar: %w", err)
		}

		if err := storage.SaveFrom(storage.Avatars, u.CustomAvatarRelativePath(), func(w io.Writer) error {
			_, err := w.Write(avatarData)
			return err
		}); err != nil {
			return fmt.Errorf("UploadAvatar: failed to save user avatar %s: %w", u.CustomAvatarRelativePath(), err)
		}

		return nil
	})
}

// DeleteAvatar deletes the user's custom avatar.
func DeleteAvatar(ctx context.Context, u *user_model.User) error {
	aPath := u.CustomAvatarRelativePath()
	log.Trace("DeleteAvatar[%d]: %s", u.ID, aPath)

	return db.WithTx(ctx, func(ctx context.Context) error {
		hasAvatar := len(u.Avatar) > 0
		u.UseCustomAvatar = false
		u.Avatar = ""
		if _, err := db.GetEngine(ctx).ID(u.ID).Cols("avatar, use_custom_avatar").Update(u); err != nil {
			return fmt.Errorf("DeleteAvatar: %w", err)
		}

		if hasAvatar {
			if err := storage.Avatars.Delete(aPath); err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("failed to remove %s: %w", aPath, err)
				}
				log.Warn("Deleting avatar %s but it doesn't exist", aPath)
			}
		}

		return nil
	})
}
