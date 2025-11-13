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

package structs

import "time"

// Secret represents a secret
// swagger:model
type Secret struct {
	// the secret's name
	Name string `json:"name"`
	// the secret's description
	Description string `json:"description"`
	// swagger:strfmt date-time
	Created time.Time `json:"created_at"`
}

// CreateOrUpdateSecretOption options when creating or updating secret
// swagger:model
type CreateOrUpdateSecretOption struct {
	// Data of the secret to update
	//
	// required: true
	Data string `json:"data" binding:"Required"`

	// Description of the secret to update
	//
	// required: false
	Description string `json:"description"`
}
