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
	"fmt"
	"testing"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/stretchr/testify/assert"
)

type Scenario struct {
	Input  string
	Tokens []string
}

func TestTokenFilter(t *testing.T) {
	scenarios := []struct {
		Input string
		Terms []string
	}{
		{
			Input: "Dockerfile",
			Terms: []string{"Dockerfile"},
		},
		{
			Input: "Dockerfile.rootless",
			Terms: []string{"Dockerfile.rootless"},
		},
		{
			Input: "a/b/c/Dockerfile.rootless",
			Terms: []string{"a", "a/b", "a/b/c", "a/b/c/Dockerfile.rootless", "Dockerfile.rootless", "Dockerfile.rootless/c", "Dockerfile.rootless/c/b", "Dockerfile.rootless/c/b/a"},
		},
		{
			Input: "",
			Terms: []string{},
		},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("ensure terms of '%s'", scenario.Input), func(t *testing.T) {
			terms := extractTerms(scenario.Input)

			assert.Len(t, terms, len(scenario.Terms))

			for _, term := range terms {
				assert.Contains(t, scenario.Terms, term)
			}
		})
	}
}

func extractTerms(input string) []string {
	tokens := tokenize(input)
	filteredTokens := filter(tokens)
	terms := make([]string, 0, len(filteredTokens))

	for _, token := range filteredTokens {
		terms = append(terms, string(token.Term))
	}

	return terms
}

func filter(input analysis.TokenStream) analysis.TokenStream {
	filter := NewTokenFilter()
	return filter.Filter(input)
}

func tokenize(input string) analysis.TokenStream {
	tokenizer := unicode.NewUnicodeTokenizer()
	return tokenizer.Tokenize([]byte(input))
}
