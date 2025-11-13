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

package auth

import (
	"context"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/log"
)

// SyncExternalUsers is used to synchronize users with external authorization source
func SyncExternalUsers(ctx context.Context, updateExisting bool) error {
	log.Trace("Doing: SyncExternalUsers")

	ls, err := db.Find[auth.Source](ctx, auth.FindSourcesOptions{})
	if err != nil {
		log.Error("SyncExternalUsers: %v", err)
		return err
	}

	for _, s := range ls {
		if !s.IsActive || !s.IsSyncEnabled {
			continue
		}
		select {
		case <-ctx.Done():
			log.Warn("SyncExternalUsers: Cancelled before update of %s", s.Name)
			return db.ErrCancelledf("Before update of %s", s.Name)
		default:
		}

		if syncable, ok := s.Cfg.(SynchronizableSource); ok {
			err := syncable.Sync(ctx, updateExisting)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
