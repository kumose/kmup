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
	"github.com/kumose/kmup/modules/auth/password"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"
	user_service "github.com/kumose/kmup/services/user"

	"github.com/urfave/cli/v3"
)

func microcmdUserChangePassword() *cli.Command {
	return &cli.Command{
		Name:   "change-password",
		Usage:  "Change a user's password",
		Action: runChangePassword,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Usage:    "The user to change password for",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "New password to set for user",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "must-change-password",
				Usage: "User must change password (can be disabled by --must-change-password=false)",
				Value: true,
			},
		},
	}
}

func runChangePassword(ctx context.Context, c *cli.Command) error {
	if !setting.IsInTesting {
		if err := initDB(ctx); err != nil {
			return err
		}
	}

	user, err := user_model.GetUserByName(ctx, c.String("username"))
	if err != nil {
		return err
	}

	opts := &user_service.UpdateAuthOptions{
		Password:           optional.Some(c.String("password")),
		MustChangePassword: optional.Some(c.Bool("must-change-password")),
	}
	if err := user_service.UpdateAuth(ctx, user, opts); err != nil {
		switch {
		case errors.Is(err, password.ErrMinLength):
			return fmt.Errorf("password is not long enough, needs to be at least %d characters", setting.MinPasswordLength)
		case errors.Is(err, password.ErrComplexity):
			return errors.New("password does not meet complexity requirements")
		case errors.Is(err, password.ErrIsPwned):
			return errors.New("the password is in a list of stolen passwords previously exposed in public data breaches, please try again with a different password, to see more details: https://haveibeenpwned.com/Passwords")
		default:
			return err
		}
	}

	fmt.Printf("%s's password has been successfully updated!\n", user.Name)
	return nil
}
