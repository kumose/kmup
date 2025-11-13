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

package openid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDiscoveredInfo struct{}

func (s *testDiscoveredInfo) ClaimedID() string {
	return "claimedID"
}

func (s *testDiscoveredInfo) OpEndpoint() string {
	return "opEndpoint"
}

func (s *testDiscoveredInfo) OpLocalID() string {
	return "opLocalID"
}

func TestTimedDiscoveryCache(t *testing.T) {
	ttl := 50 * time.Millisecond
	dc := newTimedDiscoveryCache(ttl)

	// Put some initial values
	dc.Put("foo", &testDiscoveredInfo{}) // openid.opEndpoint: "a", openid.opLocalID: "b", openid.claimedID: "c"})

	// Make sure we can retrieve them
	di := dc.Get("foo")
	require.NotNil(t, di)
	assert.Equal(t, "opEndpoint", di.OpEndpoint())
	assert.Equal(t, "opLocalID", di.OpLocalID())
	assert.Equal(t, "claimedID", di.ClaimedID())

	// Attempt to get a non-existent value
	assert.Nil(t, dc.Get("bar"))

	// Sleep for a while and try to retrieve again
	time.Sleep(ttl * 3 / 2)

	assert.Nil(t, dc.Get("foo"))
}
