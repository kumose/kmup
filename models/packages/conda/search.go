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

package conda

import (
	"context"
	"strings"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/packages"
	conda_module "github.com/kumose/kmup/modules/packages/conda"

	"xorm.io/builder"
)

type FileSearchOptions struct {
	OwnerID  int64
	Channel  string
	Subdir   string
	Filename string
}

// SearchFiles gets all files matching the search options
func SearchFiles(ctx context.Context, opts *FileSearchOptions) ([]*packages.PackageFile, error) {
	var cond builder.Cond = builder.Eq{
		"package.type":                packages.TypeConda,
		"package.owner_id":            opts.OwnerID,
		"package_version.is_internal": false,
	}

	if opts.Filename != "" {
		cond = cond.And(builder.Eq{
			"package_file.lower_name": strings.ToLower(opts.Filename),
		})
	}

	var versionPropsCond builder.Cond = builder.Eq{
		"package_property.ref_type": packages.PropertyTypePackage,
		"package_property.name":     conda_module.PropertyChannel,
		"package_property.value":    opts.Channel,
	}

	cond = cond.And(builder.In("package.id", builder.Select("package_property.ref_id").Where(versionPropsCond).From("package_property")))

	var filePropsCond builder.Cond = builder.Eq{
		"package_property.ref_type": packages.PropertyTypeFile,
		"package_property.name":     conda_module.PropertySubdir,
		"package_property.value":    opts.Subdir,
	}

	cond = cond.And(builder.In("package_file.id", builder.Select("package_property.ref_id").Where(filePropsCond).From("package_property")))

	sess := db.GetEngine(ctx).
		Select("package_file.*").
		Table("package_file").
		Join("INNER", "package_version", "package_version.id = package_file.version_id").
		Join("INNER", "package", "package.id = package_version.package_id").
		Where(cond)

	pfs := make([]*packages.PackageFile, 0, 10)
	return pfs, sess.Find(&pfs)
}
