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

package setting

import (
	"net/url"
	"path"

	"github.com/kumose/kmup/modules/log"
)

// API settings
var API = struct {
	EnableSwagger          bool
	SwaggerURL             string
	MaxResponseItems       int
	DefaultPagingNum       int
	DefaultGitTreesPerPage int
	DefaultMaxBlobSize     int64
	DefaultMaxResponseSize int64
}{
	EnableSwagger:          true,
	SwaggerURL:             "",
	MaxResponseItems:       50,
	DefaultPagingNum:       30,
	DefaultGitTreesPerPage: 1000,
	DefaultMaxBlobSize:     10485760,
	DefaultMaxResponseSize: 104857600,
}

func loadAPIFrom(rootCfg ConfigProvider) {
	mustMapSetting(rootCfg, "api", &API)

	defaultAppURL := string(Protocol) + "://" + Domain + ":" + HTTPPort
	u, err := url.Parse(rootCfg.Section("server").Key("ROOT_URL").MustString(defaultAppURL))
	if err != nil {
		log.Fatal("Invalid ROOT_URL '%s': %s", AppURL, err)
	}
	u.Path = path.Join(u.Path, "api", "swagger")
	API.SwaggerURL = u.String()
}
