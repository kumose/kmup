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

type SubmoduleWebLink struct {
	RepoWebLink, CommitWebLink string
}

// GetSubModules get all the submodules of current revision git tree
func (c *Commit) GetSubModules() (*ObjectCache[*SubModule], error) {
	if c.submoduleCache != nil {
		return c.submoduleCache, nil
	}

	entry, err := c.GetTreeEntryByPath(".gitmodules")
	if err != nil {
		if _, ok := err.(ErrNotExist); ok {
			return nil, nil
		}
		return nil, err
	}

	rd, err := entry.Blob().DataAsync()
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	// at the moment we do not strictly limit the size of the .gitmodules file because some users would have huge .gitmodules files (>1MB)
	c.submoduleCache, err = configParseSubModules(rd)
	if err != nil {
		return nil, err
	}
	return c.submoduleCache, nil
}

// GetSubModule gets the submodule by the entry name.
// It returns "nil, nil" if the submodule does not exist, caller should always remember to check the "nil"
func (c *Commit) GetSubModule(entryName string) (*SubModule, error) {
	modules, err := c.GetSubModules()
	if err != nil {
		return nil, err
	}

	if modules != nil {
		if module, has := modules.Get(entryName); has {
			return module, nil
		}
	}
	return nil, nil
}
