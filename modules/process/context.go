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

package process

import (
	"context"
)

// Context is a wrapper around context.Context and contains the current pid for this context
type Context struct {
	context.Context
	pid IDType
}

// GetPID returns the PID for this context
func (c *Context) GetPID() IDType {
	return c.pid
}

// GetParent returns the parent process context (if any)
func (c *Context) GetParent() *Context {
	return GetContext(c.Context)
}

// Value is part of the interface for context.Context. We mostly defer to the internal context - but we return this in response to the ProcessContextKey
func (c *Context) Value(key any) any {
	if key == ProcessContextKey {
		return c
	}
	return c.Context.Value(key)
}

// ProcessContextKey is the key under which process contexts are stored
var ProcessContextKey any = "process_context"

// GetContext will return a process context if one exists
func GetContext(ctx context.Context) *Context {
	if pCtx, ok := ctx.(*Context); ok {
		return pCtx
	}
	pCtxInterface := ctx.Value(ProcessContextKey)
	if pCtxInterface == nil {
		return nil
	}
	if pCtx, ok := pCtxInterface.(*Context); ok {
		return pCtx
	}
	return nil
}

// GetPID returns the PID for this context
func GetPID(ctx context.Context) IDType {
	pCtx := GetContext(ctx)
	if pCtx == nil {
		return ""
	}
	return pCtx.GetPID()
}

// GetParentPID returns the ParentPID for this context
func GetParentPID(ctx context.Context) IDType {
	var parentPID IDType
	if parentProcess := GetContext(ctx); parentProcess != nil {
		parentPID = parentProcess.GetPID()
	}
	return parentPID
}
