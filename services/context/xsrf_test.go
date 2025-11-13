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

package context

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	key      = "quay"
	userID   = "12345678"
	actionID = "POST /form"
)

var (
	now              = time.Now()
	oneMinuteFromNow = now.Add(1 * time.Minute)
)

func Test_ValidToken(t *testing.T) {
	t.Run("Validate token", func(t *testing.T) {
		tok := GenerateCsrfToken(key, userID, actionID, now)
		assert.True(t, ValidCsrfToken(tok, key, userID, actionID, oneMinuteFromNow))
		assert.True(t, ValidCsrfToken(tok, key, userID, actionID, now.Add(CsrfTokenTimeout-1*time.Nanosecond)))
		assert.True(t, ValidCsrfToken(tok, key, userID, actionID, now.Add(-1*time.Minute)))
	})
}

// Test_SeparatorReplacement tests that separators are being correctly substituted
func Test_SeparatorReplacement(t *testing.T) {
	t.Run("Test two separator replacements", func(t *testing.T) {
		assert.NotEqual(t, GenerateCsrfToken("foo:bar", "baz", "wah", now),
			GenerateCsrfToken("foo", "bar:baz", "wah", now))
	})
}

func Test_InvalidToken(t *testing.T) {
	t.Run("Test invalid tokens", func(t *testing.T) {
		invalidTokenTests := []struct {
			name, key, userID, actionID string
			t                           time.Time
		}{
			{"Bad key", "foobar", userID, actionID, oneMinuteFromNow},
			{"Bad userID", key, "foobar", actionID, oneMinuteFromNow},
			{"Bad actionID", key, userID, "foobar", oneMinuteFromNow},
			{"Expired", key, userID, actionID, now.Add(CsrfTokenTimeout)},
			{"More than 1 minute from the future", key, userID, actionID, now.Add(-1*time.Nanosecond - 1*time.Minute)},
		}

		tok := GenerateCsrfToken(key, userID, actionID, now)
		for _, itt := range invalidTokenTests {
			assert.False(t, ValidCsrfToken(tok, itt.key, itt.userID, itt.actionID, itt.t))
		}
	})
}

// Test_ValidateBadData primarily tests that no unexpected panics are triggered during parsing
func Test_ValidateBadData(t *testing.T) {
	t.Run("Validate bad data", func(t *testing.T) {
		badDataTests := []struct {
			name, tok string
		}{
			{"Invalid Base64", "ASDab24(@)$*=="},
			{"No delimiter", base64.URLEncoding.EncodeToString([]byte("foobar12345678"))},
			{"Invalid time", base64.URLEncoding.EncodeToString([]byte("foobar:foobar"))},
		}

		for _, bdt := range badDataTests {
			assert.False(t, ValidCsrfToken(bdt.tok, key, userID, actionID, oneMinuteFromNow))
		}
	})
}
