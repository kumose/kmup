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

package util

import (
	"bytes"
	"errors"
	"io"
)

type NopCloser struct {
	io.Writer
}

func (NopCloser) Close() error { return nil }

// ReadAtMost reads at most len(buf) bytes from r into buf.
// It returns the number of bytes copied. n is only less than len(buf) if r provides fewer bytes.
// If EOF or ErrUnexpectedEOF occurs while reading, err will be nil.
func ReadAtMost(r io.Reader, buf []byte) (n int, err error) {
	n, err = io.ReadFull(r, buf)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		err = nil
	}
	return n, err
}

// ReadWithLimit reads at most "limit" bytes from r into buf.
// If EOF or ErrUnexpectedEOF occurs while reading, err will be nil.
func ReadWithLimit(r io.Reader, n int) (buf []byte, err error) {
	return readWithLimit(r, 1024, n)
}

func readWithLimit(r io.Reader, batch, limit int) ([]byte, error) {
	if limit <= batch {
		buf := make([]byte, limit)
		n, err := ReadAtMost(r, buf)
		if err != nil {
			return nil, err
		}
		return buf[:n], nil
	}
	res := bytes.NewBuffer(make([]byte, 0, batch))
	bufFix := make([]byte, batch)
	eof := false
	for res.Len() < limit && !eof {
		bufTmp := bufFix
		if res.Len()+batch > limit {
			bufTmp = bufFix[:limit-res.Len()]
		}
		n, err := io.ReadFull(r, bufTmp)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			eof = true
		} else if err != nil {
			return nil, err
		}
		if _, err = res.Write(bufTmp[:n]); err != nil {
			return nil, err
		}
	}
	return res.Bytes(), nil
}

// ErrNotEmpty is an error reported when there is a non-empty reader
var ErrNotEmpty = errors.New("not-empty")

// IsEmptyReader reads a reader and ensures it is empty
func IsEmptyReader(r io.Reader) (err error) {
	var buf [1]byte

	for {
		n, err := r.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if n > 0 {
			return ErrNotEmpty
		}
	}
}

type CountingReader struct {
	io.Reader
	n int
}

var _ io.Reader = &CountingReader{}

func (w *CountingReader) Count() int {
	return w.n
}

func (w *CountingReader) Read(p []byte) (int, error) {
	n, err := w.Reader.Read(p)
	w.n += n
	return n, err
}

func NewCountingReader(rd io.Reader) *CountingReader {
	return &CountingReader{Reader: rd}
}
