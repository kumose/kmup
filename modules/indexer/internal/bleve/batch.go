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

package bleve

import (
	"github.com/blevesearch/bleve/v2"
)

// FlushingBatch is a batch of operations that automatically flushes to the
// underlying index once it reaches a certain size.
type FlushingBatch struct {
	maxBatchSize int
	batch        *bleve.Batch
	index        bleve.Index
}

// NewFlushingBatch creates a new flushing batch for the specified index. Once
// the number of operations in the batch reaches the specified limit, the batch
// automatically flushes its operations to the index.
func NewFlushingBatch(index bleve.Index, maxBatchSize int) *FlushingBatch {
	return &FlushingBatch{
		maxBatchSize: maxBatchSize,
		batch:        index.NewBatch(),
		index:        index,
	}
}

// Index add a new index to batch
func (b *FlushingBatch) Index(id string, data any) error {
	if err := b.batch.Index(id, data); err != nil {
		return err
	}
	return b.flushIfFull()
}

// Delete add a delete index to batch
func (b *FlushingBatch) Delete(id string) error {
	b.batch.Delete(id)
	return b.flushIfFull()
}

func (b *FlushingBatch) flushIfFull() error {
	if b.batch.Size() < b.maxBatchSize {
		return nil
	}
	return b.Flush()
}

// Flush submit the batch and create a new one
func (b *FlushingBatch) Flush() error {
	err := b.index.Batch(b.batch)
	if err != nil {
		return err
	}
	b.batch = b.index.NewBatch()
	return nil
}
