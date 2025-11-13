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
	"time"

	"github.com/yohcop/openid-go"
)

// For the demo, we use in-memory infinite storage nonce and discovery
// cache. In your app, do not use this as it will eat up memory and
// never
// free it. Use your own implementation, on a better database system.
// If you have multiple servers for example, you may need to share at
// least
// the nonceStore between them.
var (
	nonceStore     = openid.NewSimpleNonceStore()
	discoveryCache = newTimedDiscoveryCache(24 * time.Hour)
)

// Verify handles response from OpenID provider
func Verify(fullURL string) (id string, err error) {
	return openid.Verify(fullURL, discoveryCache, nonceStore)
}

// Normalize normalizes an OpenID URI
func Normalize(url string) (id string, err error) {
	return openid.Normalize(url)
}

// RedirectURL redirects browser
func RedirectURL(id, callbackURL, realm string) (string, error) {
	return openid.RedirectURL(id, callbackURL, realm)
}
