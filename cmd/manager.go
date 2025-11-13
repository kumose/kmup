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
	"os"
	"time"

	"github.com/kumose/kmup/modules/private"

	"github.com/urfave/cli/v3"
)

var (
	// CmdManager represents the manager command
	CmdManager = &cli.Command{
		Name:        "manager",
		Usage:       "Manage the running kmup process",
		Description: "This is a command for managing the running kmup process",
		Commands: []*cli.Command{
			subcmdShutdown,
			subcmdRestart,
			subcmdReloadTemplates,
			subcmdFlushQueues,
			subcmdLogging,
			subCmdProcesses,
		},
	}
	subcmdShutdown = &cli.Command{
		Name:  "shutdown",
		Usage: "Gracefully shutdown the running process",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
			},
		},
		Action: runShutdown,
	}
	subcmdRestart = &cli.Command{
		Name:  "restart",
		Usage: "Gracefully restart the running process - (not implemented for windows servers)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
			},
		},
		Action: runRestart,
	}
	subcmdReloadTemplates = &cli.Command{
		Name:  "reload-templates",
		Usage: "Reload template files in the running process",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
			},
		},
		Action: runReloadTemplates,
	}
	subcmdFlushQueues = &cli.Command{
		Name:   "flush-queues",
		Usage:  "Flush queues in the running process",
		Action: runFlushQueues,
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:  "timeout",
				Value: 60 * time.Second,
				Usage: "Timeout for the flushing process",
			},
			&cli.BoolFlag{
				Name:  "non-blocking",
				Usage: "Set to true to not wait for flush to complete before returning",
			},
			&cli.BoolFlag{
				Name: "debug",
			},
		},
	}
	subCmdProcesses = &cli.Command{
		Name:   "processes",
		Usage:  "Display running processes within the current process",
		Action: runProcesses,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
			},
			&cli.BoolFlag{
				Name:  "flat",
				Usage: "Show processes as flat table rather than as tree",
			},
			&cli.BoolFlag{
				Name:  "no-system",
				Usage: "Do not show system processes",
			},
			&cli.BoolFlag{
				Name:  "stacktraces",
				Usage: "Show stacktraces",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as json",
			},
			&cli.StringFlag{
				Name:  "cancel",
				Usage: "Process PID to cancel. (Only available for non-system processes.)",
			},
		},
	}
)

func runShutdown(ctx context.Context, c *cli.Command) error {
	setup(ctx, c.Bool("debug"))
	extra := private.Shutdown(ctx)
	return handleCliResponseExtra(extra)
}

func runRestart(ctx context.Context, c *cli.Command) error {
	setup(ctx, c.Bool("debug"))
	extra := private.Restart(ctx)
	return handleCliResponseExtra(extra)
}

func runReloadTemplates(ctx context.Context, c *cli.Command) error {
	setup(ctx, c.Bool("debug"))
	extra := private.ReloadTemplates(ctx)
	return handleCliResponseExtra(extra)
}

func runFlushQueues(ctx context.Context, c *cli.Command) error {
	setup(ctx, c.Bool("debug"))
	extra := private.FlushQueues(ctx, c.Duration("timeout"), c.Bool("non-blocking"))
	return handleCliResponseExtra(extra)
}

func runProcesses(ctx context.Context, c *cli.Command) error {
	setup(ctx, c.Bool("debug"))
	extra := private.Processes(ctx, os.Stdout, c.Bool("flat"), c.Bool("no-system"), c.Bool("stacktraces"), c.Bool("json"), c.String("cancel"))
	return handleCliResponseExtra(extra)
}
