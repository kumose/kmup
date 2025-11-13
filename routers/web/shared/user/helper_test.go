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

package user

import (
	"testing"

	"github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestMakeSelfOnTop(t *testing.T) {
	users := MakeSelfOnTop(nil, []*user.User{{ID: 2}, {ID: 1}})
	assert.Len(t, users, 2)
	assert.EqualValues(t, 2, users[0].ID)

	users = MakeSelfOnTop(&user.User{ID: 1}, []*user.User{{ID: 2}, {ID: 1}})
	assert.Len(t, users, 2)
	assert.EqualValues(t, 1, users[0].ID)

	users = MakeSelfOnTop(&user.User{ID: 2}, []*user.User{{ID: 2}, {ID: 1}})
	assert.Len(t, users, 2)
	assert.EqualValues(t, 2, users[0].ID)
}
