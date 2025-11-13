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
	"io"
	"regexp"
	"strings"
)

var (
	patternAuthor      = compilePattern("author")
	patternHomepage    = compilePattern("homepage")
	patternURL         = compilePattern("url")
	patternLicense     = compilePattern("license")
	patternDescription = compilePattern("description")
	patternTopics      = regexp.MustCompile(`(?im)^\s*topics\s*=\s*\((.+)\)`)
	patternTopicList   = regexp.MustCompile(`\s*['"](.+?)['"]\s*,?`)
)

func compilePattern(name string) *regexp.Regexp {
	return regexp.MustCompile(`(?im)^\s*` + name + `\s*=\s*['"\(](.+)['"\)]`)
}

func ParseConanfile(r io.Reader) (*Metadata, error) {
	buf, err := io.ReadAll(io.LimitReader(r, 1<<20))
	if err != nil {
		return nil, err
	}

	metadata := &Metadata{}

	m := patternAuthor.FindSubmatch(buf)
	if len(m) > 1 && len(m[1]) > 0 {
		metadata.Author = string(m[1])
	}
	m = patternHomepage.FindSubmatch(buf)
	if len(m) > 1 && len(m[1]) > 0 {
		metadata.ProjectURL = string(m[1])
	}
	m = patternURL.FindSubmatch(buf)
	if len(m) > 1 && len(m[1]) > 0 {
		metadata.RepositoryURL = string(m[1])
	}
	m = patternLicense.FindSubmatch(buf)
	if len(m) > 1 && len(m[1]) > 0 {
		metadata.License = strings.ReplaceAll(strings.ReplaceAll(string(m[1]), "'", ""), "\"", "")
	}
	m = patternDescription.FindSubmatch(buf)
	if len(m) > 1 && len(m[1]) > 0 {
		metadata.Description = string(m[1])
	}
	m = patternTopics.FindSubmatch(buf)
	if len(m) > 1 && len(m[1]) > 0 {
		m2 := patternTopicList.FindAllSubmatch(m[1], -1)
		if len(m2) > 0 {
			metadata.Keywords = make([]string, 0, len(m2))
			for _, g := range m2 {
				if len(g) > 1 {
					metadata.Keywords = append(metadata.Keywords, string(g[1]))
				}
			}
		}
	}
	return metadata, nil
}
