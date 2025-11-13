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

package convert

import (
	"context"

	"github.com/kumose/kmup/models/packages"
	access_model "github.com/kumose/kmup/models/perm/access"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
)

// ToPackage convert a packages.PackageDescriptor to api.Package
func ToPackage(ctx context.Context, pd *packages.PackageDescriptor, doer *user_model.User) (*api.Package, error) {
	var repo *api.Repository
	if pd.Repository != nil {
		permission, err := access_model.GetUserRepoPermission(ctx, pd.Repository, doer)
		if err != nil {
			return nil, err
		}

		if permission.HasAnyUnitAccess() {
			repo = ToRepo(ctx, pd.Repository, permission)
		}
	}

	return &api.Package{
		ID:         pd.Version.ID,
		Owner:      ToUser(ctx, pd.Owner, doer),
		Repository: repo,
		Creator:    ToUser(ctx, pd.Creator, doer),
		Type:       string(pd.Package.Type),
		Name:       pd.Package.Name,
		Version:    pd.Version.Version,
		CreatedAt:  pd.Version.CreatedUnix.AsTime(),
		HTMLURL:    pd.VersionHTMLURL(ctx),
	}, nil
}

// ToPackageFile converts packages.PackageFileDescriptor to api.PackageFile
func ToPackageFile(pfd *packages.PackageFileDescriptor) *api.PackageFile {
	return &api.PackageFile{
		ID:         pfd.File.ID,
		Size:       pfd.Blob.Size,
		Name:       pfd.File.Name,
		HashMD5:    pfd.Blob.HashMD5,
		HashSHA1:   pfd.Blob.HashSHA1,
		HashSHA256: pfd.Blob.HashSHA256,
		HashSHA512: pfd.Blob.HashSHA512,
	}
}
