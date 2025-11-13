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

package log

import (
	"context"
	"runtime"
	"strings"

	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/util/rotatingfilewriter"
)

var projectPackagePrefix string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	projectPackagePrefix = strings.TrimSuffix(filename, "modules/log/init.go")
	if projectPackagePrefix == filename {
		// in case the source code file is moved, we can not trim the suffix, the code above should also be updated.
		panic("unable to detect correct package prefix, please update file: " + filename)
	}

	rotatingfilewriter.ErrorPrintf = FallbackErrorf

	process.TraceCallback = func(skip int, start bool, pid process.IDType, description string, parentPID process.IDType, typ string) {
		if start && parentPID != "" {
			Log(skip+1, TRACE, "Start %s: %s (from %s) (%s)", NewColoredValue(pid, FgHiYellow), description, NewColoredValue(parentPID, FgYellow), NewColoredValue(typ, Reset))
		} else if start {
			Log(skip+1, TRACE, "Start %s: %s (%s)", NewColoredValue(pid, FgHiYellow), description, NewColoredValue(typ, Reset))
		} else {
			Log(skip+1, TRACE, "Done %s: %s", NewColoredValue(pid, FgHiYellow), NewColoredValue(description, Reset))
		}
	}
}

func newProcessTypedContext(parent context.Context, desc string) (ctx context.Context, cancel context.CancelFunc) {
	// the "process manager" also calls "log.Trace()" to output logs, so if we want to create new contexts by the manager, we need to disable the trace temporarily
	process.TraceLogDisable(true)
	defer process.TraceLogDisable(false)
	ctx, _, cancel = process.GetManager().AddTypedContext(parent, desc, process.SystemProcessType, false)
	return ctx, cancel
}
