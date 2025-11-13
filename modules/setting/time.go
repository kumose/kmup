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
	"time"

	"github.com/kumose/kmup/modules/log"
)

// DefaultUILocation is the location on the UI, so that we can display the time on UI.
var DefaultUILocation = time.Local

func loadTimeFrom(rootCfg ConfigProvider) {
	zone := rootCfg.Section("time").Key("DEFAULT_UI_LOCATION").String()
	if zone != "" {
		var err error
		DefaultUILocation, err = time.LoadLocation(zone)
		if err != nil {
			log.Fatal("Load time zone failed: %v", err)
		}
	}
	if DefaultUILocation == nil {
		DefaultUILocation = time.Local
	}
}
