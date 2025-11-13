// Copyright 2014 The Gogs Authors. All rights reserved.
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

package main

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/kumose/kmup/cmd"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	// register supported doc types
	_ "github.com/kumose/kmup/modules/markup/asciicast"
	_ "github.com/kumose/kmup/modules/markup/console"
	_ "github.com/kumose/kmup/modules/markup/csv"
	_ "github.com/kumose/kmup/modules/markup/markdown"
	_ "github.com/kumose/kmup/modules/markup/orgmode"

	"github.com/urfave/cli/v3"
)

// these flags will be set by the build flags
var (
	Version     = "development" // program version for this build
	Tags        = ""            // the Golang build tags
	MakeVersion = ""            // "make" program version if built with make
)

func init() {
	setting.AppVer = Version
	setting.AppBuiltWith = formatBuiltWith()
	setting.AppStartTime = time.Now().UTC()
}

func main() {
	cli.OsExiter = func(code int) {
		log.GetManager().Close()
		os.Exit(code)
	}
	app := cmd.NewMainApp(cmd.AppVersion{Version: Version, Extra: formatBuiltWith()})
	_ = cmd.RunMainApp(app, os.Args...) // all errors should have been handled by the RunMainApp
	// flush the queued logs before exiting, it is a MUST, otherwise there will be log loss
	log.GetManager().Close()
}

func formatBuiltWith() string {
	version := runtime.Version()
	if len(MakeVersion) > 0 {
		version = MakeVersion + ", " + runtime.Version()
	}
	if len(Tags) == 0 {
		return " built with " + version
	}

	return " built with " + version + " : " + strings.ReplaceAll(Tags, " ", ", ")
}
