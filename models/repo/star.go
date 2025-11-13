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

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/timeutil"
)

// Star represents a starred repo by an user.
type Star struct {
	ID          int64              `xorm:"pk autoincr"`
	UID         int64              `xorm:"UNIQUE(s)"`
	RepoID      int64              `xorm:"UNIQUE(s)"`
	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
}

func init() {
	db.RegisterModel(new(Star))
}

// StarRepo or unstar repository.
func StarRepo(ctx context.Context, doer *user_model.User, repo *Repository, star bool) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		staring := IsStaring(ctx, doer.ID, repo.ID)

		if star {
			if user_model.IsUserBlockedBy(ctx, doer, repo.OwnerID) {
				return user_model.ErrBlockedUser
			}

			if staring {
				return nil
			}

			if err := db.Insert(ctx, &Star{UID: doer.ID, RepoID: repo.ID}); err != nil {
				return err
			}
			if _, err := db.Exec(ctx, "UPDATE `repository` SET num_stars = num_stars + 1 WHERE id = ?", repo.ID); err != nil {
				return err
			}
			if _, err := db.Exec(ctx, "UPDATE `user` SET num_stars = num_stars + 1 WHERE id = ?", doer.ID); err != nil {
				return err
			}
		} else {
			if !staring {
				return nil
			}

			if _, err := db.DeleteByBean(ctx, &Star{UID: doer.ID, RepoID: repo.ID}); err != nil {
				return err
			}
			if _, err := db.Exec(ctx, "UPDATE `repository` SET num_stars = num_stars - 1 WHERE id = ?", repo.ID); err != nil {
				return err
			}
			if _, err := db.Exec(ctx, "UPDATE `user` SET num_stars = num_stars - 1 WHERE id = ?", doer.ID); err != nil {
				return err
			}
		}

		return nil
	})
}

// IsStaring checks if user has starred given repository.
func IsStaring(ctx context.Context, userID, repoID int64) bool {
	has, _ := db.GetEngine(ctx).Get(&Star{UID: userID, RepoID: repoID})
	return has
}

// GetStargazers returns the users that starred the repo.
func GetStargazers(ctx context.Context, repo *Repository, opts db.ListOptions) ([]*user_model.User, error) {
	sess := db.GetEngine(ctx).Where("star.repo_id = ?", repo.ID).
		Join("LEFT", "star", "`user`.id = star.uid")
	if opts.Page > 0 {
		sess = db.SetSessionPagination(sess, &opts)

		users := make([]*user_model.User, 0, opts.PageSize)
		return users, sess.Find(&users)
	}

	users := make([]*user_model.User, 0, 8)
	return users, sess.Find(&users)
}

// ClearRepoStars clears all stars for a repository and from the user that starred it.
// Used when a repository is set to private.
func ClearRepoStars(ctx context.Context, repoID int64) error {
	if _, err := db.Exec(ctx, "UPDATE `user` SET num_stars=num_stars-1 WHERE id IN (SELECT `uid` FROM `star` WHERE repo_id = ?)", repoID); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, "UPDATE `repository` SET num_stars = 0 WHERE id = ?", repoID); err != nil {
		return err
	}

	return db.DeleteBeans(ctx, Star{RepoID: repoID})
}
