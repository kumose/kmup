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

package db

import (
	"context"
	"time"

	"github.com/kumose/kmup/modules/gtprof"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/xorm/contexts"
)

type EngineHook struct {
	Threshold time.Duration
	Logger    log.Logger
}

var _ contexts.Hook = (*EngineHook)(nil)

func (*EngineHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	ctx, _ := gtprof.GetTracer().Start(c.Ctx, gtprof.TraceSpanDatabase)
	return ctx, nil
}

func (h *EngineHook) AfterProcess(c *contexts.ContextHook) error {
	span := gtprof.GetContextSpan(c.Ctx)
	if span != nil {
		// Do not record SQL parameters here:
		// * It shouldn't expose the parameters because they contain sensitive information, end users need to report the trace details safely.
		// * Some parameters contain quite long texts, waste memory and are difficult to display.
		span.SetAttributeString(gtprof.TraceAttrDbSQL, c.SQL)
		span.End()
	} else {
		setting.PanicInDevOrTesting("span in database engine hook is nil")
	}
	if c.ExecuteTime >= h.Threshold {
		// 8 is the amount of skips passed to runtime.Caller, so that in the log the correct function
		// is being displayed (the function that ultimately wants to execute the query in the code)
		// instead of the function of the slow query hook being called.
		h.Logger.Log(8, &log.Event{Level: log.WARN}, "[Slow SQL Query] %s %v - %v", c.SQL, c.Args, c.ExecuteTime)
	}
	return nil
}
