// Copyright 2016 The Gogs Authors. All rights reserved.
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

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kumose/kmup/modules/generate"

	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v3"
)

var (
	// CmdGenerate represents the available generate sub-command.
	CmdGenerate = &cli.Command{
		Name:  "generate",
		Usage: "Generate Kmup's secrets/keys/tokens",
		Commands: []*cli.Command{
			subcmdSecret,
		},
	}

	subcmdSecret = &cli.Command{
		Name:  "secret",
		Usage: "Generate a secret token",
		Commands: []*cli.Command{
			microcmdGenerateInternalToken,
			microcmdGenerateLfsJwtSecret,
			microcmdGenerateSecretKey,
		},
	}

	microcmdGenerateInternalToken = &cli.Command{
		Name:   "INTERNAL_TOKEN",
		Usage:  "Generate a new INTERNAL_TOKEN",
		Action: runGenerateInternalToken,
	}

	microcmdGenerateLfsJwtSecret = &cli.Command{
		Name:    "JWT_SECRET",
		Aliases: []string{"LFS_JWT_SECRET"},
		Usage:   "Generate a new JWT_SECRET",
		Action:  runGenerateLfsJwtSecret,
	}

	microcmdGenerateSecretKey = &cli.Command{
		Name:   "SECRET_KEY",
		Usage:  "Generate a new SECRET_KEY",
		Action: runGenerateSecretKey,
	}
)

func runGenerateInternalToken(_ context.Context, c *cli.Command) error {
	internalToken, err := generate.NewInternalToken()
	if err != nil {
		return err
	}

	fmt.Printf("%s", internalToken)

	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("\n")
	}

	return nil
}

func runGenerateLfsJwtSecret(_ context.Context, c *cli.Command) error {
	_, jwtSecretBase64, err := generate.NewJwtSecretWithBase64()
	if err != nil {
		return err
	}

	fmt.Printf("%s", jwtSecretBase64)

	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("\n")
	}

	return nil
}

func runGenerateSecretKey(_ context.Context, c *cli.Command) error {
	secretKey, err := generate.NewSecretKey()
	if err != nil {
		return err
	}

	// codeql[disable-next-line=go/clear-text-logging]
	fmt.Printf("%s", secretKey)

	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("\n")
	}

	return nil
}
