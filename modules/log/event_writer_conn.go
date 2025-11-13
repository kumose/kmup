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
	"net"
)

type WriterConnOption struct {
	Addr           string
	Protocol       string
	Reconnect      bool
	ReconnectOnMsg bool
}

type eventWriterConn struct {
	*EventWriterBaseImpl
	connWriter connWriter
}

var _ EventWriter = (*eventWriterConn)(nil)

func NewEventWriterConn(writerName string, writerMode WriterMode) EventWriter {
	w := &eventWriterConn{EventWriterBaseImpl: NewEventWriterBase(writerName, "conn", writerMode)}
	opt := writerMode.WriterOption.(WriterConnOption)
	w.connWriter = connWriter{
		ReconnectOnMsg: opt.ReconnectOnMsg,
		Reconnect:      opt.Reconnect,
		Net:            opt.Protocol,
		Addr:           opt.Addr,
	}
	w.OutputWriteCloser = &w.connWriter
	return w
}

func init() {
	RegisterEventWriter("conn", NewEventWriterConn)
}

// below is copied from old code

type connWriter struct {
	innerWriter io.WriteCloser

	ReconnectOnMsg bool
	Reconnect      bool
	Net            string `json:"net"`
	Addr           string `json:"addr"`
}

var _ io.WriteCloser = (*connWriter)(nil)

// Close the inner writer
func (i *connWriter) Close() error {
	if i.innerWriter != nil {
		return i.innerWriter.Close()
	}
	return nil
}

// Write the data to the connection
func (i *connWriter) Write(p []byte) (int, error) {
	if i.neededConnectOnMsg() {
		if err := i.connect(); err != nil {
			return 0, err
		}
	}

	if i.ReconnectOnMsg {
		defer i.innerWriter.Close()
	}

	return i.innerWriter.Write(p)
}

func (i *connWriter) neededConnectOnMsg() bool {
	if i.Reconnect {
		i.Reconnect = false
		return true
	}

	if i.innerWriter == nil {
		return true
	}

	return i.ReconnectOnMsg
}

func (i *connWriter) connect() error {
	if i.innerWriter != nil {
		_ = i.innerWriter.Close()
		i.innerWriter = nil
	}

	conn, err := net.Dial(i.Net, i.Addr)
	if err != nil {
		return err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		err = tcpConn.SetKeepAlive(true)
		if err != nil {
			return err
		}
	}

	i.innerWriter = conn
	return nil
}
