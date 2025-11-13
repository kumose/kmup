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

package log

import (
	"io"

	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/modules/util/rotatingfilewriter"
)

type WriterFileOption struct {
	FileName         string
	MaxSize          int64
	LogRotate        bool
	DailyRotate      bool
	MaxDays          int
	Compress         bool
	CompressionLevel int
}

type eventWriterFile struct {
	*EventWriterBaseImpl
	fileWriter io.WriteCloser
}

var _ EventWriter = (*eventWriterFile)(nil)

func NewEventWriterFile(name string, mode WriterMode) EventWriter {
	w := &eventWriterFile{EventWriterBaseImpl: NewEventWriterBase(name, "file", mode)}
	opt := mode.WriterOption.(WriterFileOption)
	var err error
	w.fileWriter, err = rotatingfilewriter.Open(opt.FileName, &rotatingfilewriter.Options{
		Rotate:           opt.LogRotate,
		MaximumSize:      opt.MaxSize,
		RotateDaily:      opt.DailyRotate,
		KeepDays:         opt.MaxDays,
		Compress:         opt.Compress,
		CompressionLevel: opt.CompressionLevel,
	})
	if err != nil {
		// if the log file can't be opened, what should it do? panic/exit? ignore logs? fallback to stderr?
		// it seems that "fallback to stderr" is slightly better than others ....
		FallbackErrorf("unable to open log file %q: %v", opt.FileName, err)
		w.fileWriter = util.NopCloser{Writer: LoggerToWriter(FallbackErrorf)}
	}
	w.OutputWriteCloser = w.fileWriter
	return w
}

func init() {
	RegisterEventWriter("file", NewEventWriterFile)
}
