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

package migrations

import (
	"crypto/tls"
	"net/http"

	"github.com/kumose/kmup/modules/hostmatcher"
	"github.com/kumose/kmup/modules/proxy"
	"github.com/kumose/kmup/modules/setting"
)

// NewMigrationHTTPClient returns a HTTP client for migration
func NewMigrationHTTPClient() *http.Client {
	return &http.Client{
		Transport: NewMigrationHTTPTransport(),
	}
}

// NewMigrationHTTPTransport returns a HTTP transport for migration
func NewMigrationHTTPTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: setting.Migrations.SkipTLSVerify},
		Proxy:           proxy.Proxy(),
		DialContext:     hostmatcher.NewDialContext("migration", allowList, blockList, setting.Proxy.ProxyURLFixed),
	}
}
