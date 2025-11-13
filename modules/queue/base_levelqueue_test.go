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
	"testing"

	"github.com/kumose/kmup/modules/queue/lqinternal"
	"github.com/kumose/kmup/modules/setting"

	"github.com/kumose-go/levelqueue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestBaseLevelDB(t *testing.T) {
	_, err := newBaseLevelQueueGeneric(&BaseConfig{ConnStr: "redis://"}, false)
	assert.ErrorContains(t, err, "invalid leveldb connection string")

	_, err = newBaseLevelQueueGeneric(&BaseConfig{DataFullDir: "relative"}, false)
	assert.ErrorContains(t, err, "invalid leveldb data dir")

	testQueueBasic(t, newBaseLevelQueueSimple, toBaseConfig("baseLevelQueue", setting.QueueSettings{Datadir: t.TempDir() + "/queue-test", Length: 10}), false)
	testQueueBasic(t, newBaseLevelQueueUnique, toBaseConfig("baseLevelQueueUnique", setting.QueueSettings{ConnStr: "leveldb://" + t.TempDir() + "/queue-test", Length: 10}), true)
}

func TestCorruptedLevelQueue(t *testing.T) {
	// sometimes the levelqueue could be in a corrupted state, this test is to make sure it can recover from it
	dbDir := t.TempDir() + "/levelqueue-test"
	db, err := leveldb.OpenFile(dbDir, nil)
	require.NoError(t, err)
	defer db.Close()

	assert.NoError(t, db.Put([]byte("other-key"), []byte("other-value"), nil))

	nameQueuePrefix := []byte("queue_name")
	nameSetPrefix := []byte("set_name")
	lq, err := levelqueue.NewUniqueQueue(db, nameQueuePrefix, nameSetPrefix, false)
	assert.NoError(t, err)
	assert.NoError(t, lq.RPush([]byte("item-1")))

	itemKey := lqinternal.QueueItemKeyBytes(nameQueuePrefix, 1)
	itemValue, err := db.Get(itemKey, nil)
	assert.NoError(t, err)
	assert.Equal(t, []byte("item-1"), itemValue)

	// there should be 5 keys in db: queue low, queue high, 1 queue item, 1 set item, and "other-key"
	keys := lqinternal.ListLevelQueueKeys(db)
	assert.Len(t, keys, 5)

	// delete the queue item key, to corrupt the queue
	assert.NoError(t, db.Delete(itemKey, nil))
	// now the queue is corrupted, it never works again
	_, err = lq.LPop()
	assert.ErrorIs(t, err, levelqueue.ErrNotFound)
	assert.NoError(t, lq.Close())

	// remove all the queue related keys to reset the queue
	lqinternal.RemoveLevelQueueKeys(db, nameQueuePrefix)
	lqinternal.RemoveLevelQueueKeys(db, nameSetPrefix)
	// now there should be only 1 key in db: "other-key"
	keys = lqinternal.ListLevelQueueKeys(db)
	assert.Len(t, keys, 1)
	assert.Equal(t, []byte("other-key"), keys[0])

	// re-create a queue from db
	lq, err = levelqueue.NewUniqueQueue(db, nameQueuePrefix, nameSetPrefix, false)
	assert.NoError(t, err)
	assert.NoError(t, lq.RPush([]byte("item-new-1")))
	// now the queue works again
	itemValue, err = lq.LPop()
	assert.NoError(t, err)
	assert.Equal(t, []byte("item-new-1"), itemValue)
	assert.NoError(t, lq.Close())
}
