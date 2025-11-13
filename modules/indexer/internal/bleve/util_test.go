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
	"fmt"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestBleveGuessFuzzinessByKeyword(t *testing.T) {
	defer test.MockVariableValue(&setting.Indexer.TypeBleveMaxFuzzniess, 2)()

	scenarios := []struct {
		Input     string
		Fuzziness int // See util.go for the definition of fuzziness in this particular context
	}{
		{
			Input:     "",
			Fuzziness: 0,
		},
		{
			Input:     "Avocado",
			Fuzziness: 1,
		},
		{
			Input:     "Geschwindigkeit",
			Fuzziness: 2,
		},
		{
			Input:     "non-exist",
			Fuzziness: 0,
		},
		{
			Input:     "갃갃갃",
			Fuzziness: 0,
		},
		{
			Input:     "repo1",
			Fuzziness: 0,
		},
		{
			Input:     "avocado.md",
			Fuzziness: 0,
		},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("Fuziniess:%s=%d", scenario.Input, scenario.Fuzziness), func(t *testing.T) {
			assert.Equal(t, scenario.Fuzziness, GuessFuzzinessByKeyword(scenario.Input))
		})
	}
}
