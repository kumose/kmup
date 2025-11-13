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

// Label a label to an issue or a pr
// swagger:model
type Label struct {
	// ID is the unique identifier for the label
	ID int64 `json:"id"`
	// Name is the display name of the label
	Name string `json:"name"`
	// example: false
	Exclusive bool `json:"exclusive"`
	// example: false
	IsArchived bool `json:"is_archived"`
	// example: 00aabb
	Color string `json:"color"`
	// Description provides additional context about the label's purpose
	Description string `json:"description"`
	// URL is the API endpoint for accessing this label
	URL string `json:"url"`
}

// CreateLabelOption options for creating a label
type CreateLabelOption struct {
	// required:true
	// Name is the display name for the new label
	Name string `json:"name" binding:"Required"`
	// example: false
	Exclusive bool `json:"exclusive"`
	// required:true
	// example: #00aabb
	Color string `json:"color" binding:"Required"`
	// Description provides additional context about the label's purpose
	Description string `json:"description"`
	// example: false
	IsArchived bool `json:"is_archived"`
}

// EditLabelOption options for editing a label
type EditLabelOption struct {
	// Name is the new display name for the label
	Name *string `json:"name"`
	// example: false
	Exclusive *bool `json:"exclusive"`
	// example: #00aabb
	Color *string `json:"color"`
	// Description provides additional context about the label's purpose
	Description *string `json:"description"`
	// example: false
	IsArchived *bool `json:"is_archived"`
}

// IssueLabelsOption a collection of labels
type IssueLabelsOption struct {
	// Labels can be a list of integers representing label IDs
	// or a list of strings representing label names
	Labels []any `json:"labels"`
}

// LabelTemplate info of a Label template
type LabelTemplate struct {
	// Name is the display name of the label template
	Name string `json:"name"`
	// example: false
	Exclusive bool `json:"exclusive"`
	// example: 00aabb
	Color string `json:"color"`
	// Description provides additional context about the label template's purpose
	Description string `json:"description"`
}
