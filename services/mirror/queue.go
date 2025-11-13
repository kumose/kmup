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

package mirror

import (
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/queue"
	"github.com/kumose/kmup/modules/setting"
)

var mirrorQueue *queue.WorkerPoolQueue[*SyncRequest]

// SyncType type of sync request
type SyncType int

const (
	// PullMirrorType for pull mirrors
	PullMirrorType SyncType = iota
	// PushMirrorType for push mirrors
	PushMirrorType
)

// SyncRequest for the mirror queue
type SyncRequest struct {
	Type        SyncType
	ReferenceID int64 // RepoID for pull mirror, MirrorID for push mirror
}

func queueHandler(items ...*SyncRequest) []*SyncRequest {
	for _, req := range items {
		doMirrorSync(graceful.GetManager().ShutdownContext(), req)
	}
	return nil
}

// StartSyncMirrors starts a go routine to sync the mirrors
func StartSyncMirrors() {
	if !setting.Mirror.Enabled {
		return
	}
	mirrorQueue = queue.CreateUniqueQueue(graceful.GetManager().ShutdownContext(), "mirror", queueHandler)
	if mirrorQueue == nil {
		log.Fatal("Unable to create mirror queue")
	}
	go graceful.GetManager().RunWithCancel(mirrorQueue)
}

// AddPullMirrorToQueue adds repoID to mirror queue
func AddPullMirrorToQueue(repoID int64) {
	addMirrorToQueue(PullMirrorType, repoID)
}

// AddPushMirrorToQueue adds the push mirror to the queue
func AddPushMirrorToQueue(mirrorID int64) {
	addMirrorToQueue(PushMirrorType, mirrorID)
}

func addMirrorToQueue(syncType SyncType, referenceID int64) {
	if !setting.Mirror.Enabled {
		return
	}
	go func() {
		if err := PushToQueue(syncType, referenceID); err != nil {
			log.Error("Unable to push sync request for to the queue for pull mirror repo[%d]. Error: %v", referenceID, err)
		}
	}()
}

// PushToQueue adds the sync request to the queue
func PushToQueue(mirrorType SyncType, referenceID int64) error {
	return mirrorQueue.Push(&SyncRequest{
		Type:        mirrorType,
		ReferenceID: referenceID,
	})
}
