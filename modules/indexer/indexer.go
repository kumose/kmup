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

package indexer

type SearchModeType string

const (
	SearchModeExact  SearchModeType = "exact"
	SearchModeWords  SearchModeType = "words"
	SearchModeFuzzy  SearchModeType = "fuzzy"
	SearchModeRegexp SearchModeType = "regexp"
)

type SearchMode struct {
	ModeValue    SearchModeType
	TooltipTrKey string
	TitleTrKey   string
}

func SearchModesExactWords() []SearchMode {
	return []SearchMode{
		{
			ModeValue:    SearchModeExact,
			TooltipTrKey: "search.exact_tooltip",
			TitleTrKey:   "search.exact",
		},
		{
			ModeValue:    SearchModeWords,
			TooltipTrKey: "search.words_tooltip",
			TitleTrKey:   "search.words",
		},
	}
}

func SearchModesExactWordsFuzzy() []SearchMode {
	return append(SearchModesExactWords(), []SearchMode{
		{
			ModeValue:    SearchModeFuzzy,
			TooltipTrKey: "search.fuzzy_tooltip",
			TitleTrKey:   "search.fuzzy",
		},
	}...)
}

func GitGrepSupportedSearchModes() []SearchMode {
	return append(SearchModesExactWords(), []SearchMode{
		{
			ModeValue:    SearchModeRegexp,
			TooltipTrKey: "search.regexp_tooltip",
			TitleTrKey:   "search.regexp",
		},
	}...)
}
