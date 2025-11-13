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

// CreateVariableOption the option when creating variable
// swagger:model
type CreateVariableOption struct {
	// Value of the variable to create
	//
	// required: true
	Value string `json:"value" binding:"Required"`

	// Description of the variable to create
	//
	// required: false
	Description string `json:"description"`
}

// UpdateVariableOption the option when updating variable
// swagger:model
type UpdateVariableOption struct {
	// New name for the variable. If the field is empty, the variable name won't be updated.
	Name string `json:"name"`
	// Value of the variable to update
	//
	// required: true
	Value string `json:"value" binding:"Required"`

	// Description of the variable to update
	//
	// required: false
	Description string `json:"description"`
}

// ActionVariable return value of the query API
// swagger:model
type ActionVariable struct {
	// the owner to which the variable belongs
	OwnerID int64 `json:"owner_id"`
	// the repository to which the variable belongs
	RepoID int64 `json:"repo_id"`
	// the name of the variable
	Name string `json:"name"`
	// the value of the variable
	Data string `json:"data"`
	// the description of the variable
	Description string `json:"description"`
}
