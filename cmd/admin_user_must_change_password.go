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
	"errors"
	"fmt"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"

	"github.com/urfave/cli/v3"
)

func microcmdUserMustChangePassword() *cli.Command {
	return &cli.Command{
		Name:   "must-change-password",
		Usage:  "Set the must change password flag for the provided users or all users",
		Action: runMustChangePassword,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"A"},
				Usage:   "All users must change password, except those explicitly excluded with --exclude",
			},
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"e"},
				Usage:   "Do not change the must-change-password flag for these users",
			},
			&cli.BoolFlag{
				Name:  "unset",
				Usage: "Instead of setting the must-change-password flag, unset it",
			},
		},
	}
}

func runMustChangePassword(ctx context.Context, c *cli.Command) error {
	if c.NArg() == 0 && !c.IsSet("all") {
		return errors.New("either usernames or --all must be provided")
	}

	mustChangePassword := !c.Bool("unset")
	all := c.Bool("all")
	exclude := c.StringSlice("exclude")

	if !setting.IsInTesting {
		if err := initDB(ctx); err != nil {
			return err
		}
	}

	n, err := user_model.SetMustChangePassword(ctx, all, mustChangePassword, c.Args().Slice(), exclude)
	if err != nil {
		return err
	}

	// codeql[disable-next-line=go/clear-text-logging]
	fmt.Printf("Updated %d users setting MustChangePassword to %t\n", n, mustChangePassword)
	return nil
}
