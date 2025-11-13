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

package packages

import (
	"context"
	"strconv"
	"time"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/perm"
	"github.com/kumose/kmup/models/unit"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"

	"xorm.io/builder"
)

// ErrPackageBlobNotExist indicates a package blob not exist error
var ErrPackageBlobNotExist = util.NewNotExistErrorf("package blob does not exist")

func init() {
	db.RegisterModel(new(PackageBlob))
}

// PackageBlob represents a package blob
type PackageBlob struct {
	ID          int64              `xorm:"pk autoincr"`
	Size        int64              `xorm:"NOT NULL DEFAULT 0"`
	HashMD5     string             `xorm:"hash_md5 char(32) UNIQUE(md5) INDEX NOT NULL"`
	HashSHA1    string             `xorm:"hash_sha1 char(40) UNIQUE(sha1) INDEX NOT NULL"`
	HashSHA256  string             `xorm:"hash_sha256 char(64) UNIQUE(sha256) INDEX NOT NULL"`
	HashSHA512  string             `xorm:"hash_sha512 char(128) UNIQUE(sha512) INDEX NOT NULL"`
	CreatedUnix timeutil.TimeStamp `xorm:"created INDEX NOT NULL"`
}

// GetOrInsertBlob inserts a blob. If the blob exists already the existing blob is returned
func GetOrInsertBlob(ctx context.Context, pb *PackageBlob) (*PackageBlob, bool, error) {
	e := db.GetEngine(ctx)

	existing := &PackageBlob{}

	has, err := e.Where(builder.Eq{
		"size":        pb.Size,
		"hash_md5":    pb.HashMD5,
		"hash_sha1":   pb.HashSHA1,
		"hash_sha256": pb.HashSHA256,
		"hash_sha512": pb.HashSHA512,
	}).Get(existing)
	if err != nil {
		return nil, false, err
	}
	if has {
		return existing, true, nil
	}
	if _, err = e.Insert(pb); err != nil {
		return nil, false, err
	}
	return pb, false, nil
}

// GetBlobByID gets a blob by id
func GetBlobByID(ctx context.Context, blobID int64) (*PackageBlob, error) {
	pb := &PackageBlob{}

	has, err := db.GetEngine(ctx).ID(blobID).Get(pb)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrPackageBlobNotExist
	}
	return pb, nil
}

// ExistPackageBlobWithSHA returns if a package blob exists with the provided sha
func ExistPackageBlobWithSHA(ctx context.Context, blobSha256 string) (bool, error) {
	return db.GetEngine(ctx).Exist(&PackageBlob{
		HashSHA256: blobSha256,
	})
}

// FindExpiredUnreferencedBlobs gets all blobs without associated files older than the specific duration
func FindExpiredUnreferencedBlobs(ctx context.Context, olderThan time.Duration) ([]*PackageBlob, error) {
	pbs := make([]*PackageBlob, 0, 10)
	return pbs, db.GetEngine(ctx).
		Table("package_blob").
		Join("LEFT", "package_file", "package_file.blob_id = package_blob.id").
		Where("package_file.id IS NULL AND package_blob.created_unix < ?", time.Now().Add(-olderThan).Unix()).
		Find(&pbs)
}

// DeleteBlobByID deletes a blob by id
func DeleteBlobByID(ctx context.Context, blobID int64) error {
	_, err := db.GetEngine(ctx).ID(blobID).Delete(&PackageBlob{})
	return err
}

// GetTotalBlobSize returns the total blobs size in bytes
func GetTotalBlobSize(ctx context.Context) (int64, error) {
	return db.GetEngine(ctx).
		SumInt(&PackageBlob{}, "size")
}

// GetTotalUnreferencedBlobSize returns the total size of all unreferenced blobs in bytes
func GetTotalUnreferencedBlobSize(ctx context.Context) (int64, error) {
	return db.GetEngine(ctx).
		Table("package_blob").
		Join("LEFT", "package_file", "package_file.blob_id = package_blob.id").
		Where("package_file.id IS NULL").
		SumInt(&PackageBlob{}, "size")
}

// IsBlobAccessibleForUser tests if the user has access to the blob
func IsBlobAccessibleForUser(ctx context.Context, blobID int64, user *user_model.User) (bool, error) {
	if user.IsAdmin {
		return true, nil
	}

	maxTeamAuthorize := builder.
		Select("max(team.authorize)").
		From("team").
		InnerJoin("team_user", "team_user.team_id = team.id").
		Where(builder.Eq{"team_user.uid": user.ID}.And(builder.Expr("team_user.org_id = `user`.id")))

	maxTeamUnitAccessMode := builder.
		Select("max(team_unit.access_mode)").
		From("team").
		InnerJoin("team_user", "team_user.team_id = team.id").
		InnerJoin("team_unit", "team_unit.team_id = team.id").
		Where(builder.Eq{"team_user.uid": user.ID, "team_unit.type": unit.TypePackages}.And(builder.Expr("team_user.org_id = `user`.id")))

	cond := builder.Eq{"package_blob.id": blobID}.And(
		// owner = user
		builder.Eq{"`user`.id": user.ID}.
			// user can see owner
			Or(builder.Eq{"`user`.visibility": structs.VisibleTypePublic}.Or(builder.Eq{"`user`.visibility": structs.VisibleTypeLimited})).
			// owner is an organization and user has access to it
			Or(builder.Eq{"`user`.type": user_model.UserTypeOrganization}.
				And(builder.Lte{strconv.Itoa(int(perm.AccessModeRead)): maxTeamAuthorize}.Or(builder.Lte{strconv.Itoa(int(perm.AccessModeRead)): maxTeamUnitAccessMode}))),
	)

	return db.GetEngine(ctx).
		Table("package_blob").
		Join("INNER", "package_file", "package_file.blob_id = package_blob.id").
		Join("INNER", "package_version", "package_version.id = package_file.version_id").
		Join("INNER", "package", "package.id = package_version.package_id").
		Join("INNER", "user", "`user`.id = package.owner_id").
		Where(cond).
		Exist(&PackageBlob{})
}
