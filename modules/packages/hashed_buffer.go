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

package packages

import (
	"io"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util/filebuffer"
)

// HashedSizeReader provide methods to read, sum hashes and a Size method
type HashedSizeReader interface {
	io.Reader
	HashSummer
	Size() int64
}

// HashedBuffer is buffer which calculates multiple checksums
type HashedBuffer struct {
	*filebuffer.FileBackedBuffer

	hash *MultiHasher

	combinedWriter io.Writer
}

const DefaultMemorySize = 32 * 1024 * 1024

// NewHashedBuffer creates a hashed buffer with the default memory size
func NewHashedBuffer() (*HashedBuffer, error) {
	return NewHashedBufferWithSize(DefaultMemorySize)
}

// NewHashedBufferWithSize creates a hashed buffer with a specific memory size
func NewHashedBufferWithSize(maxMemorySize int) (*HashedBuffer, error) {
	tempDir, err := setting.AppDataTempDir("package-hashed-buffer").MkdirAllSub("")
	if err != nil {
		return nil, err
	}
	b := filebuffer.New(maxMemorySize, tempDir)
	hash := NewMultiHasher()

	combinedWriter := io.MultiWriter(b, hash)

	return &HashedBuffer{
		b,
		hash,
		combinedWriter,
	}, nil
}

// CreateHashedBufferFromReader creates a hashed buffer with the default memory size and copies the provided reader data into it.
func CreateHashedBufferFromReader(r io.Reader) (*HashedBuffer, error) {
	return CreateHashedBufferFromReaderWithSize(r, DefaultMemorySize)
}

// CreateHashedBufferFromReaderWithSize creates a hashed buffer and copies the provided reader data into it.
func CreateHashedBufferFromReaderWithSize(r io.Reader, maxMemorySize int) (*HashedBuffer, error) {
	b, err := NewHashedBufferWithSize(maxMemorySize)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(b, r)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Write implements io.Writer
func (b *HashedBuffer) Write(p []byte) (int, error) {
	return b.combinedWriter.Write(p)
}

// Sums gets the MD5, SHA1, SHA256 and SHA512 checksums of the data
func (b *HashedBuffer) Sums() (hashMD5, hashSHA1, hashSHA256, hashSHA512 []byte) {
	return b.hash.Sums()
}
