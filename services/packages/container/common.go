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

package container

import (
	"context"
	"io"
	"strings"

	packages_model "github.com/kumose/kmup/models/packages"
	container_service "github.com/kumose/kmup/models/packages/container"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/packages"
	container_module "github.com/kumose/kmup/modules/packages/container"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// UpdateRepositoryNames updates the repository name property for all packages of the specific owner
func UpdateRepositoryNames(ctx context.Context, owner *user_model.User, newOwnerName string) error {
	ps, err := packages_model.GetPackagesByType(ctx, owner.ID, packages_model.TypeContainer)
	if err != nil {
		return err
	}

	newOwnerName = strings.ToLower(newOwnerName)

	for _, p := range ps {
		if err := packages_model.DeletePropertiesByName(ctx, packages_model.PropertyTypePackage, p.ID, container_module.PropertyRepository); err != nil {
			return err
		}

		if _, err := packages_model.InsertProperty(ctx, packages_model.PropertyTypePackage, p.ID, container_module.PropertyRepository, newOwnerName+"/"+p.LowerName); err != nil {
			return err
		}
	}

	return nil
}

func ParseManifestMetadata(ctx context.Context, rd io.Reader, ownerID int64, imageName string) (*v1.Manifest, *packages_model.PackageFileDescriptor, *container_module.Metadata, error) {
	var manifest v1.Manifest
	if err := json.NewDecoder(rd).Decode(&manifest); err != nil {
		return nil, nil, nil, err
	}
	configDescriptor, err := container_service.GetContainerBlob(ctx, &container_service.BlobSearchOptions{
		OwnerID: ownerID,
		Image:   imageName,
		Digest:  manifest.Config.Digest.String(),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	configReader, err := packages.NewContentStore().OpenBlob(packages.BlobHash256Key(configDescriptor.Blob.HashSHA256))
	if err != nil {
		return nil, nil, nil, err
	}
	defer configReader.Close()
	metadata, err := container_module.ParseImageConfig(manifest.Config.MediaType, configReader)
	return &manifest, configDescriptor, metadata, err
}
