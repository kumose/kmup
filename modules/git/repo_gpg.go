// Copyright 2015 The Gogs Authors. All rights reserved.
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

package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/process"
)

// LoadPublicKeyContent will load the key from gpg
func (gpgSettings *GPGSettings) LoadPublicKeyContent() error {
	if gpgSettings.Format == SigningKeyFormatSSH {
		content, err := os.ReadFile(gpgSettings.KeyID)
		if err != nil {
			return fmt.Errorf("unable to read SSH public key file: %s, %w", gpgSettings.KeyID, err)
		}
		gpgSettings.PublicKeyContent = string(content)
		return nil
	}
	content, stderr, err := process.GetManager().Exec(
		"gpg -a --export",
		"gpg", "-a", "--export", gpgSettings.KeyID)
	if err != nil {
		return fmt.Errorf("unable to get default signing key: %s, %s, %w", gpgSettings.KeyID, stderr, err)
	}
	gpgSettings.PublicKeyContent = content
	return nil
}

// GetDefaultPublicGPGKey will return and cache the default public GPG settings for this repository
func (repo *Repository) GetDefaultPublicGPGKey(forceUpdate bool) (*GPGSettings, error) {
	if repo.gpgSettings != nil && !forceUpdate {
		return repo.gpgSettings, nil
	}

	gpgSettings := &GPGSettings{
		Sign: true,
	}

	value, _, _ := gitcmd.NewCommand("config", "--get", "commit.gpgsign").WithDir(repo.Path).RunStdString(repo.Ctx)
	sign, valid := ParseBool(strings.TrimSpace(value))
	if !sign || !valid {
		gpgSettings.Sign = false
		repo.gpgSettings = gpgSettings
		return gpgSettings, nil
	}

	signingKey, _, _ := gitcmd.NewCommand("config", "--get", "user.signingkey").WithDir(repo.Path).RunStdString(repo.Ctx)
	gpgSettings.KeyID = strings.TrimSpace(signingKey)

	format, _, _ := gitcmd.NewCommand("config", "--default", SigningKeyFormatOpenPGP, "--get", "gpg.format").WithDir(repo.Path).RunStdString(repo.Ctx)
	gpgSettings.Format = strings.TrimSpace(format)

	defaultEmail, _, _ := gitcmd.NewCommand("config", "--get", "user.email").WithDir(repo.Path).RunStdString(repo.Ctx)
	gpgSettings.Email = strings.TrimSpace(defaultEmail)

	defaultName, _, _ := gitcmd.NewCommand("config", "--get", "user.name").WithDir(repo.Path).RunStdString(repo.Ctx)
	gpgSettings.Name = strings.TrimSpace(defaultName)

	if err := gpgSettings.LoadPublicKeyContent(); err != nil {
		return nil, err
	}
	repo.gpgSettings = gpgSettings
	return repo.gpgSettings, nil
}
