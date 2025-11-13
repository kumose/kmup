// Copyright 2015 The Gogs Authors. All rights reserved.
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
	"github.com/kumose/kmup/modules/git/gitcmd"
)

// BranchPrefix base dir of the branch information file store on git
const BranchPrefix = "refs/heads/"

// AddRemote adds a new remote to repository.
func (repo *Repository) AddRemote(name, url string, fetch bool) error {
	cmd := gitcmd.NewCommand("remote", "add")
	if fetch {
		cmd.AddArguments("-f")
	}
	_, _, err := cmd.AddDynamicArguments(name, url).
		WithDir(repo.Path).
		RunStdString(repo.Ctx)
	return err
}
