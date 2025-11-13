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

import api "github.com/kumose/kmup/modules/structs"

// GeneralRepoSettings
// swagger:response GeneralRepoSettings
type swaggerResponseGeneralRepoSettings struct {
	// in:body
	Body api.GeneralRepoSettings `json:"body"`
}

// GeneralUISettings
// swagger:response GeneralUISettings
type swaggerResponseGeneralUISettings struct {
	// in:body
	Body api.GeneralUISettings `json:"body"`
}

// GeneralAPISettings
// swagger:response GeneralAPISettings
type swaggerResponseGeneralAPISettings struct {
	// in:body
	Body api.GeneralAPISettings `json:"body"`
}

// GeneralAttachmentSettings
// swagger:response GeneralAttachmentSettings
type swaggerResponseGeneralAttachmentSettings struct {
	// in:body
	Body api.GeneralAttachmentSettings `json:"body"`
}
