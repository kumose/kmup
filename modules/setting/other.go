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

import "github.com/kumose/kmup/modules/log"

type OtherConfig struct {
	ShowFooterVersion          bool
	ShowFooterTemplateLoadTime bool
	ShowFooterPoweredBy        bool
	EnableFeed                 bool
	EnableSitemap              bool
}

var Other = OtherConfig{
	ShowFooterVersion:          true,
	ShowFooterTemplateLoadTime: true,
	ShowFooterPoweredBy:        true,
	EnableSitemap:              true,
	EnableFeed:                 true,
}

func loadOtherFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("other")
	if err := sec.MapTo(&Other); err != nil {
		log.Fatal("Failed to map [other] settings: %v", err)
	}
}
