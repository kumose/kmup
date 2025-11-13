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

package doctor

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	asymkey_model "github.com/kumose/kmup/models/asymkey"
	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	asymkey_service "github.com/kumose/kmup/services/asymkey"
)

func checkAuthorizedKeys(ctx context.Context, logger log.Logger, autofix bool) error {
	if setting.SSH.StartBuiltinServer || !setting.SSH.CreateAuthorizedKeysFile {
		return nil
	}

	fPath := filepath.Join(setting.SSH.RootPath, "authorized_keys")
	f, err := os.Open(fPath)
	if err != nil {
		if !autofix {
			logger.Critical("Unable to open authorized_keys file. ERROR: %v", err)
			return fmt.Errorf("Unable to open authorized_keys file. ERROR: %w", err)
		}
		logger.Warn("Unable to open authorized_keys. (ERROR: %v). Attempting to rewrite...", err)
		if err = asymkey_service.RewriteAllPublicKeys(ctx); err != nil {
			logger.Critical("Unable to rewrite authorized_keys file. ERROR: %v", err)
			return fmt.Errorf("Unable to rewrite authorized_keys file. ERROR: %w", err)
		}
	}
	defer f.Close()

	linesInAuthorizedKeys := make(container.Set[string])

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, asymkey_model.AuthorizedStringCommentPrefix) {
			continue
		}
		linesInAuthorizedKeys.Add(line)
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	// although there is a "defer close" above, here close explicitly before the generating, because it needs to open the file for writing again
	_ = f.Close()

	// now we regenerate and check if there are any lines missing
	regenerated := &bytes.Buffer{}
	if err := asymkey_model.RegeneratePublicKeys(ctx, regenerated); err != nil {
		logger.Critical("Unable to regenerate authorized_keys file. ERROR: %v", err)
		return fmt.Errorf("Unable to regenerate authorized_keys file. ERROR: %w", err)
	}
	scanner = bufio.NewScanner(regenerated)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, asymkey_model.AuthorizedStringCommentPrefix) {
			continue
		}
		if linesInAuthorizedKeys.Contains(line) {
			continue
		}
		if !autofix {
			logger.Critical(
				"authorized_keys file %q is out of date.\nRegenerate it with:\n\t\"%s\"\nor\n\t\"%s\"",
				fPath,
				"kmup admin regenerate keys",
				"kmup doctor --run authorized-keys --fix")
			return errors.New(`authorized_keys is out of date and should be regenerated with "kmup admin regenerate keys" or "kmup doctor --run authorized-keys --fix"`)
		}
		logger.Warn("authorized_keys is out of date. Attempting rewrite...")
		err = asymkey_service.RewriteAllPublicKeys(ctx)
		if err != nil {
			logger.Critical("Unable to rewrite authorized_keys file. ERROR: %v", err)
			return fmt.Errorf("Unable to rewrite authorized_keys file. ERROR: %w", err)
		}
	}
	return nil
}

func init() {
	Register(&Check{
		Title:     "Check if OpenSSH authorized_keys file is up-to-date",
		Name:      "authorized-keys",
		IsDefault: true,
		Run:       checkAuthorizedKeys,
		Priority:  4,
	})
}
