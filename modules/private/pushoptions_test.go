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

package private

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitPushOptions(t *testing.T) {
	o := GitPushOptions{}

	v := o.Bool("no-such")
	assert.False(t, v.Has())
	assert.False(t, v.Value())

	o.AddFromKeyValue("opt1=a=b")
	o.AddFromKeyValue("opt2=false")
	o.AddFromKeyValue("opt3=true")
	o.AddFromKeyValue("opt4")

	assert.Equal(t, "a=b", o["opt1"])
	assert.False(t, o.Bool("opt1").Value())
	assert.True(t, o.Bool("opt2").Has())
	assert.False(t, o.Bool("opt2").Value())
	assert.True(t, o.Bool("opt3").Value())
	assert.True(t, o.Bool("opt4").Value())
}
