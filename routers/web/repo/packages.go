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
	"net/http"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/packages"
	"github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/services/context"
)

const (
	tplPackagesList templates.TplName = "repo/packages"
)

// Packages displays a list of all packages in the repository
func Packages(ctx *context.Context) {
	page := max(ctx.FormInt("page"), 1)
	query := ctx.FormTrim("q")
	packageType := ctx.FormTrim("type")

	pvs, total, err := packages.SearchLatestVersions(ctx, &packages.PackageSearchOptions{
		Paginator: &db.ListOptions{
			PageSize: setting.UI.PackagesPagingNum,
			Page:     page,
		},
		OwnerID:    ctx.ContextUser.ID,
		RepoID:     ctx.Repo.Repository.ID,
		Type:       packages.Type(packageType),
		Name:       packages.SearchValue{Value: query},
		IsInternal: optional.Some(false),
	})
	if err != nil {
		ctx.ServerError("SearchLatestVersions", err)
		return
	}

	pds, err := packages.GetPackageDescriptors(ctx, pvs)
	if err != nil {
		ctx.ServerError("GetPackageDescriptors", err)
		return
	}

	hasPackages, err := packages.HasRepositoryPackages(ctx, ctx.Repo.Repository.ID)
	if err != nil {
		ctx.ServerError("HasRepositoryPackages", err)
		return
	}

	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["IsPackagesPage"] = true
	ctx.Data["Query"] = query
	ctx.Data["PackageType"] = packageType
	ctx.Data["AvailableTypes"] = packages.TypeList
	ctx.Data["HasPackages"] = hasPackages
	ctx.Data["CanWritePackages"] = ctx.Repo.CanWrite(unit.TypePackages) || ctx.IsUserSiteAdmin()
	ctx.Data["PackageDescriptors"] = pds
	ctx.Data["Total"] = total
	ctx.Data["RepositoryAccessMap"] = map[int64]bool{ctx.Repo.Repository.ID: true} // There is only the current repository

	pager := context.NewPagination(int(total), setting.UI.PackagesPagingNum, page, 5)
	pager.AddParamFromRequest(ctx.Req)
	ctx.Data["Page"] = pager

	ctx.HTML(http.StatusOK, tplPackagesList)
}
