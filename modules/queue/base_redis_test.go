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
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/kumose/kmup/modules/nosql"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func waitRedisReady(conn string, dur time.Duration) (ready bool) {
	ctxTimed, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for t := time.Now(); ; time.Sleep(50 * time.Millisecond) {
		ret := nosql.GetManager().GetRedisClient(conn).Ping(ctxTimed)
		if ret.Err() == nil {
			return true
		}
		if time.Since(t) > dur {
			return false
		}
	}
}

func redisServerCmd(t *testing.T) *exec.Cmd {
	redisServerProg, err := exec.LookPath("redis-server")
	if err != nil {
		return nil
	}
	c := &exec.Cmd{
		Path:   redisServerProg,
		Args:   []string{redisServerProg, "--bind", "127.0.0.1", "--port", "6379"},
		Dir:    t.TempDir(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return c
}

func TestBaseRedis(t *testing.T) {
	var redisServer *exec.Cmd
	defer func() {
		if redisServer != nil {
			_ = redisServer.Process.Signal(os.Interrupt)
			_ = redisServer.Wait()
		}
	}()
	if !waitRedisReady("redis://127.0.0.1:6379/0", 0) {
		redisServer = redisServerCmd(t)
		if redisServer == nil && os.Getenv("CI") == "" {
			t.Skip("redis-server not found")
			return
		}
		assert.NoError(t, redisServer.Start())
		require.True(t, waitRedisReady("redis://127.0.0.1:6379/0", 5*time.Second), "start redis-server")
	}

	testQueueBasic(t, newBaseRedisSimple, toBaseConfig("baseRedis", setting.QueueSettings{Length: 10}), false)
	testQueueBasic(t, newBaseRedisUnique, toBaseConfig("baseRedisUnique", setting.QueueSettings{Length: 10}), true)
}
