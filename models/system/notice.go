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

package system

import (
	"context"
	"fmt"
	"time"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/storage"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"
)

// NoticeType describes the notice type
type NoticeType int

const (
	// NoticeRepository type
	NoticeRepository NoticeType = iota + 1
	// NoticeTask type
	NoticeTask
)

// Notice represents a system notice for admin.
type Notice struct {
	ID          int64 `xorm:"pk autoincr"`
	Type        NoticeType
	Description string             `xorm:"LONGTEXT"`
	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
}

func init() {
	db.RegisterModel(new(Notice))
}

// TrStr returns a translation format string.
func (n *Notice) TrStr() string {
	return fmt.Sprintf("admin.notices.type_%d", n.Type)
}

// CreateNotice creates new system notice.
func CreateNotice(ctx context.Context, tp NoticeType, desc string, args ...any) error {
	if len(args) > 0 {
		desc = fmt.Sprintf(desc, args...)
	}
	n := &Notice{
		Type:        tp,
		Description: desc,
	}
	return db.Insert(ctx, n)
}

// CreateRepositoryNotice creates new system notice with type NoticeRepository.
func CreateRepositoryNotice(desc string, args ...any) error {
	return CreateNotice(graceful.GetManager().ShutdownContext(), NoticeRepository, desc, args...)
}

// RemoveAllWithNotice removes all directories in given path and
// creates a system notice when error occurs.
func RemoveAllWithNotice(ctx context.Context, title, path string) {
	if err := util.RemoveAll(path); err != nil {
		desc := fmt.Sprintf("%s [%s]: %v", title, path, err)
		log.Warn(title+" [%s]: %v", path, err)
		if err = CreateNotice(graceful.GetManager().ShutdownContext(), NoticeRepository, desc); err != nil {
			log.Error("CreateRepositoryNotice: %v", err)
		}
	}
}

// RemoveStorageWithNotice removes a file from the storage and
// creates a system notice when error occurs.
func RemoveStorageWithNotice(ctx context.Context, bucket storage.ObjectStorage, title, path string) {
	if err := bucket.Delete(path); err != nil {
		desc := fmt.Sprintf("%s [%s]: %v", title, path, err)
		log.Warn(title+" [%s]: %v", path, err)

		if err = CreateNotice(graceful.GetManager().ShutdownContext(), NoticeRepository, desc); err != nil {
			log.Error("CreateRepositoryNotice: %v", err)
		}
	}
}

// CountNotices returns number of notices.
func CountNotices(ctx context.Context) int64 {
	count, _ := db.GetEngine(ctx).Count(new(Notice))
	return count
}

// Notices returns notices in given page.
func Notices(ctx context.Context, page, pageSize int) ([]*Notice, error) {
	notices := make([]*Notice, 0, pageSize)
	return notices, db.GetEngine(ctx).
		Limit(pageSize, (page-1)*pageSize).
		Desc("created_unix").
		Find(&notices)
}

// DeleteNotices deletes all notices with ID from start to end (inclusive).
func DeleteNotices(ctx context.Context, start, end int64) error {
	if start == 0 && end == 0 {
		_, err := db.GetEngine(ctx).Exec("DELETE FROM notice")
		return err
	}

	sess := db.GetEngine(ctx).Where("id >= ?", start)
	if end > 0 {
		sess.And("id <= ?", end)
	}
	_, err := sess.Delete(new(Notice))
	return err
}

// DeleteOldSystemNotices deletes all old system notices from database.
func DeleteOldSystemNotices(ctx context.Context, olderThan time.Duration) (err error) {
	if olderThan <= 0 {
		return nil
	}

	_, err = db.GetEngine(ctx).Where("created_unix < ?", time.Now().Add(-olderThan).Unix()).Delete(&Notice{})
	return err
}
