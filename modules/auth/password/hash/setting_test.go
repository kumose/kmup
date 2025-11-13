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

package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSettingPasswordHashAlgorithm(t *testing.T) {
	t.Run("pbkdf2 is pbkdf2_v2", func(t *testing.T) {
		pbkdf2v2Config, pbkdf2v2Algo := SetDefaultPasswordHashAlgorithm("pbkdf2_v2")
		pbkdf2Config, pbkdf2Algo := SetDefaultPasswordHashAlgorithm("pbkdf2")

		assert.Equal(t, pbkdf2v2Config, pbkdf2Config)
		assert.Equal(t, pbkdf2v2Algo.Specification, pbkdf2Algo.Specification)
	})

	for a, b := range aliasAlgorithmNames {
		t.Run(a+"="+b, func(t *testing.T) {
			aConfig, aAlgo := SetDefaultPasswordHashAlgorithm(a)
			bConfig, bAlgo := SetDefaultPasswordHashAlgorithm(b)

			assert.Equal(t, bConfig, aConfig)
			assert.Equal(t, aAlgo.Specification, bAlgo.Specification)
		})
	}

	t.Run("pbkdf2_v2 is the default when default password hash algorithm is empty", func(t *testing.T) {
		emptyConfig, emptyAlgo := SetDefaultPasswordHashAlgorithm("")
		pbkdf2v2Config, pbkdf2v2Algo := SetDefaultPasswordHashAlgorithm("pbkdf2_v2")

		assert.Equal(t, pbkdf2v2Config, emptyConfig)
		assert.Equal(t, pbkdf2v2Algo.Specification, emptyAlgo.Specification)
	})
}
