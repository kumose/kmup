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

package arch

import (
	"context"

	packages_model "github.com/kumose/kmup/models/packages"
	arch_module "github.com/kumose/kmup/modules/packages/arch"
)

// GetRepositories gets all available repositories
func GetRepositories(ctx context.Context, ownerID int64) ([]string, error) {
	return packages_model.GetDistinctPropertyValues(
		ctx,
		packages_model.TypeArch,
		ownerID,
		packages_model.PropertyTypeFile,
		arch_module.PropertyRepository,
		nil,
	)
}

// GetArchitectures gets all available architectures for the given repository
func GetArchitectures(ctx context.Context, ownerID int64, repository string) ([]string, error) {
	return packages_model.GetDistinctPropertyValues(
		ctx,
		packages_model.TypeArch,
		ownerID,
		packages_model.PropertyTypeFile,
		arch_module.PropertyArchitecture,
		&packages_model.DistinctPropertyDependency{
			Name:  arch_module.PropertyRepository,
			Value: repository,
		},
	)
}
