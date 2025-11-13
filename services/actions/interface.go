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

package actions

import "github.com/kumose/kmup/services/context"

// API for actions of a repository or organization
type API interface {
	// ListActionsSecrets list secrets
	ListActionsSecrets(*context.APIContext)
	// CreateOrUpdateSecret create or update a secret
	CreateOrUpdateSecret(*context.APIContext)
	// DeleteSecret delete a secret
	DeleteSecret(*context.APIContext)
	// ListVariables list variables
	ListVariables(*context.APIContext)
	// GetVariable get a variable
	GetVariable(*context.APIContext)
	// DeleteVariable delete a variable
	DeleteVariable(*context.APIContext)
	// CreateVariable create a variable
	CreateVariable(*context.APIContext)
	// UpdateVariable update a variable
	UpdateVariable(*context.APIContext)
	// GetRegistrationToken get registration token
	GetRegistrationToken(*context.APIContext)
	// CreateRegistrationToken get registration token
	CreateRegistrationToken(*context.APIContext)
	// ListRunners list runners
	ListRunners(*context.APIContext)
	// GetRunner get a runner
	GetRunner(*context.APIContext)
	// DeleteRunner delete runner
	DeleteRunner(*context.APIContext)
	// ListWorkflowJobs list jobs
	ListWorkflowJobs(*context.APIContext)
	// ListWorkflowRuns list runs
	ListWorkflowRuns(*context.APIContext)
}
