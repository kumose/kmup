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

package charset

import (
	"sort"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestAmbiguousCharacters(t *testing.T) {
	for locale, ambiguous := range AmbiguousCharacters {
		assert.Equal(t, locale, ambiguous.Locale)
		assert.Len(t, ambiguous.With, len(ambiguous.Confusable))
		assert.True(t, sort.SliceIsSorted(ambiguous.Confusable, func(i, j int) bool {
			return ambiguous.Confusable[i] < ambiguous.Confusable[j]
		}))

		for _, confusable := range ambiguous.Confusable {
			assert.True(t, unicode.Is(ambiguous.RangeTable, confusable))
			i := sort.Search(len(ambiguous.Confusable), func(j int) bool {
				return ambiguous.Confusable[j] >= confusable
			})
			found := i < len(ambiguous.Confusable) && ambiguous.Confusable[i] == confusable
			assert.True(t, found, "%c is not in %d", confusable, i)
		}
	}
}
