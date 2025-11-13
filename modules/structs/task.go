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

// TaskType defines task type
type TaskType int

const TaskTypeMigrateRepo TaskType = iota // migrate repository from external or local disk

// Name returns the task type name
func (taskType TaskType) Name() string {
	switch taskType {
	case TaskTypeMigrateRepo:
		return "Migrate Repository"
	}
	return ""
}

// TaskStatus defines task status
type TaskStatus int

// enumerate all the kinds of task status
const (
	TaskStatusQueued   TaskStatus = iota // 0 task is queued
	TaskStatusRunning                    // 1 task is running
	TaskStatusStopped                    // 2 task is stopped (never used)
	TaskStatusFailed                     // 3 task is failed
	TaskStatusFinished                   // 4 task is finished
)
