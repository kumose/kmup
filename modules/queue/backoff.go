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

package queue

import (
	"context"
	"time"
)

var (
	backoffBegin = 50 * time.Millisecond
	backoffUpper = 2 * time.Second
)

type (
	backoffFuncRetErr[T any] func() (retry bool, ret T, err error)
	backoffFuncErr           func() (retry bool, err error)
)

func mockBackoffDuration(d time.Duration) func() {
	oldBegin, oldUpper := backoffBegin, backoffUpper
	backoffBegin, backoffUpper = d, d
	return func() {
		backoffBegin, backoffUpper = oldBegin, oldUpper
	}
}

func backoffRetErr[T any](ctx context.Context, begin, upper time.Duration, end <-chan time.Time, fn backoffFuncRetErr[T]) (ret T, err error) {
	d := begin
	for {
		// check whether the context has been cancelled or has reached the deadline, return early
		select {
		case <-ctx.Done():
			return ret, ctx.Err()
		case <-end:
			return ret, context.DeadlineExceeded
		default:
		}

		// call the target function
		retry, ret, err := fn()
		if err != nil {
			return ret, err
		}
		if !retry {
			return ret, nil
		}

		// wait for a while before retrying, and also respect the context & deadline
		select {
		case <-ctx.Done():
			return ret, ctx.Err()
		case <-time.After(d):
			d *= 2
			if d > upper {
				d = upper
			}
		case <-end:
			return ret, context.DeadlineExceeded
		}
	}
}

func backoffErr(ctx context.Context, begin, upper time.Duration, end <-chan time.Time, fn backoffFuncErr) error {
	_, err := backoffRetErr(ctx, begin, upper, end, func() (retry bool, ret any, err error) {
		retry, err = fn()
		return retry, nil, err
	})
	return err
}
