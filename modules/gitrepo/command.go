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

package gitrepo

import (
	"context"

	"github.com/kumose/kmup/modules/git/gitcmd"
)

func RunCmd(ctx context.Context, repo Repository, cmd *gitcmd.Command) error {
	return cmd.WithDir(repoPath(repo)).WithParentCallerInfo().Run(ctx)
}

func RunCmdString(ctx context.Context, repo Repository, cmd *gitcmd.Command) (string, error) {
	res, _, err := cmd.WithDir(repoPath(repo)).WithParentCallerInfo().RunStdString(ctx)
	return res, err
}

func RunCmdBytes(ctx context.Context, repo Repository, cmd *gitcmd.Command) ([]byte, []byte, error) {
	return cmd.WithDir(repoPath(repo)).WithParentCallerInfo().RunStdBytes(ctx)
}
