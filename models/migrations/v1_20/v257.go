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

package v1_20

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func CreateActionArtifactTable(x *xorm.Engine) error {
	// ActionArtifact is a file that is stored in the artifact storage.
	type ActionArtifact struct {
		ID                 int64 `xorm:"pk autoincr"`
		RunID              int64 `xorm:"index UNIQUE(runid_name)"` // The run id of the artifact
		RunnerID           int64
		RepoID             int64 `xorm:"index"`
		OwnerID            int64
		CommitSHA          string
		StoragePath        string             // The path to the artifact in the storage
		FileSize           int64              // The size of the artifact in bytes
		FileCompressedSize int64              // The size of the artifact in bytes after gzip compression
		ContentEncoding    string             // The content encoding of the artifact
		ArtifactPath       string             // The path to the artifact when runner uploads it
		ArtifactName       string             `xorm:"UNIQUE(runid_name)"` // The name of the artifact when runner uploads it
		Status             int64              `xorm:"index"`              // The status of the artifact
		CreatedUnix        timeutil.TimeStamp `xorm:"created"`
		UpdatedUnix        timeutil.TimeStamp `xorm:"updated index"`
	}

	return x.Sync(new(ActionArtifact))
}
