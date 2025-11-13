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

package nuget

import (
	"context"
	"strings"

	"github.com/kumose/kmup/models/db"
	packages_model "github.com/kumose/kmup/models/packages"

	"xorm.io/builder"
)

// SearchVersions gets all versions of packages matching the search options
func SearchVersions(ctx context.Context, opts *packages_model.PackageSearchOptions) ([]*packages_model.PackageVersion, int64, error) {
	cond := toConds(opts)

	e := db.GetEngine(ctx)

	total, err := e.
		Where(cond).
		Count(&packages_model.Package{})
	if err != nil {
		return nil, 0, err
	}

	inner := builder.
		Dialect(db.BuilderDialect()). // builder needs the sql dialect to build the Limit() below
		Select("*").
		From("package").
		Where(cond).
		OrderBy("package.name ASC")
	if opts.Paginator != nil {
		skip, take := opts.Paginator.GetSkipTake()
		inner = inner.Limit(take, skip)
	}

	sess := e.
		Where(opts.ToConds()).
		Table("package_version").
		Join("INNER", inner, "package.id = package_version.package_id")

	pvs := make([]*packages_model.PackageVersion, 0, 10)
	return pvs, total, sess.Find(&pvs)
}

// CountPackages counts all packages matching the search options
func CountPackages(ctx context.Context, opts *packages_model.PackageSearchOptions) (int64, error) {
	return db.GetEngine(ctx).
		Where(toConds(opts)).
		Count(&packages_model.Package{})
}

func toConds(opts *packages_model.PackageSearchOptions) builder.Cond {
	var cond builder.Cond = builder.Eq{
		"package.is_internal": opts.IsInternal.Value(),
		"package.owner_id":    opts.OwnerID,
		"package.type":        packages_model.TypeNuGet,
	}
	if opts.Name.Value != "" {
		if opts.Name.ExactMatch {
			cond = cond.And(builder.Eq{"package.lower_name": strings.ToLower(opts.Name.Value)})
		} else {
			cond = cond.And(builder.Like{"package.lower_name", strings.ToLower(opts.Name.Value)})
		}
	}
	return cond
}
