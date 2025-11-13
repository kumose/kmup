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
	"text/tabwriter"

	user_model "github.com/kumose/kmup/models/user"

	"github.com/urfave/cli/v3"
)

var microcmdUserList = &cli.Command{
	Name:   "list",
	Usage:  "List users",
	Action: runListUsers,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "admin",
			Usage: "List only admin users",
		},
	},
}

func runListUsers(ctx context.Context, c *cli.Command) error {
	if err := initDB(ctx); err != nil {
		return err
	}

	users, err := user_model.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 0, 1, ' ', 0)

	if c.IsSet("admin") {
		fmt.Fprintf(w, "ID\tUsername\tEmail\tIsActive\n")
		for _, u := range users {
			if u.IsAdmin {
				fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", u.ID, u.Name, u.Email, u.IsActive)
			}
		}
	} else {
		twofa := user_model.UserList(users).GetTwoFaStatus(ctx)
		fmt.Fprintf(w, "ID\tUsername\tEmail\tIsActive\tIsAdmin\t2FA\n")
		for _, u := range users {
			fmt.Fprintf(w, "%d\t%s\t%s\t%t\t%t\t%t\n", u.ID, u.Name, u.Email, u.IsActive, u.IsAdmin, twofa[u.ID])
		}
	}

	w.Flush()
	return nil
}
