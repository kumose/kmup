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

package swagger

import (
	api "github.com/kumose/kmup/modules/structs"
)

// ServerVersion
// swagger:response ServerVersion
type swaggerResponseServerVersion struct {
	// in:body
	Body api.ServerVersion `json:"body"`
}

// GitignoreTemplateList
// swagger:response GitignoreTemplateList
type swaggerResponseGitignoreTemplateList struct {
	// in:body
	Body []string `json:"body"`
}

// GitignoreTemplateInfo
// swagger:response GitignoreTemplateInfo
type swaggerResponseGitignoreTemplateInfo struct {
	// in:body
	Body api.GitignoreTemplateInfo `json:"body"`
}

// LicenseTemplateList
// swagger:response LicenseTemplateList
type swaggerResponseLicensesTemplateList struct {
	// in:body
	Body []api.LicensesTemplateListEntry `json:"body"`
}

// LicenseTemplateInfo
// swagger:response LicenseTemplateInfo
type swaggerResponseLicenseTemplateInfo struct {
	// in:body
	Body api.LicenseTemplateInfo `json:"body"`
}

// StringSlice
// swagger:response StringSlice
type swaggerResponseStringSlice struct {
	// in:body
	Body []string `json:"body"`
}

// LabelTemplateList
// swagger:response LabelTemplateList
type swaggerResponseLabelTemplateList struct {
	// in:body
	Body []string `json:"body"`
}

// LabelTemplateInfo
// swagger:response LabelTemplateInfo
type swaggerResponseLabelTemplateInfo struct {
	// in:body
	Body []api.LabelTemplate `json:"body"`
}
