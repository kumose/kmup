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

package path

import (
	"slices"
	"strings"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

const (
	Name = "kmup/path"
)

type TokenFilter struct{}

func NewTokenFilter() *TokenFilter {
	return &TokenFilter{}
}

func TokenFilterConstructor(config map[string]any, cache *registry.Cache) (analysis.TokenFilter, error) {
	return NewTokenFilter(), nil
}

func (s *TokenFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	if len(input) == 1 {
		// if there is only one token, we dont need to generate the reversed chain
		return generatePathTokens(input, false)
	}

	normal := generatePathTokens(input, false)
	reversed := generatePathTokens(input, true)

	return append(normal, reversed...)
}

// Generates path tokens from the input tokens.
// This mimics the behavior of the path hierarchy tokenizer in ES. It takes the input tokens and combine them, generating a term for each component
// in tree (e.g., foo/bar/baz.md will generate foo, foo/bar, and foo/bar/baz.md).
//
// If the reverse flag is set, the order of the tokens is reversed (the same input will generate baz.md, baz.md/bar, baz.md/bar/foo). This is useful
// to efficiently search for filenames without supplying the fullpath.
func generatePathTokens(input analysis.TokenStream, reversed bool) analysis.TokenStream {
	terms := make([]string, 0, len(input))
	longestTerm := 0

	if reversed {
		slices.Reverse(input)
	}

	for i := range input {
		var sb strings.Builder
		sb.Write(input[0].Term)

		for j := 1; j < i; j++ {
			sb.WriteString("/")
			sb.Write(input[j].Term)
		}

		term := sb.String()

		if longestTerm < len(term) {
			longestTerm = len(term)
		}

		terms = append(terms, term)
	}

	output := make(analysis.TokenStream, 0, len(terms))

	for _, term := range terms {
		var start, end int

		if reversed {
			start = 0
			end = len(term)
		} else {
			start = longestTerm - len(term)
			end = longestTerm
		}

		token := analysis.Token{
			Position: 1,
			Start:    start,
			End:      end,
			Type:     analysis.AlphaNumeric,
			Term:     []byte(term),
		}

		output = append(output, &token)
	}

	return output
}

func init() {
	// FIXME: move it to the bleve's init function, but do not call it in global init
	err := registry.RegisterTokenFilter(Name, TokenFilterConstructor)
	if err != nil {
		panic(err)
	}
}
