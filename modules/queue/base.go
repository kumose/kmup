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

var pushBlockTime = 5 * time.Second

type baseQueue interface {
	PushItem(ctx context.Context, data []byte) error
	PopItem(ctx context.Context) ([]byte, error)
	HasItem(ctx context.Context, data []byte) (bool, error)
	Len(ctx context.Context) (int, error)
	Close() error
	RemoveAll(ctx context.Context) error
}

func popItemByChan(ctx context.Context, popItemFn func(ctx context.Context) ([]byte, error)) (chanItem chan []byte, chanErr chan error) {
	chanItem = make(chan []byte)
	chanErr = make(chan error)
	go func() {
		for {
			it, err := popItemFn(ctx)
			if err != nil {
				close(chanItem)
				chanErr <- err
				return
			}
			if it == nil {
				close(chanItem)
				close(chanErr)
				return
			}
			chanItem <- it
		}
	}()
	return chanItem, chanErr
}
