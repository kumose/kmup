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

package runner

import (
	"context"
	"crypto/subtle"
	"errors"
	"strings"

	actions_model "github.com/kumose/kmup/models/actions"
	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	uuidHeaderKey  = "x-runner-uuid"
	tokenHeaderKey = "x-runner-token"
)

var withRunner = connect.WithInterceptors(connect.UnaryInterceptorFunc(func(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		methodName := getMethodName(request)
		if methodName == "Register" {
			return unaryFunc(ctx, request)
		}
		uuid := request.Header().Get(uuidHeaderKey)
		token := request.Header().Get(tokenHeaderKey)

		runner, err := actions_model.GetRunnerByUUID(ctx, uuid)
		if err != nil {
			if errors.Is(err, util.ErrNotExist) {
				return nil, status.Error(codes.Unauthenticated, "unregistered runner")
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		if subtle.ConstantTimeCompare([]byte(runner.TokenHash), []byte(auth_model.HashToken(token, runner.TokenSalt))) != 1 {
			return nil, status.Error(codes.Unauthenticated, "unregistered runner")
		}

		cols := []string{"last_online"}
		runner.LastOnline = timeutil.TimeStampNow()
		if methodName == "UpdateTask" || methodName == "UpdateLog" {
			runner.LastActive = timeutil.TimeStampNow()
			cols = append(cols, "last_active")
		}
		if err := actions_model.UpdateRunner(ctx, runner, cols...); err != nil {
			log.Error("can't update runner status: %v", err)
		}

		ctx = context.WithValue(ctx, runnerCtxKey{}, runner)
		return unaryFunc(ctx, request)
	}
}))

func getMethodName(req connect.AnyRequest) string {
	splits := strings.Split(req.Spec().Procedure, "/")
	if len(splits) > 0 {
		return splits[len(splits)-1]
	}
	return ""
}

type runnerCtxKey struct{}

func GetRunner(ctx context.Context) *actions_model.ActionRunner {
	if v := ctx.Value(runnerCtxKey{}); v != nil {
		if r, ok := v.(*actions_model.ActionRunner); ok {
			return r
		}
	}
	return nil
}
