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

package types

import "github.com/kumose/kmup/modules/translation"

type OwnerType string

const (
	OwnerTypeSystemGlobal = "system-global"
	OwnerTypeIndividual   = "individual"
	OwnerTypeRepository   = "repository"
	OwnerTypeOrganization = "organization"
)

func (o OwnerType) LocaleString(locale translation.Locale) string {
	switch o {
	case OwnerTypeSystemGlobal:
		return locale.TrString("concept_system_global")
	case OwnerTypeIndividual:
		return locale.TrString("concept_user_individual")
	case OwnerTypeRepository:
		return locale.TrString("concept_code_repository")
	case OwnerTypeOrganization:
		return locale.TrString("concept_user_organization")
	}
	return locale.TrString("unknown")
}
