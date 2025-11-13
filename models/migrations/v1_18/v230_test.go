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

package v1_18

import (
	"testing"

	"github.com/kumose/kmup/models/migrations/base"

	"github.com/stretchr/testify/assert"
)

func Test_AddConfidentialClientColumnToOAuth2ApplicationTable(t *testing.T) {
	// premigration
	type oauth2Application struct {
		ID int64
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(oauth2Application))
	defer deferable()
	if x == nil || t.Failed() {
		return
	}

	if err := AddConfidentialClientColumnToOAuth2ApplicationTable(x); err != nil {
		assert.NoError(t, err)
		return
	}

	// postmigration
	type ExpectedOAuth2Application struct {
		ID                 int64
		ConfidentialClient bool
	}

	got := []ExpectedOAuth2Application{}
	if err := x.Table("oauth2_application").Select("id, confidential_client").Find(&got); !assert.NoError(t, err) {
		return
	}

	assert.NotEmpty(t, got)
	for _, e := range got {
		assert.True(t, e.ConfidentialClient)
	}
}
