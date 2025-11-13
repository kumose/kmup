// Copyright 2014 The Gogs Authors. All rights reserved.
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

package cron

import (
	"context"
	"runtime/pprof"
	"time"

	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/translation"

	"github.com/go-co-op/gocron"
)

var scheduler = gocron.NewScheduler(time.Local)

// Init begins cron tasks
// Each cron task is run within the shutdown context as a running server
// AtShutdown the cron server is stopped
func Init(original context.Context) {
	defer pprof.SetGoroutineLabels(original)
	_, _, finished := process.GetManager().AddTypedContext(graceful.GetManager().ShutdownContext(), "Service: Cron", process.SystemProcessType, true)
	initBasicTasks()
	initExtendedTasks()
	initActionsTasks()

	lock.Lock()
	for _, task := range tasks {
		if task.IsEnabled() && task.DoRunAtStart() {
			go task.Run()
		}
	}

	scheduler.StartAsync()
	started = true
	lock.Unlock()
	graceful.GetManager().RunAtShutdown(context.Background(), func() {
		scheduler.Stop()
		lock.Lock()
		started = false
		lock.Unlock()
		finished()
	})
}

// TaskTableRow represents a task row in the tasks table
type TaskTableRow struct {
	Name        string
	Spec        string
	Next        time.Time
	Prev        time.Time
	Status      string
	LastMessage string
	LastDoer    string
	ExecTimes   int64
	task        *Task
}

func (t *TaskTableRow) FormatLastMessage(locale translation.Locale) string {
	if t.Status == "finished" {
		return t.task.GetConfig().FormatMessage(locale, t.Name, t.Status, t.LastDoer)
	}

	return t.task.GetConfig().FormatMessage(locale, t.Name, t.Status, t.LastDoer, t.LastMessage)
}

// TaskTable represents a table of tasks
type TaskTable []*TaskTableRow

// ListTasks returns all running cron tasks.
func ListTasks() TaskTable {
	jobs := scheduler.Jobs()
	jobMap := map[string]*gocron.Job{}
	for _, job := range jobs {
		// the first tag is the task name
		tags := job.Tags()
		if len(tags) == 0 { // should never happen
			continue
		}
		jobMap[job.Tags()[0]] = job
	}

	lock.Lock()
	defer lock.Unlock()

	tTable := make([]*TaskTableRow, 0, len(tasks))
	for _, task := range tasks {
		spec := "-"
		var (
			next time.Time
			prev time.Time
		)
		if e, ok := jobMap[task.Name]; ok {
			tags := e.Tags()
			if len(tags) > 1 {
				spec = tags[1] // the second tag is the task spec
			}
			next = e.NextRun()
			prev = e.PreviousRun()
		}

		task.lock.Lock()
		// If the manual run is after the cron run, use that instead.
		if prev.Before(task.LastRun) {
			prev = task.LastRun
		}
		tTable = append(tTable, &TaskTableRow{
			Name:        task.Name,
			Spec:        spec,
			Next:        next,
			Prev:        prev,
			ExecTimes:   task.ExecTimes,
			LastMessage: task.LastMessage,
			Status:      task.Status,
			LastDoer:    task.LastDoer,
			task:        task,
		})
		task.lock.Unlock()
	}

	return tTable
}
