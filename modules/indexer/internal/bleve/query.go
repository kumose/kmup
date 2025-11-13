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

package bleve

import (
	"github.com/kumose/kmup/modules/optional"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

// NumericEqualityQuery generates a numeric equality query for the given value and field
func NumericEqualityQuery(value int64, field string) *query.NumericRangeQuery {
	f := float64(value)
	tru := true
	q := bleve.NewNumericRangeInclusiveQuery(&f, &f, &tru, &tru)
	q.SetField(field)
	return q
}

// MatchPhraseQuery generates a match phrase query for the given phrase, field and analyzer
func MatchPhraseQuery(matchPhrase, field, analyzer string, fuzziness int) *query.MatchPhraseQuery {
	q := bleve.NewMatchPhraseQuery(matchPhrase)
	q.FieldVal = field
	q.Analyzer = analyzer
	q.Fuzziness = fuzziness
	return q
}

// MatchAndQuery generates a match query for the given phrase, field and analyzer
func MatchAndQuery(matchPhrase, field, analyzer string, fuzziness int) *query.MatchQuery {
	q := bleve.NewMatchQuery(matchPhrase)
	q.FieldVal = field
	q.Analyzer = analyzer
	q.Fuzziness = fuzziness
	q.Operator = query.MatchQueryOperatorAnd
	return q
}

// BoolFieldQuery generates a bool field query for the given value and field
func BoolFieldQuery(value bool, field string) *query.BoolFieldQuery {
	q := bleve.NewBoolFieldQuery(value)
	q.SetField(field)
	return q
}

func NumericRangeInclusiveQuery(minOption, maxOption optional.Option[int64], field string) *query.NumericRangeQuery {
	var minF, maxF *float64
	var minI, maxI *bool
	if minOption.Has() {
		minF = new(float64)
		*minF = float64(minOption.Value())
		minI = new(bool)
		*minI = true
	}
	if maxOption.Has() {
		maxF = new(float64)
		*maxF = float64(maxOption.Value())
		maxI = new(bool)
		*maxI = true
	}
	q := bleve.NewNumericRangeInclusiveQuery(minF, maxF, minI, maxI)
	q.SetField(field)
	return q
}
