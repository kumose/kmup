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

package asymkey

import (
	"os"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
)

func GetDisplaySigningKey(key *git.SigningKey) string {
	if key == nil || key.Format == "" {
		return ""
	}

	switch key.Format {
	case git.SigningKeyFormatOpenPGP:
		return key.KeyID
	case git.SigningKeyFormatSSH:
		content, err := os.ReadFile(key.KeyID)
		if err != nil {
			log.Error("Unable to read SSH key %s: %v", key.KeyID, err)
			return "(Unable to read SSH key)"
		}
		display, err := CalcFingerprint(string(content))
		if err != nil {
			log.Error("Unable to calculate fingerprint for SSH key %s: %v", key.KeyID, err)
			return "(Unable to calculate fingerprint for SSH key)"
		}
		return display
	}
	setting.PanicInDevOrTesting("Unknown signing key format: %s", key.Format)
	return "(Unknown key format)"
}
