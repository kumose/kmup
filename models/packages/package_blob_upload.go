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
	"strings"
	"time"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"
)

// ErrPackageBlobUploadNotExist indicates a package blob upload not exist error
var ErrPackageBlobUploadNotExist = util.NewNotExistErrorf("package blob upload does not exist")

func init() {
	db.RegisterModel(new(PackageBlobUpload))
}

// PackageBlobUpload represents a package blob upload
type PackageBlobUpload struct {
	ID             string             `xorm:"pk"`
	BytesReceived  int64              `xorm:"NOT NULL DEFAULT 0"`
	HashStateBytes []byte             `xorm:"BLOB"`
	CreatedUnix    timeutil.TimeStamp `xorm:"created NOT NULL"`
	UpdatedUnix    timeutil.TimeStamp `xorm:"updated INDEX NOT NULL"`
}

// CreateBlobUpload inserts a blob upload
func CreateBlobUpload(ctx context.Context) (*PackageBlobUpload, error) {
	id, err := util.CryptoRandomString(25)
	if err != nil {
		return nil, err
	}

	pbu := &PackageBlobUpload{
		ID: strings.ToLower(id),
	}

	_, err = db.GetEngine(ctx).Insert(pbu)
	return pbu, err
}

// GetBlobUploadByID gets a blob upload by id
func GetBlobUploadByID(ctx context.Context, id string) (*PackageBlobUpload, error) {
	pbu := &PackageBlobUpload{}

	has, err := db.GetEngine(ctx).ID(id).Get(pbu)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrPackageBlobUploadNotExist
	}
	return pbu, nil
}

// UpdateBlobUpload updates the blob upload
func UpdateBlobUpload(ctx context.Context, pbu *PackageBlobUpload) error {
	_, err := db.GetEngine(ctx).ID(pbu.ID).Update(pbu)
	return err
}

// DeleteBlobUploadByID deletes the blob upload
func DeleteBlobUploadByID(ctx context.Context, id string) error {
	_, err := db.GetEngine(ctx).ID(id).Delete(&PackageBlobUpload{})
	return err
}

// FindExpiredBlobUploads gets all expired blob uploads
func FindExpiredBlobUploads(ctx context.Context, olderThan time.Duration) ([]*PackageBlobUpload, error) {
	pbus := make([]*PackageBlobUpload, 0, 10)
	return pbus, db.GetEngine(ctx).
		Where("updated_unix < ?", time.Now().Add(-olderThan).Unix()).
		Find(&pbus)
}
