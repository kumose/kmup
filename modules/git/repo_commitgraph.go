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

package git

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/modules/git/gitcmd"
)

// WriteCommitGraph write commit graph to speed up repo access
// this requires git v2.18 to be installed
func WriteCommitGraph(ctx context.Context, repoPath string) error {
	if DefaultFeatures().CheckVersionAtLeast("2.18") {
		if _, _, err := gitcmd.NewCommand("commit-graph", "write").WithDir(repoPath).RunStdString(ctx); err != nil {
			return fmt.Errorf("unable to write commit-graph for '%s' : %w", repoPath, err)
		}
	}
	return nil
}
