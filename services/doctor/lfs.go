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

package doctor

import (
	"context"
	"errors"
	"time"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/repository"
)

func init() {
	Register(&Check{
		Title:                      "Garbage collect LFS",
		Name:                       "gc-lfs",
		IsDefault:                  false,
		Run:                        garbageCollectLFSCheck,
		AbortIfFailed:              false,
		SkipDatabaseInitialization: false,
		Priority:                   1,
	})
}

func garbageCollectLFSCheck(ctx context.Context, logger log.Logger, autofix bool) error {
	if !setting.LFS.StartServer {
		return errors.New("LFS support is disabled")
	}

	if err := repository.GarbageCollectLFSMetaObjects(ctx, repository.GarbageCollectLFSMetaObjectsOptions{
		LogDetail: logger.Info,
		AutoFix:   autofix,
		// Only attempt to garbage collect lfs meta objects older than a week as the order of git lfs upload
		// and git object upload is not necessarily guaranteed. It's possible to imagine a situation whereby
		// an LFS object is uploaded but the git branch is not uploaded immediately, or there are some rapid
		// changes in new branches that might lead to lfs objects becoming temporarily unassociated with git
		// objects.
		//
		// It is likely that a week is potentially excessive but it should definitely be enough that any
		// unassociated LFS object is genuinely unassociated.
		OlderThan: time.Now().Add(-24 * time.Hour * 7),
		// We don't set the UpdatedLessRecentlyThan because we want to do a full GC
	}); err != nil {
		return err
	}

	return checkStorage(&checkStorageOptions{LFS: true})(ctx, logger, autofix)
}
