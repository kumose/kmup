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

package attribute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Attribute(t *testing.T) {
	assert.Empty(t, Attribute("").ToString().Value())
	assert.Empty(t, Attribute("unspecified").ToString().Value())
	assert.Equal(t, "python", Attribute("python").ToString().Value())
	assert.Equal(t, "Java", Attribute("Java").ToString().Value())

	attributes := Attributes{
		m: map[string]Attribute{
			LinguistGenerated:     "true",
			LinguistDocumentation: "false",
			LinguistDetectable:    "set",
			LinguistLanguage:      "Python",
			GitlabLanguage:        "Java",
			"filter":              "unspecified",
			"test":                "",
		},
	}

	assert.Empty(t, attributes.Get("test").ToString().Value())
	assert.Empty(t, attributes.Get("filter").ToString().Value())
	assert.Equal(t, "Python", attributes.Get(LinguistLanguage).ToString().Value())
	assert.Equal(t, "Java", attributes.Get(GitlabLanguage).ToString().Value())
	assert.True(t, attributes.Get(LinguistGenerated).ToBool().Value())
	assert.False(t, attributes.Get(LinguistDocumentation).ToBool().Value())
	assert.True(t, attributes.Get(LinguistDetectable).ToBool().Value())
}
