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

package v1_15

import (
	"strings"
	"testing"

	"github.com/kumose/kmup/models/migrations/base"

	"github.com/stretchr/testify/assert"
)

func Test_AddPrimaryEmail2EmailAddress(t *testing.T) {
	type User struct {
		ID       int64
		Email    string
		IsActive bool
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(User))
	if x == nil || t.Failed() {
		defer deferable()
		return
	}
	defer deferable()

	err := AddPrimaryEmail2EmailAddress(x)
	assert.NoError(t, err)

	type EmailAddress struct {
		ID          int64  `xorm:"pk autoincr"`
		UID         int64  `xorm:"INDEX NOT NULL"`
		Email       string `xorm:"UNIQUE NOT NULL"`
		LowerEmail  string `xorm:"UNIQUE NOT NULL"`
		IsActivated bool
		IsPrimary   bool `xorm:"DEFAULT(false) NOT NULL"`
	}

	users := make([]User, 0, 20)
	err = x.Find(&users)
	assert.NoError(t, err)

	for _, user := range users {
		var emailAddress EmailAddress
		has, err := x.Where("lower_email=?", strings.ToLower(user.Email)).Get(&emailAddress)
		assert.NoError(t, err)
		assert.True(t, has)
		assert.True(t, emailAddress.IsPrimary)
		assert.Equal(t, user.IsActive, emailAddress.IsActivated)
		assert.Equal(t, user.ID, emailAddress.UID)
	}
}
