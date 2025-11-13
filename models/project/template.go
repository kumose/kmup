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

package project

type (
	// TemplateType is used to represent a project template type
	TemplateType uint8

	// TemplateConfig is used to identify the template type of project that is being created
	TemplateConfig struct {
		TemplateType TemplateType
		Translation  string
	}
)

const (
	// TemplateTypeNone is a project template type that has no predefined columns
	TemplateTypeNone TemplateType = iota

	// TemplateTypeBasicKanban is a project template type that has basic predefined columns
	TemplateTypeBasicKanban

	// TemplateTypeBugTriage is a project template type that has predefined columns suited to hunting down bugs
	TemplateTypeBugTriage
)

// GetTemplateConfigs retrieves the template configs of configurations project columns could have
func GetTemplateConfigs() []TemplateConfig {
	return []TemplateConfig{
		{TemplateTypeNone, "repo.projects.type.none"},
		{TemplateTypeBasicKanban, "repo.projects.type.basic_kanban"},
		{TemplateTypeBugTriage, "repo.projects.type.bug_triage"},
	}
}

// IsTemplateTypeValid checks if the project template type is valid
func IsTemplateTypeValid(p TemplateType) bool {
	switch p {
	case TemplateTypeNone, TemplateTypeBasicKanban, TemplateTypeBugTriage:
		return true
	default:
		return false
	}
}
