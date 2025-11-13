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

package lfstransfer

import (
	"context"
	"fmt"
	"os"

	"github.com/kumose/kmup/modules/lfstransfer/backend"

	"github.com/charmbracelet/git-lfs-transfer/transfer"
)

func Main(ctx context.Context, repo, verb, token string) error {
	logger := newLogger()
	pktline := transfer.NewPktline(os.Stdin, os.Stdout, logger)
	kmupBackend, err := backend.New(ctx, repo, verb, token, logger)
	if err != nil {
		return err
	}

	for _, cap := range backend.Capabilities {
		if err := pktline.WritePacketText(cap); err != nil {
			logger.Log("error sending capability due to error:", err)
		}
	}
	if err := pktline.WriteFlush(); err != nil {
		logger.Log("error flushing capabilities:", err)
	}
	p := transfer.NewProcessor(pktline, kmupBackend, logger)
	defer logger.Log("done processing commands")
	switch verb {
	case "upload":
		return p.ProcessCommands(transfer.UploadOperation)
	case "download":
		return p.ProcessCommands(transfer.DownloadOperation)
	default:
		return fmt.Errorf("unknown operation %q", verb)
	}
}
