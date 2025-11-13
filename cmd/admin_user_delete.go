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
	"strings"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/storage"
	user_service "github.com/kumose/kmup/services/user"

	"github.com/urfave/cli/v3"
)

func microcmdUserDelete() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete specific user by id, name or email",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "id",
				Usage: "ID of user of the user to delete",
			},
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "Username of the user to delete",
			},
			&cli.StringFlag{
				Name:    "email",
				Aliases: []string{"e"},
				Usage:   "Email of the user to delete",
			},
			&cli.BoolFlag{
				Name:  "purge",
				Usage: "Purge user, all their repositories, organizations and comments",
			},
		},
		Action: runDeleteUser,
	}
}

func runDeleteUser(ctx context.Context, c *cli.Command) error {
	if !c.IsSet("id") && !c.IsSet("username") && !c.IsSet("email") {
		return errors.New("You must provide the id, username or email of a user to delete")
	}

	if !setting.IsInTesting {
		if err := initDB(ctx); err != nil {
			return err
		}
	}

	if err := storage.Init(); err != nil {
		return err
	}

	var err error
	var user *user_model.User
	if c.IsSet("email") {
		user, err = user_model.GetUserByEmail(ctx, c.String("email"))
	} else if c.IsSet("username") {
		user, err = user_model.GetUserByName(ctx, c.String("username"))
	} else {
		user, err = user_model.GetUserByID(ctx, c.Int64("id"))
	}
	if err != nil {
		return err
	}
	if c.IsSet("username") && user.LowerName != strings.ToLower(strings.TrimSpace(c.String("username"))) {
		return fmt.Errorf("the user %s who has email %s does not match the provided username %s", user.Name, c.String("email"), c.String("username"))
	}

	if c.IsSet("id") && user.ID != c.Int64("id") {
		return fmt.Errorf("the user %s does not match the provided id %d", user.Name, c.Int64("id"))
	}

	return user_service.DeleteUser(ctx, user, c.Bool("purge"))
}
