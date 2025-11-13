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

package repo

import (
	"context"
	"fmt"
	"image/png"
	"io"
	"net/url"
	"strconv"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/avatar"
	"github.com/kumose/kmup/modules/httplib"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/storage"
)

// CustomAvatarRelativePath returns repository custom avatar file path.
func (repo *Repository) CustomAvatarRelativePath() string {
	return repo.Avatar
}

// ExistsWithAvatarAtStoragePath returns true if there is a user with this Avatar
func ExistsWithAvatarAtStoragePath(ctx context.Context, storagePath string) (bool, error) {
	// See func (repo *Repository) CustomAvatarRelativePath()
	// repo.Avatar is used directly as the storage path - therefore we can check for existence directly using the path
	return db.GetEngine(ctx).Where("`avatar`=?", storagePath).Exist(new(Repository))
}

// RelAvatarLink returns a relative link to the repository's avatar.
func (repo *Repository) RelAvatarLink(ctx context.Context) string {
	return repo.relAvatarLink(ctx)
}

// generateRandomAvatar generates a random avatar for repository.
func generateRandomAvatar(ctx context.Context, repo *Repository) error {
	idToString := strconv.FormatInt(repo.ID, 10)

	seed := idToString
	img, err := avatar.RandomImage([]byte(seed))
	if err != nil {
		return fmt.Errorf("RandomImage: %w", err)
	}

	repo.Avatar = idToString

	if err := storage.SaveFrom(storage.RepoAvatars, repo.CustomAvatarRelativePath(), func(w io.Writer) error {
		if err := png.Encode(w, img); err != nil {
			log.Error("Encode: %v", err)
		}
		return err
	}); err != nil {
		return fmt.Errorf("Failed to create dir %s: %w", repo.CustomAvatarRelativePath(), err)
	}

	log.Info("New random avatar created for repository: %d", repo.ID)

	if _, err := db.GetEngine(ctx).ID(repo.ID).Cols("avatar").NoAutoTime().Update(repo); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) relAvatarLink(ctx context.Context) string {
	// If no avatar - path is empty
	avatarPath := repo.CustomAvatarRelativePath()
	if len(avatarPath) == 0 {
		switch mode := setting.RepoAvatar.Fallback; mode {
		case "image":
			return setting.RepoAvatar.FallbackImage
		case "random":
			if err := generateRandomAvatar(ctx, repo); err != nil {
				log.Error("generateRandomAvatar: %v", err)
			}
		default:
			// default behaviour: do not display avatar
			return ""
		}
	}
	return setting.AppSubURL + "/repo-avatars/" + url.PathEscape(repo.Avatar)
}

// AvatarLink returns the full avatar url with http host or the empty string if the repo doesn't have an avatar.
//
// TODO: refactor it to a relative URL, but it is still used in API response at the moment
func (repo *Repository) AvatarLink(ctx context.Context) string {
	relLink := repo.relAvatarLink(ctx)
	if relLink != "" {
		return httplib.MakeAbsoluteURL(ctx, relLink)
	}
	return ""
}
