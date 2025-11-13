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

package process

import (
	"context"
	"time"
)

var (
	SystemProcessType  = "system"
	RequestProcessType = "request"
	NormalProcessType  = "normal"
	NoneProcessType    = "none"
)

// process represents a working process inheriting from Kmup.
type process struct {
	PID         IDType // Process ID, not system one.
	ParentPID   IDType
	Description string
	Start       time.Time
	Cancel      context.CancelFunc
	Type        string
}

// ToProcess converts a process to a externally usable Process
func (p *process) toProcess() *Process {
	process := &Process{
		PID:         p.PID,
		ParentPID:   p.ParentPID,
		Description: p.Description,
		Start:       p.Start,
		Type:        p.Type,
	}
	return process
}
