// Copyright 2014 The Gogs Authors. All rights reserved.
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

package forms

import (
	"net/http"

	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/web/middleware"
	"github.com/kumose/kmup/services/context"

	"github.com/kumose-go/chi/binding"
)

// ________                            .__                __  .__
// \_____  \_______  _________    ____ |__|____________ _/  |_|__| ____   ____
//  /   |   \_  __ \/ ___\__  \  /    \|  \___   /\__  \\   __\  |/  _ \ /    \
// /    |    \  | \/ /_/  > __ \|   |  \  |/    /  / __ \|  | |  (  <_> )   |  \
// \_______  /__|  \___  (____  /___|  /__/_____ \(____  /__| |__|\____/|___|  /
//         \/     /_____/     \/     \/         \/     \/                    \/

// CreateOrgForm form for creating organization
type CreateOrgForm struct {
	OrgName                   string `binding:"Required;Username;MaxSize(40)" locale:"org.org_name_holder"`
	Visibility                structs.VisibleType
	RepoAdminChangeTeamAccess bool
}

// Validate validates the fields
func (f *CreateOrgForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

// UpdateOrgSettingForm form for updating organization settings
type UpdateOrgSettingForm struct {
	FullName                  string `binding:"MaxSize(100)"`
	Email                     string `binding:"MaxSize(255)"`
	Description               string `binding:"MaxSize(255)"`
	Website                   string `binding:"ValidUrl;MaxSize(255)"`
	Location                  string `binding:"MaxSize(50)"`
	MaxRepoCreation           int
	RepoAdminChangeTeamAccess bool
}

// Validate validates the fields
func (f *UpdateOrgSettingForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

type RenameOrgForm struct {
	OrgName    string `binding:"Required"`
	NewOrgName string `binding:"Required;Username;MaxSize(40)" locale:"org.org_name_holder"`
}

// ___________
// \__    ___/___ _____    _____
//   |    |_/ __ \\__  \  /     \
//   |    |\  ___/ / __ \|  Y Y  \
//   |____| \___  >____  /__|_|  /
//              \/     \/      \/

// CreateTeamForm form for creating team
type CreateTeamForm struct {
	TeamName         string `binding:"Required;AlphaDashDot;MaxSize(255)"`
	Description      string `binding:"MaxSize(255)"`
	Permission       string
	RepoAccess       string
	CanCreateOrgRepo bool
}

// Validate validates the fields
func (f *CreateTeamForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}
