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
	"io"
	"os"
	"strings"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	"github.com/urfave/cli/v3"
)

var cliHelpPrinterOld = cli.HelpPrinter

func init() {
	cli.HelpPrinter = cliHelpPrinterNew
}

// cliHelpPrinterNew helps to print "DEFAULT CONFIGURATION" for the following cases ( "-c" can apper in any position):
// * ./kmup -c /dev/null -h
// * ./kmup -c help /dev/null help
// * ./kmup help -c /dev/null
// * ./kmup help -c /dev/null web
// * ./kmup help web -c /dev/null
// * ./kmup web help -c /dev/null
// * ./kmup web -h -c /dev/null
func cliHelpPrinterNew(out io.Writer, templ string, data any) {
	cmd, _ := data.(*cli.Command)
	if cmd != nil {
		prepareWorkPathAndCustomConf(cmd)
	}
	cliHelpPrinterOld(out, templ, data)
	if setting.CustomConf != "" {
		_, _ = fmt.Fprintf(out, `
DEFAULT CONFIGURATION:
   AppPath:    %s
   WorkPath:   %s
   CustomPath: %s
   ConfigFile: %s

`, setting.AppPath, setting.AppWorkPath, setting.CustomPath, setting.CustomConf)
	}
}

func prepareSubcommandWithGlobalFlags(originCmd *cli.Command) {
	originBefore := originCmd.Before
	originCmd.Before = func(ctxOrig context.Context, cmd *cli.Command) (ctx context.Context, err error) {
		ctx = ctxOrig
		if originBefore != nil {
			ctx, err = originBefore(ctx, cmd)
			if err != nil {
				return ctx, err
			}
		}
		prepareWorkPathAndCustomConf(cmd)
		return ctx, nil
	}
}

// prepareWorkPathAndCustomConf tries to prepare the work path, custom path and custom config from various inputs:
// command line flags, environment variables, config file
func prepareWorkPathAndCustomConf(cmd *cli.Command) {
	var args setting.ArgWorkPathAndCustomConf
	if cmd.IsSet("work-path") {
		args.WorkPath = cmd.String("work-path")
	}
	if cmd.IsSet("custom-path") {
		args.CustomPath = cmd.String("custom-path")
	}
	if cmd.IsSet("config") {
		args.CustomConf = cmd.String("config")
	}
	setting.InitWorkPathAndCommonConfig(os.Getenv, args)
}

type AppVersion struct {
	Version string
	Extra   string
}

func NewMainApp(appVer AppVersion) *cli.Command {
	app := &cli.Command{}
	app.Name = "kmup" // must be lower-cased because it appears in the "USAGE" section like "kmup doctor [command [command options]]"
	app.Usage = "A painless self-hosted Git service"
	app.Description = `Kmup program contains "web" and other subcommands. If no subcommand is given, it starts the web server by default. Use "web" subcommand for more web server arguments, use other subcommands for other purposes.`
	app.Version = appVer.Version + appVer.Extra
	app.EnableShellCompletion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:      "work-path",
			Aliases:   []string{"w"},
			TakesFile: true,
			Usage:     "Set Kmup's working path (defaults to the Kmup's binary directory)",
		},
		&cli.StringFlag{
			Name:      "config",
			Aliases:   []string{"c"},
			TakesFile: true,
			Value:     setting.CustomConf,
			Usage:     "Set custom config file (defaults to '{WorkPath}/custom/conf/app.ini')",
		},
		&cli.StringFlag{
			Name:      "custom-path",
			Aliases:   []string{"C"},
			TakesFile: true,
			Usage:     "Set custom path (defaults to '{WorkPath}/custom')",
		},
	}
	// these sub-commands need to use a config file
	subCmdWithConfig := []*cli.Command{
		CmdWeb,
		CmdServ,
		CmdHook,
		CmdKeys,
		CmdDump,
		CmdAdmin,
		CmdMigrate,
		CmdDoctor,
		CmdManager,
		CmdEmbedded,
		CmdMigrateStorage,
		CmdDumpRepository,
		CmdRestoreRepository,
		CmdActions,
	}

	// these sub-commands do not need the config file, and they do not depend on any path or environment variable.
	subCmdStandalone := []*cli.Command{
		cmdConfig(),
		cmdCert(),
		CmdGenerate,
		CmdDocs,
	}

	// TODO: we should eventually drop the default command,
	// but not sure whether it would break Windows users who used to double-click the EXE to run.
	app.DefaultCommand = CmdWeb.Name

	app.Before = PrepareConsoleLoggerLevel(log.INFO)
	for i := range subCmdWithConfig {
		prepareSubcommandWithGlobalFlags(subCmdWithConfig[i])
	}
	app.Commands = append(app.Commands, subCmdWithConfig...)
	app.Commands = append(app.Commands, subCmdStandalone...)

	setting.InitKmupEnvVars()
	return app
}

func RunMainApp(app *cli.Command, args ...string) error {
	ctx, cancel := installSignals()
	defer cancel()
	err := app.Run(ctx, args)
	if err == nil {
		return nil
	}
	if strings.HasPrefix(err.Error(), "flag provided but not defined:") {
		// the cli package should already have output the error message, so just exit
		cli.OsExiter(1)
		return err
	}
	_, _ = fmt.Fprintf(app.ErrWriter, "Command error: %v\n", err)
	cli.OsExiter(1)
	return err
}
