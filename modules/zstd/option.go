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

package zstd

import "github.com/klauspost/compress/zstd"

type WriterOption = zstd.EOption

var (
	WithEncoderCRC               = zstd.WithEncoderCRC
	WithEncoderConcurrency       = zstd.WithEncoderConcurrency
	WithWindowSize               = zstd.WithWindowSize
	WithEncoderPadding           = zstd.WithEncoderPadding
	WithEncoderLevel             = zstd.WithEncoderLevel
	WithZeroFrames               = zstd.WithZeroFrames
	WithAllLitEntropyCompression = zstd.WithAllLitEntropyCompression
	WithNoEntropyCompression     = zstd.WithNoEntropyCompression
	WithSingleSegment            = zstd.WithSingleSegment
	WithLowerEncoderMem          = zstd.WithLowerEncoderMem
	WithEncoderDict              = zstd.WithEncoderDict
	WithEncoderDictRaw           = zstd.WithEncoderDictRaw
)

type EncoderLevel = zstd.EncoderLevel

const (
	SpeedFastest           EncoderLevel = zstd.SpeedFastest
	SpeedDefault           EncoderLevel = zstd.SpeedDefault
	SpeedBetterCompression EncoderLevel = zstd.SpeedBetterCompression
	SpeedBestCompression   EncoderLevel = zstd.SpeedBestCompression
)

type ReaderOption = zstd.DOption

var (
	WithDecoderLowmem      = zstd.WithDecoderLowmem
	WithDecoderConcurrency = zstd.WithDecoderConcurrency
	WithDecoderMaxMemory   = zstd.WithDecoderMaxMemory
	WithDecoderDicts       = zstd.WithDecoderDicts
	WithDecoderDictRaw     = zstd.WithDecoderDictRaw
	WithDecoderMaxWindow   = zstd.WithDecoderMaxWindow
	WithDecodeAllCapLimit  = zstd.WithDecodeAllCapLimit
	WithDecodeBuffersBelow = zstd.WithDecodeBuffersBelow
	IgnoreChecksum         = zstd.IgnoreChecksum
)
