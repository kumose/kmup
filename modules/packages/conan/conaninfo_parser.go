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

package conan

import (
	"bufio"
	"io"
	"strings"

	"github.com/kumose/kmup/modules/util"
)

// Conaninfo represents infos of a Conan package
type Conaninfo struct {
	Settings     map[string]string   `json:"settings"`
	FullSettings map[string]string   `json:"full_settings"`
	Requires     []string            `json:"requires"`
	FullRequires []string            `json:"full_requires"`
	Options      map[string]string   `json:"options"`
	FullOptions  map[string]string   `json:"full_options"`
	RecipeHash   string              `json:"recipe_hash"`
	Environment  map[string][]string `json:"environment"`
}

func ParseConaninfo(r io.Reader) (*Conaninfo, error) {
	sections, err := readSections(io.LimitReader(r, 1<<20))
	if err != nil {
		return nil, err
	}

	info := &Conaninfo{}
	for section, lines := range sections {
		if len(lines) == 0 {
			continue
		}
		switch section {
		case "settings":
			info.Settings = toMap(lines)
		case "full_settings":
			info.FullSettings = toMap(lines)
		case "options":
			info.Options = toMap(lines)
		case "full_options":
			info.FullOptions = toMap(lines)
		case "requires":
			info.Requires = lines
		case "full_requires":
			info.FullRequires = lines
		case "recipe_hash":
			info.RecipeHash = lines[0]
		case "env":
			info.Environment = toMapArray(lines)
		}
	}
	return info, nil
}

func readSections(r io.Reader) (map[string][]string, error) {
	sections := make(map[string][]string)

	section := ""
	lines := make([]string, 0, 5)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if section != "" {
				sections[section] = lines
			}
			section = line[1 : len(line)-1]
			lines = make([]string, 0, 5)
			continue
		}
		if section != "" {
			if line != "" {
				lines = append(lines, line)
			}
			continue
		}
		if line != "" {
			return nil, util.NewInvalidArgumentErrorf("invalid conaninfo.txt")
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if section != "" {
		sections[section] = lines
	}
	return sections, nil
}

func toMap(lines []string) map[string]string {
	result := make(map[string]string)
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
			continue
		}
		result[parts[0]] = parts[1]
	}
	return result
}

func toMapArray(lines []string) map[string][]string {
	result := make(map[string][]string)
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
			continue
		}
		var items []string
		if strings.HasPrefix(parts[1], "[") && strings.HasSuffix(parts[1], "]") {
			items = strings.Split(parts[1], ",")
		} else {
			items = []string{parts[1]}
		}
		result[parts[0]] = items
	}
	return result
}
