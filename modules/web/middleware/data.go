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

package middleware

import (
	"context"
	"time"

	"github.com/kumose/kmup/modules/reqctx"
	"github.com/kumose/kmup/modules/setting"
)

const ContextDataKeySignedUser = "SignedUser"

func GetContextData(c context.Context) reqctx.ContextData {
	if rc := reqctx.GetRequestDataStore(c); rc != nil {
		return rc.GetData()
	}
	return nil
}

func CommonTemplateContextData() reqctx.ContextData {
	return reqctx.ContextData{
		"PageTitleCommon": setting.AppName,

		"IsLandingPageOrganizations": setting.LandingPageURL == setting.LandingPageOrganizations,

		"ShowRegistrationButton":        setting.Service.ShowRegistrationButton,
		"ShowMilestonesDashboardPage":   setting.Service.ShowMilestonesDashboardPage,
		"ShowFooterVersion":             setting.Other.ShowFooterVersion,
		"DisableDownloadSourceArchives": setting.Repository.DisableDownloadSourceArchives,

		"EnableSwagger":      setting.API.EnableSwagger,
		"EnableOpenIDSignIn": setting.Service.EnableOpenIDSignIn,
		"PageStartTime":      time.Now(),

		"RunModeIsProd": setting.IsProd,
	}
}
