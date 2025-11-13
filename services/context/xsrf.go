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
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CsrfTokenTimeout represents the duration that XSRF tokens are valid.
// It is exported so clients may set cookie timeouts that match generated tokens.
const CsrfTokenTimeout = 24 * time.Hour

// CsrfTokenRegenerationInterval is the interval between token generations, old tokens are still valid before CsrfTokenTimeout
var CsrfTokenRegenerationInterval = 10 * time.Minute

var csrfTokenSep = []byte(":")

// GenerateCsrfToken returns a URL-safe secure XSRF token that expires in CsrfTokenTimeout hours.
// key is a secret key for your application.
// userID is a unique identifier for the user.
// actionID is the action the user is taking (e.g. POSTing to a particular path).
func GenerateCsrfToken(key, userID, actionID string, now time.Time) string {
	nowUnixNano := now.UnixNano()
	nowUnixNanoStr := strconv.FormatInt(nowUnixNano, 10)
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(strings.ReplaceAll(userID, ":", "_")))
	h.Write(csrfTokenSep)
	h.Write([]byte(strings.ReplaceAll(actionID, ":", "_")))
	h.Write(csrfTokenSep)
	h.Write([]byte(nowUnixNanoStr))
	tok := fmt.Sprintf("%s:%s", h.Sum(nil), nowUnixNanoStr)
	return base64.RawURLEncoding.EncodeToString([]byte(tok))
}

func ParseCsrfToken(token string) (issueTime time.Time, ok bool) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, false
	}

	pos := bytes.LastIndex(data, csrfTokenSep)
	if pos == -1 {
		return time.Time{}, false
	}
	nanos, err := strconv.ParseInt(string(data[pos+1:]), 10, 64)
	if err != nil {
		return time.Time{}, false
	}
	return time.Unix(0, nanos), true
}

// ValidCsrfToken returns true if token is a valid and unexpired token returned by Generate.
func ValidCsrfToken(token, key, userID, actionID string, now time.Time) bool {
	issueTime, ok := ParseCsrfToken(token)
	if !ok {
		return false
	}

	// Check that the token is not expired.
	if now.Sub(issueTime) >= CsrfTokenTimeout {
		return false
	}

	// Check that the token is not from the future.
	// Allow 1-minute grace period in case the token is being verified on a
	// machine whose clock is behind the machine that issued the token.
	if issueTime.After(now.Add(1 * time.Minute)) {
		return false
	}

	expected := GenerateCsrfToken(key, userID, actionID, issueTime)

	// Check that the token matches the expected value.
	// Use constant time comparison to avoid timing attacks.
	return subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1
}
