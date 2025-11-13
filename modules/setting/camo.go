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
	"strconv"

	"github.com/kumose/kmup/modules/log"
)

var Camo = struct {
	Enabled   bool
	ServerURL string `ini:"SERVER_URL"`
	HMACKey   string `ini:"HMAC_KEY"`
	Always    bool
}{}

func loadCamoFrom(rootCfg ConfigProvider) {
	mustMapSetting(rootCfg, "camo", &Camo)
	if Camo.Enabled {
		oldValue := rootCfg.Section("camo").Key("ALLWAYS").MustString("")
		if oldValue != "" {
			log.Warn("camo.ALLWAYS is deprecated, use camo.ALWAYS instead")
			Camo.Always, _ = strconv.ParseBool(oldValue)
		}

		if Camo.ServerURL == "" || Camo.HMACKey == "" {
			log.Fatal(`Camo settings require "SERVER_URL" and HMAC_KEY`)
		}
	}
}
