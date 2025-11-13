// Copyright 2016 The Gogs Authors. All rights reserved.
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

package repo

import (
	"context"
	"time"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"
)

// ErrMirrorNotExist mirror does not exist error
var ErrMirrorNotExist = util.NewNotExistErrorf("Mirror does not exist")

// Mirror represents mirror information of a repository.
type Mirror struct {
	ID          int64       `xorm:"pk autoincr"`
	RepoID      int64       `xorm:"INDEX"`
	Repo        *Repository `xorm:"-"`
	Interval    time.Duration
	EnablePrune bool `xorm:"NOT NULL DEFAULT true"`

	UpdatedUnix    timeutil.TimeStamp `xorm:"INDEX"`
	NextUpdateUnix timeutil.TimeStamp `xorm:"INDEX"`

	LFS         bool   `xorm:"lfs_enabled NOT NULL DEFAULT false"`
	LFSEndpoint string `xorm:"lfs_endpoint TEXT"`

	RemoteAddress string `xorm:"VARCHAR(2048)"`
}

func init() {
	db.RegisterModel(new(Mirror))
}

// BeforeInsert will be invoked by XORM before inserting a record
func (m *Mirror) BeforeInsert() {
	if m != nil {
		m.UpdatedUnix = timeutil.TimeStampNow()
		m.NextUpdateUnix = timeutil.TimeStampNow()
	}
}

// GetRepository returns the repository.
func (m *Mirror) GetRepository(ctx context.Context) *Repository {
	if m.Repo != nil {
		return m.Repo
	}
	var err error
	m.Repo, err = GetRepositoryByID(ctx, m.RepoID)
	if err != nil {
		log.Error("getRepositoryByID[%d]: %v", m.ID, err)
	}
	return m.Repo
}

// GetRemoteName returns the name of the remote.
func (m *Mirror) GetRemoteName() string {
	return "origin"
}

// ScheduleNextUpdate calculates and sets next update time.
func (m *Mirror) ScheduleNextUpdate() {
	if m.Interval != 0 {
		m.NextUpdateUnix = timeutil.TimeStampNow().AddDuration(m.Interval)
	} else {
		m.NextUpdateUnix = 0
	}
}

// GetMirrorByRepoID returns mirror information of a repository.
func GetMirrorByRepoID(ctx context.Context, repoID int64) (*Mirror, error) {
	m := &Mirror{RepoID: repoID}
	has, err := db.GetEngine(ctx).Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrMirrorNotExist
	}
	return m, nil
}

// UpdateMirror updates the mirror
func UpdateMirror(ctx context.Context, m *Mirror) error {
	_, err := db.GetEngine(ctx).ID(m.ID).AllCols().Update(m)
	return err
}

// TouchMirror updates the mirror updatedUnix
func TouchMirror(ctx context.Context, m *Mirror) error {
	m.UpdatedUnix = timeutil.TimeStampNow()
	_, err := db.GetEngine(ctx).ID(m.ID).Cols("updated_unix").Update(m)
	return err
}

// DeleteMirrorByRepoID deletes a mirror by repoID
func DeleteMirrorByRepoID(ctx context.Context, repoID int64) error {
	_, err := db.GetEngine(ctx).Delete(&Mirror{RepoID: repoID})
	return err
}

// MirrorsIterate iterates all mirror repositories.
func MirrorsIterate(ctx context.Context, limit int, f func(idx int, bean any) error) error {
	sess := db.GetEngine(ctx).
		Where("next_update_unix<=?", time.Now().Unix()).
		And("next_update_unix!=0").
		OrderBy("updated_unix ASC")
	if limit > 0 {
		sess = sess.Limit(limit)
	}
	return sess.Iterate(new(Mirror), f)
}

// InsertMirror inserts a mirror to database
func InsertMirror(ctx context.Context, mirror *Mirror) error {
	_, err := db.GetEngine(ctx).Insert(mirror)
	return err
}
