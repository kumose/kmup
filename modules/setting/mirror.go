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

// Mirror settings
var Mirror = struct {
	Enabled         bool
	DisableNewPull  bool
	DisableNewPush  bool
	DefaultInterval time.Duration
	MinInterval     time.Duration
}{
	Enabled:         true,
	DisableNewPull:  false,
	DisableNewPush:  false,
	MinInterval:     10 * time.Minute,
	DefaultInterval: 8 * time.Hour,
}

func loadMirrorFrom(rootCfg ConfigProvider) {
	// Handle old configuration through `[repository]` `DISABLE_MIRRORS`
	// - please note this was badly named and only disabled the creation of new pull mirrors
	// DEPRECATED should not be removed because users maybe upgrade from lower version to the latest version
	// if these are removed, the warning will not be shown
	deprecatedSetting(rootCfg, "repository", "DISABLE_MIRRORS", "mirror", "ENABLED", "v1.19.0")
	if ConfigSectionKeyBool(rootCfg.Section("repository"), "DISABLE_MIRRORS") {
		Mirror.DisableNewPull = true
	}

	if err := rootCfg.Section("mirror").MapTo(&Mirror); err != nil {
		log.Fatal("Failed to map Mirror settings: %v", err)
	}

	if !Mirror.Enabled {
		Mirror.DisableNewPull = true
		Mirror.DisableNewPush = true
	}

	if Mirror.MinInterval.Minutes() < 1 {
		log.Warn("Mirror.MinInterval is too low, set to 1 minute")
		Mirror.MinInterval = 1 * time.Minute
	}
	if Mirror.DefaultInterval < Mirror.MinInterval {
		Mirror.DefaultInterval = max(time.Hour*8, Mirror.MinInterval)
		log.Warn("Mirror.DefaultInterval is less than Mirror.MinInterval, set to %s", Mirror.DefaultInterval.String())
	}
}
