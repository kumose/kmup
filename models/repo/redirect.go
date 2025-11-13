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
	"strings"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/util"
)

// ErrRedirectNotExist represents a "RedirectNotExist" kind of error.
type ErrRedirectNotExist struct {
	OwnerID  int64
	RepoName string
}

// IsErrRedirectNotExist check if an error is an ErrRepoRedirectNotExist.
func IsErrRedirectNotExist(err error) bool {
	_, ok := err.(ErrRedirectNotExist)
	return ok
}

func (err ErrRedirectNotExist) Error() string {
	return fmt.Sprintf("repository redirect does not exist [uid: %d, name: %s]", err.OwnerID, err.RepoName)
}

func (err ErrRedirectNotExist) Unwrap() error {
	return util.ErrNotExist
}

// Redirect represents that a repo name should be redirected to another
type Redirect struct {
	ID             int64  `xorm:"pk autoincr"`
	OwnerID        int64  `xorm:"UNIQUE(s)"`
	LowerName      string `xorm:"UNIQUE(s) INDEX NOT NULL"`
	RedirectRepoID int64  // repoID to redirect to
}

// TableName represents real table name in database
func (Redirect) TableName() string {
	return "repo_redirect"
}

func init() {
	db.RegisterModel(new(Redirect))
}

// LookupRedirect look up if a repository has a redirect name
func LookupRedirect(ctx context.Context, ownerID int64, repoName string) (int64, error) {
	repoName = strings.ToLower(repoName)
	redirect := &Redirect{OwnerID: ownerID, LowerName: repoName}
	if has, err := db.GetEngine(ctx).Get(redirect); err != nil {
		return 0, err
	} else if !has {
		return 0, ErrRedirectNotExist{OwnerID: ownerID, RepoName: repoName}
	}
	return redirect.RedirectRepoID, nil
}

// NewRedirect create a new repo redirect
func NewRedirect(ctx context.Context, ownerID, repoID int64, oldRepoName, newRepoName string) error {
	oldRepoName = strings.ToLower(oldRepoName)
	newRepoName = strings.ToLower(newRepoName)

	if err := DeleteRedirect(ctx, ownerID, newRepoName); err != nil {
		return err
	}

	return db.Insert(ctx, &Redirect{
		OwnerID:        ownerID,
		LowerName:      oldRepoName,
		RedirectRepoID: repoID,
	})
}

// DeleteRedirect delete any redirect from the specified repo name to
// anything else
func DeleteRedirect(ctx context.Context, ownerID int64, repoName string) error {
	repoName = strings.ToLower(repoName)
	_, err := db.GetEngine(ctx).Delete(&Redirect{OwnerID: ownerID, LowerName: repoName})
	return err
}
