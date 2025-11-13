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

package routing

import (
	"context"
	"net/http"

	"github.com/kumose/kmup/modules/gtprof"
	"github.com/kumose/kmup/modules/reqctx"
)

type contextKeyType struct{}

var contextKey contextKeyType

// RecordFuncInfo records a func info into context
func RecordFuncInfo(ctx context.Context, funcInfo *FuncInfo) (end func()) {
	end = func() {}
	if reqCtx := reqctx.FromContext(ctx); reqCtx != nil {
		var traceSpan *gtprof.TraceSpan
		traceSpan, end = gtprof.GetTracer().StartInContext(reqCtx, "http.func")
		traceSpan.SetAttributeString("func", funcInfo.shortName)
	}
	if record, ok := ctx.Value(contextKey).(*requestRecord); ok {
		record.lock.Lock()
		record.funcInfo = funcInfo
		record.lock.Unlock()
	}
	return end
}

// MarkLongPolling marks the request is a long-polling request, and the logger may output different message for it
func MarkLongPolling(resp http.ResponseWriter, req *http.Request) {
	record, ok := req.Context().Value(contextKey).(*requestRecord)
	if !ok {
		return
	}

	record.lock.Lock()
	record.isLongPolling = true
	record.lock.Unlock()
}

// UpdatePanicError updates a context's error info, a panic may be recovered by other middlewares, but we still need to know that.
func UpdatePanicError(ctx context.Context, err any) {
	record, ok := ctx.Value(contextKey).(*requestRecord)
	if !ok {
		return
	}

	record.lock.Lock()
	record.panicError = err
	record.lock.Unlock()
}
