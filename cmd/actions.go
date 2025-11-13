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

	"github.com/kumose/kmup/modules/private"
	"github.com/kumose/kmup/modules/setting"

	"github.com/urfave/cli/v3"
)

var (
	// CmdActions represents the available actions sub-commands.
	CmdActions = &cli.Command{
		Name:  "actions",
		Usage: "Manage Kmup Actions",
		Commands: []*cli.Command{
			subcmdActionsGenRunnerToken,
		},
	}

	subcmdActionsGenRunnerToken = &cli.Command{
		Name:    "generate-runner-token",
		Usage:   "Generate a new token for a runner to use to register with the server",
		Action:  runGenerateActionsRunnerToken,
		Aliases: []string{"grt"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "scope",
				Aliases: []string{"s"},
				Value:   "",
				Usage:   "{owner}[/{repo}] - leave empty for a global runner",
			},
		},
	}
)

func runGenerateActionsRunnerToken(ctx context.Context, c *cli.Command) error {
	setting.MustInstalled()

	scope := c.String("scope")

	respText, extra := private.GenerateActionsRunnerToken(ctx, scope)
	if extra.HasError() {
		return handleCliResponseExtra(extra)
	}
	_, _ = fmt.Printf("%s\n", respText.Text)
	return nil
}
