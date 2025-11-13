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

package doctor

import (
	"slices"
	"testing"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsistencyCheck(t *testing.T) {
	checks := prepareDBConsistencyChecks()
	idx := slices.IndexFunc(checks, func(check consistencyCheck) bool {
		return check.Name == "Orphaned OAuth2Application without existing User"
	})
	require.NotEqual(t, -1, idx)

	_ = db.TruncateBeans(t.Context(), &auth.OAuth2Application{}, &user.User{})
	_ = db.TruncateBeans(t.Context(), &auth.OAuth2Application{}, &auth.OAuth2Application{})

	err := db.Insert(t.Context(), &user.User{ID: 1})
	assert.NoError(t, err)
	err = db.Insert(t.Context(), &auth.OAuth2Application{Name: "test-oauth2-app-1", ClientID: "client-id-1"})
	assert.NoError(t, err)
	err = db.Insert(t.Context(), &auth.OAuth2Application{Name: "test-oauth2-app-2", ClientID: "client-id-2", UID: 1})
	assert.NoError(t, err)
	err = db.Insert(t.Context(), &auth.OAuth2Application{Name: "test-oauth2-app-3", ClientID: "client-id-3", UID: 99999999})
	assert.NoError(t, err)

	unittest.AssertExistsAndLoadBean(t, &auth.OAuth2Application{ClientID: "client-id-1"})
	unittest.AssertExistsAndLoadBean(t, &auth.OAuth2Application{ClientID: "client-id-2"})
	unittest.AssertExistsAndLoadBean(t, &auth.OAuth2Application{ClientID: "client-id-3"})

	oauth2AppCheck := checks[idx]
	err = oauth2AppCheck.Run(t.Context(), log.GetManager().GetLogger(log.DEFAULT), true)
	assert.NoError(t, err)

	unittest.AssertExistsAndLoadBean(t, &auth.OAuth2Application{ClientID: "client-id-1"})
	unittest.AssertExistsAndLoadBean(t, &auth.OAuth2Application{ClientID: "client-id-2"})
	unittest.AssertNotExistsBean(t, &auth.OAuth2Application{ClientID: "client-id-3"})
}
