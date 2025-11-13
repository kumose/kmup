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

package repository

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/perm"
	access_model "github.com/kumose/kmup/models/perm/access"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"

	"xorm.io/builder"
)

func AddOrUpdateCollaborator(ctx context.Context, repo *repo_model.Repository, u *user_model.User, mode perm.AccessMode) error {
	// only allow valid access modes, read, write and admin
	if mode < perm.AccessModeRead || mode > perm.AccessModeAdmin {
		return perm.ErrInvalidAccessMode
	}

	if err := repo.LoadOwner(ctx); err != nil {
		return err
	}

	if user_model.IsUserBlockedBy(ctx, u, repo.OwnerID) || user_model.IsUserBlockedBy(ctx, repo.Owner, u.ID) {
		return user_model.ErrBlockedUser
	}

	return db.WithTx(ctx, func(ctx context.Context) error {
		collaboration, has, err := db.Get[repo_model.Collaboration](ctx, builder.Eq{
			"repo_id": repo.ID,
			"user_id": u.ID,
		})
		if err != nil {
			return err
		} else if has {
			if collaboration.Mode == mode {
				return nil
			}
			if _, err = db.GetEngine(ctx).
				Where("repo_id=?", repo.ID).
				And("user_id=?", u.ID).
				Cols("mode").
				Update(&repo_model.Collaboration{
					Mode: mode,
				}); err != nil {
				return err
			}
		} else if err = db.Insert(ctx, &repo_model.Collaboration{
			RepoID: repo.ID,
			UserID: u.ID,
			Mode:   mode,
		}); err != nil {
			return err
		}

		return access_model.RecalculateUserAccess(ctx, repo, u.ID)
	})
}

// DeleteCollaboration removes collaboration relation between the user and repository.
func DeleteCollaboration(ctx context.Context, repo *repo_model.Repository, collaborator *user_model.User) (err error) {
	collaboration := &repo_model.Collaboration{
		RepoID: repo.ID,
		UserID: collaborator.ID,
	}

	return db.WithTx(ctx, func(ctx context.Context) error {
		if has, err := db.GetEngine(ctx).Delete(collaboration); err != nil {
			return err
		} else if has == 0 {
			return nil
		}

		if err := repo.LoadOwner(ctx); err != nil {
			return err
		}

		if err = access_model.RecalculateAccesses(ctx, repo); err != nil {
			return err
		}

		if err = repo_model.WatchRepo(ctx, collaborator, repo, false); err != nil {
			return err
		}

		if err = ReconsiderWatches(ctx, repo, collaborator); err != nil {
			return err
		}

		// Unassign a user from any issue (s)he has been assigned to in the repository
		return ReconsiderRepoIssuesAssignee(ctx, repo, collaborator)
	})
}

func ReconsiderRepoIssuesAssignee(ctx context.Context, repo *repo_model.Repository, user *user_model.User) error {
	if canAssigned, err := access_model.CanBeAssigned(ctx, user, repo, true); err != nil || canAssigned {
		return err
	}

	if _, err := db.GetEngine(ctx).Where(builder.Eq{"assignee_id": user.ID}).
		In("issue_id", builder.Select("id").From("issue").Where(builder.Eq{"repo_id": repo.ID})).
		Delete(&issues_model.IssueAssignees{}); err != nil {
		return fmt.Errorf("Could not delete assignee[%d] %w", user.ID, err)
	}
	return nil
}

func ReconsiderWatches(ctx context.Context, repo *repo_model.Repository, user *user_model.User) error {
	if has, err := access_model.HasAnyUnitAccess(ctx, user.ID, repo); err != nil || has {
		return err
	}
	if err := repo_model.WatchRepo(ctx, user, repo, false); err != nil {
		return err
	}

	// Remove all IssueWatches a user has subscribed to in the repository
	return issues_model.RemoveIssueWatchersByRepoID(ctx, user.ID, repo.ID)
}
