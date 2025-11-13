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

package setting

import "github.com/kumose/kmup/modules/glob"

type GlobMatcher struct {
	compiledGlob  glob.Glob
	patternString string
}

var _ glob.Glob = (*GlobMatcher)(nil)

func (g *GlobMatcher) Match(s string) bool {
	return g.compiledGlob.Match(s)
}

func (g *GlobMatcher) PatternString() string {
	return g.patternString
}

func GlobMatcherCompile(pattern string, separators ...rune) (*GlobMatcher, error) {
	g, err := glob.Compile(pattern, separators...)
	if err != nil {
		return nil, err
	}
	return &GlobMatcher{
		compiledGlob:  g,
		patternString: pattern,
	}, nil
}
