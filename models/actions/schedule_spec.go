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

import (
	"context"
	"strings"
	"time"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/timeutil"

	"github.com/robfig/cron/v3"
)

// ActionScheduleSpec represents a schedule spec of a workflow file
type ActionScheduleSpec struct {
	ID         int64
	RepoID     int64                  `xorm:"index"`
	Repo       *repo_model.Repository `xorm:"-"`
	ScheduleID int64                  `xorm:"index"`
	Schedule   *ActionSchedule        `xorm:"-"`

	// Next time the job will run, or the zero time if Cron has not been
	// started or this entry's schedule is unsatisfiable
	Next timeutil.TimeStamp `xorm:"index"`
	// Prev is the last time this job was run, or the zero time if never.
	Prev timeutil.TimeStamp
	Spec string

	Created timeutil.TimeStamp `xorm:"created"`
	Updated timeutil.TimeStamp `xorm:"updated"`
}

// Parse parses the spec and returns a cron.Schedule
// Unlike the default cron parser, Parse uses UTC timezone as the default if none is specified.
func (s *ActionScheduleSpec) Parse() (cron.Schedule, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := parser.Parse(s.Spec)
	if err != nil {
		return nil, err
	}

	// If the spec has specified a timezone, use it
	if strings.HasPrefix(s.Spec, "TZ=") || strings.HasPrefix(s.Spec, "CRON_TZ=") {
		return schedule, nil
	}

	specSchedule, ok := schedule.(*cron.SpecSchedule)
	// If it's not a spec schedule, like "@every 5m", timezone is not relevant
	if !ok {
		return schedule, nil
	}

	// Set the timezone to UTC
	specSchedule.Location = time.UTC
	return specSchedule, nil
}

func init() {
	db.RegisterModel(new(ActionScheduleSpec))
}

func UpdateScheduleSpec(ctx context.Context, spec *ActionScheduleSpec, cols ...string) error {
	sess := db.GetEngine(ctx).ID(spec.ID)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	_, err := sess.Update(spec)
	return err
}
