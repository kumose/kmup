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
	"fmt"
	"image/png"
	"io"

	"github.com/kumose/kmup/models/avatars"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/avatar"
	"github.com/kumose/kmup/modules/httplib"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/storage"
)

// CustomAvatarRelativePath returns user custom avatar relative path.
func (u *User) CustomAvatarRelativePath() string {
	return u.Avatar
}

// GenerateRandomAvatar generates a random avatar for user.
func GenerateRandomAvatar(ctx context.Context, u *User) error {
	seed := u.Email
	if len(seed) == 0 {
		seed = u.Name
	}

	img, err := avatar.RandomImage([]byte(seed))
	if err != nil {
		return fmt.Errorf("RandomImage: %w", err)
	}

	u.Avatar = avatars.HashEmail(seed)

	_, err = storage.Avatars.Stat(u.CustomAvatarRelativePath())
	if err != nil {
		// If unable to Stat the avatar file (usually it means non-existing), then try to save a new one
		// Don't share the images so that we can delete them easily
		if err := storage.SaveFrom(storage.Avatars, u.CustomAvatarRelativePath(), func(w io.Writer) error {
			if err := png.Encode(w, img); err != nil {
				log.Error("Encode: %v", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to save avatar %s: %w", u.CustomAvatarRelativePath(), err)
		}
	}

	if _, err := db.GetEngine(ctx).ID(u.ID).Cols("avatar").Update(u); err != nil {
		return err
	}

	return nil
}

// AvatarLinkWithSize returns a link to the user's avatar with size. size <= 0 means default size
func (u *User) AvatarLinkWithSize(ctx context.Context, size int) string {
	// ghost user was deleted, Kmup actions is a bot user, 0 means the user should be a virtual user
	// which comes from git configure information
	if u.IsGhost() || u.IsKmupActions() || u.ID <= 0 {
		return avatars.DefaultAvatarLink()
	}

	useLocalAvatar := false
	autoGenerateAvatar := false

	disableGravatar := setting.Config().Picture.DisableGravatar.Value(ctx)

	switch {
	case u.UseCustomAvatar:
		useLocalAvatar = true
	case disableGravatar, setting.OfflineMode:
		useLocalAvatar = true
		autoGenerateAvatar = true
	}

	if useLocalAvatar {
		if u.Avatar == "" && autoGenerateAvatar {
			if err := GenerateRandomAvatar(ctx, u); err != nil {
				log.Error("GenerateRandomAvatar: %v", err)
			}
		}
		if u.Avatar == "" {
			return avatars.DefaultAvatarLink()
		}
		return avatars.GenerateUserAvatarImageLink(u.Avatar, size)
	}
	return avatars.GenerateEmailAvatarFastLink(ctx, u.AvatarEmail, size)
}

// AvatarLink returns the full avatar url with http host.
// TODO: refactor it to a relative URL, but it is still used in API response at the moment
func (u *User) AvatarLink(ctx context.Context) string {
	relLink := u.AvatarLinkWithSize(ctx, 0) // it can't be empty
	return httplib.MakeAbsoluteURL(ctx, relLink)
}

// IsUploadAvatarChanged returns true if the current user's avatar would be changed with the provided data
func (u *User) IsUploadAvatarChanged(data []byte) bool {
	if !u.UseCustomAvatar || len(u.Avatar) == 0 {
		return true
	}
	avatarID := avatar.HashAvatar(u.ID, data)
	return u.Avatar != avatarID
}

// ExistsWithAvatarAtStoragePath returns true if there is a user with this Avatar
func ExistsWithAvatarAtStoragePath(ctx context.Context, storagePath string) (bool, error) {
	// See func (u *User) CustomAvatarRelativePath()
	// u.Avatar is used directly as the storage path - therefore we can check for existence directly using the path
	return db.GetEngine(ctx).Where("`avatar`=?", storagePath).Exist(new(User))
}
