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

package gtprof

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kumose/kmup/modules/tailmsg"
)

type traceBuiltinStarter struct{}

type traceBuiltinSpan struct {
	ts *TraceSpan

	internalSpanIdx int
}

func (t *traceBuiltinSpan) addEvent(name string, cfg *EventConfig) {
	// No-op because builtin tracer doesn't need it.
	// In the future we might use it to mark the time point between backend logic and network response.
}

func (t *traceBuiltinSpan) recordError(err error, cfg *EventConfig) {
	// No-op because builtin tracer doesn't need it.
	// Actually Kmup doesn't handle err this way in most cases
}

func (t *traceBuiltinSpan) toString(out *strings.Builder, indent int) {
	t.ts.mu.RLock()
	defer t.ts.mu.RUnlock()

	out.WriteString(strings.Repeat(" ", indent))
	out.WriteString(t.ts.name)
	if t.ts.endTime.IsZero() {
		out.WriteString(" duration: (not ended)")
	} else {
		fmt.Fprintf(out, " start=%s duration=%.4fs", t.ts.startTime.Format("2006-01-02 15:04:05"), t.ts.endTime.Sub(t.ts.startTime).Seconds())
	}
	for _, a := range t.ts.attributes {
		out.WriteString(" ")
		out.WriteString(a.Key)
		out.WriteString("=")
		value := a.Value.AsString()
		if strings.ContainsAny(value, " \t\r\n") {
			quoted := false
			for _, c := range "\"'`" {
				if quoted = !strings.Contains(value, string(c)); quoted {
					value = string(c) + value + string(c)
					break
				}
			}
			if !quoted {
				value = fmt.Sprintf("%q", value)
			}
		}
		out.WriteString(value)
	}
	out.WriteString("\n")
	for _, c := range t.ts.children {
		span := c.internalSpans[t.internalSpanIdx].(*traceBuiltinSpan)
		span.toString(out, indent+2)
	}
}

func (t *traceBuiltinSpan) end() {
	if t.ts.parent == nil {
		// TODO: debug purpose only
		// TODO: it should distinguish between http response network lag and actual processing time
		threshold := time.Duration(traceBuiltinThreshold.Load())
		if threshold != 0 && t.ts.endTime.Sub(t.ts.startTime) > threshold {
			sb := &strings.Builder{}
			t.toString(sb, 0)
			tailmsg.GetManager().GetTraceRecorder().Record(sb.String())
		}
	}
}

func (t *traceBuiltinStarter) start(ctx context.Context, traceSpan *TraceSpan, internalSpanIdx int) (context.Context, traceSpanInternal) {
	return ctx, &traceBuiltinSpan{ts: traceSpan, internalSpanIdx: internalSpanIdx}
}

func init() {
	globalTraceStarters = append(globalTraceStarters, &traceBuiltinStarter{})
}

var traceBuiltinThreshold atomic.Int64

func EnableBuiltinTracer(threshold time.Duration) {
	traceBuiltinThreshold.Store(int64(threshold))
}
