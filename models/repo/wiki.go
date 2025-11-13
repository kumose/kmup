// Copyright 2015 The Gogs Authors. All rights reserved.
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

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/util"
)

// ErrWikiAlreadyExist represents a "WikiAlreadyExist" kind of error.
type ErrWikiAlreadyExist struct {
	Title string
}

// IsErrWikiAlreadyExist checks if an error is an ErrWikiAlreadyExist.
func IsErrWikiAlreadyExist(err error) bool {
	_, ok := err.(ErrWikiAlreadyExist)
	return ok
}

func (err ErrWikiAlreadyExist) Error() string {
	return fmt.Sprintf("wiki page already exists [title: %s]", err.Title)
}

func (err ErrWikiAlreadyExist) Unwrap() error {
	return util.ErrAlreadyExist
}

// ErrWikiReservedName represents a reserved name error.
type ErrWikiReservedName struct {
	Title string
}

// IsErrWikiReservedName checks if an error is an ErrWikiReservedName.
func IsErrWikiReservedName(err error) bool {
	_, ok := err.(ErrWikiReservedName)
	return ok
}

func (err ErrWikiReservedName) Error() string {
	return "wiki title is reserved: " + err.Title
}

func (err ErrWikiReservedName) Unwrap() error {
	return util.ErrInvalidArgument
}

// ErrWikiInvalidFileName represents an invalid wiki file name.
type ErrWikiInvalidFileName struct {
	FileName string
}

// IsErrWikiInvalidFileName checks if an error is an ErrWikiInvalidFileName.
func IsErrWikiInvalidFileName(err error) bool {
	_, ok := err.(ErrWikiInvalidFileName)
	return ok
}

func (err ErrWikiInvalidFileName) Error() string {
	return "Invalid wiki filename: " + err.FileName
}

func (err ErrWikiInvalidFileName) Unwrap() error {
	return util.ErrInvalidArgument
}

// WikiCloneLink returns clone URLs of repository wiki.
func (repo *Repository) WikiCloneLink(ctx context.Context, doer *user_model.User) *CloneLink {
	return repo.cloneLink(ctx, doer, repo.Name+".wiki")
}

func RelativeWikiPath(ownerName, repoName string) string {
	return strings.ToLower(ownerName) + "/" + strings.ToLower(repoName) + ".wiki.git"
}

// WikiStorageRepo returns the storage repo for the wiki
// The wiki repository should have the same object format as the code repository
func (repo *Repository) WikiStorageRepo() StorageRepo {
	return StorageRepo(RelativeWikiPath(repo.OwnerName, repo.Name))
}
