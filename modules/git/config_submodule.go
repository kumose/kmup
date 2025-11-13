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
	"bufio"
	"fmt"
	"io"
	"strings"
)

// SubModule is a reference on git repository
type SubModule struct {
	Path   string
	URL    string
	Branch string // this field is newly added but not really used
}

// configParseSubModules this is not a complete parse for gitmodules file, it only
// parses the url and path of submodules. At the moment it only parses well-formed gitmodules files.
// In the future, there should be a complete implementation of https://git-scm.com/docs/git-config#_syntax
func configParseSubModules(r io.Reader) (*ObjectCache[*SubModule], error) {
	var subModule *SubModule
	subModules := newObjectCache[*SubModule]()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Section header [section]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if subModule != nil {
				subModules.Set(subModule.Path, subModule)
			}
			if strings.HasPrefix(line, "[submodule") {
				subModule = &SubModule{}
			} else {
				subModule = nil
			}
			continue
		}

		if subModule == nil {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "path":
			subModule.Path = value
		case "url":
			subModule.URL = value
		case "branch":
			subModule.Branch = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	if subModule != nil {
		subModules.Set(subModule.Path, subModule)
	}
	return subModules, nil
}
