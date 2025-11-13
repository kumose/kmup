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

package ssh

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

func Init() error {
	if setting.SSH.Disabled {
		builtinUnused()
		return nil
	}

	if setting.SSH.StartBuiltinServer {
		Listen(setting.SSH.ListenHost, setting.SSH.ListenPort, setting.SSH.ServerCiphers, setting.SSH.ServerKeyExchanges, setting.SSH.ServerMACs)
		log.Info("SSH server started on %q. Ciphers: %v, key exchange algorithms: %v, MACs: %v",
			net.JoinHostPort(setting.SSH.ListenHost, strconv.Itoa(setting.SSH.ListenPort)),
			util.Iif[any](setting.SSH.ServerCiphers == nil, "default", setting.SSH.ServerCiphers),
			util.Iif[any](setting.SSH.ServerKeyExchanges == nil, "default", setting.SSH.ServerKeyExchanges),
			util.Iif[any](setting.SSH.ServerMACs == nil, "default", setting.SSH.ServerMACs),
		)
		return nil
	}

	builtinUnused()

	if len(setting.SSH.TrustedUserCAKeys) > 0 && setting.SSH.AuthorizedPrincipalsEnabled {
		caKeysFileName := setting.SSH.TrustedUserCAKeysFile
		caKeysFileDir := filepath.Dir(caKeysFileName)

		err := os.MkdirAll(caKeysFileDir, 0o700) // SSH.RootPath by default (That is `~/.ssh` in most cases)
		if err != nil {
			return fmt.Errorf("failed to create directory %q for ssh trusted ca keys: %w", caKeysFileDir, err)
		}

		if err := os.WriteFile(caKeysFileName, []byte(strings.Join(setting.SSH.TrustedUserCAKeys, "\n")), 0o600); err != nil {
			return fmt.Errorf("failed to write ssh trusted ca keys to %q: %w", caKeysFileName, err)
		}
	}

	return nil
}
