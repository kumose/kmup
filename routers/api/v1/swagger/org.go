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

// Organization
// swagger:response Organization
type swaggerResponseOrganization struct {
	// in:body
	Body api.Organization `json:"body"`
}

// OrganizationList
// swagger:response OrganizationList
type swaggerResponseOrganizationList struct {
	// in:body
	Body []api.Organization `json:"body"`
}

// Team
// swagger:response Team
type swaggerResponseTeam struct {
	// in:body
	Body api.Team `json:"body"`
}

// TeamList
// swagger:response TeamList
type swaggerResponseTeamList struct {
	// in:body
	Body []api.Team `json:"body"`
}

// OrganizationPermissions
// swagger:response OrganizationPermissions
type swaggerResponseOrganizationPermissions struct {
	// in:body
	Body api.OrganizationPermissions `json:"body"`
}
