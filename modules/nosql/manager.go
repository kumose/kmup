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

package nosql

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/kumose/kmup/modules/process"

	"github.com/redis/go-redis/v9"
	"github.com/syndtr/goleveldb/leveldb"
)

var manager *Manager

// Manager is the nosql connection manager
type Manager struct {
	ctx      context.Context
	finished context.CancelFunc
	mutex    sync.Mutex

	RedisConnections   map[string]*redisClientHolder
	LevelDBConnections map[string]*levelDBHolder
}

type redisClientHolder struct {
	redis.UniversalClient
	name  []string
	count int64
}

func (r *redisClientHolder) Close() error {
	return manager.CloseRedisClient(r.name[0])
}

type levelDBHolder struct {
	name  []string
	count int64
	db    *leveldb.DB
}

func init() {
	_ = GetManager()
}

// GetManager returns a Manager and initializes one as singleton is there's none yet
func GetManager() *Manager {
	if manager == nil {
		ctx, _, finished := process.GetManager().AddTypedContext(context.Background(), "Service: NoSQL", process.SystemProcessType, false)
		manager = &Manager{
			ctx:                ctx,
			finished:           finished,
			RedisConnections:   make(map[string]*redisClientHolder),
			LevelDBConnections: make(map[string]*levelDBHolder),
		}
	}
	return manager
}

func valToTimeDuration(vs []string) (result time.Duration) {
	var err error
	for _, v := range vs {
		result, err = time.ParseDuration(v)
		if err != nil {
			var val int
			val, err = strconv.Atoi(v)
			result = time.Duration(val)
		}
		if err == nil {
			return result
		}
	}
	return result
}
