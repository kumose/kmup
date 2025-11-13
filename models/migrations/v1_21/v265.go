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

package v1_21

import (
	"xorm.io/xorm"
)

func AlterActionArtifactTable(x *xorm.Engine) error {
	// ActionArtifact is a file that is stored in the artifact storage.
	type ActionArtifact struct {
		RunID        int64  `xorm:"index unique(runid_name_path)"` // The run id of the artifact
		ArtifactPath string `xorm:"index unique(runid_name_path)"` // The path to the artifact when runner uploads it
		ArtifactName string `xorm:"index unique(runid_name_path)"` // The name of the artifact when
	}

	return x.Sync(new(ActionArtifact))
}
