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

package admin

import (
	"archive/zip"
	"fmt"
	"runtime/pprof"
	"time"

	"github.com/kumose/kmup/modules/httplib"
	"github.com/kumose/kmup/modules/tailmsg"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
)

func MonitorDiagnosis(ctx *context.Context) {
	seconds := min(max(ctx.FormInt64("seconds"), 1), 300)

	httplib.ServeSetHeaders(ctx.Resp, &httplib.ServeHeaderOptions{
		ContentType: "application/zip",
		Disposition: "attachment",
		Filename:    fmt.Sprintf("kmup-diagnosis-%s.zip", time.Now().Format("20060102-150405")),
	})

	zipWriter := zip.NewWriter(ctx.Resp)
	defer zipWriter.Close()

	f, err := zipWriter.CreateHeader(&zip.FileHeader{Name: "goroutine-before.txt", Method: zip.Deflate, Modified: time.Now()})
	if err != nil {
		ctx.ServerError("Failed to create zip file", err)
		return
	}
	_ = pprof.Lookup("goroutine").WriteTo(f, 1)

	f, err = zipWriter.CreateHeader(&zip.FileHeader{Name: "cpu-profile.dat", Method: zip.Deflate, Modified: time.Now()})
	if err != nil {
		ctx.ServerError("Failed to create zip file", err)
		return
	}

	err = pprof.StartCPUProfile(f)
	if err == nil {
		time.Sleep(time.Duration(seconds) * time.Second)
		pprof.StopCPUProfile()
	} else {
		_, _ = f.Write([]byte(err.Error()))
	}

	f, err = zipWriter.CreateHeader(&zip.FileHeader{Name: "goroutine-after.txt", Method: zip.Deflate, Modified: time.Now()})
	if err != nil {
		ctx.ServerError("Failed to create zip file", err)
		return
	}
	_ = pprof.Lookup("goroutine").WriteTo(f, 1)

	f, err = zipWriter.CreateHeader(&zip.FileHeader{Name: "heap.dat", Method: zip.Deflate, Modified: time.Now()})
	if err != nil {
		ctx.ServerError("Failed to create zip file", err)
		return
	}
	_ = pprof.Lookup("heap").WriteTo(f, 0)

	f, err = zipWriter.CreateHeader(&zip.FileHeader{Name: "perftrace.txt", Method: zip.Deflate, Modified: time.Now()})
	if err != nil {
		ctx.ServerError("Failed to create zip file", err)
		return
	}
	for _, record := range tailmsg.GetManager().GetTraceRecorder().GetRecords() {
		_, _ = f.Write(util.UnsafeStringToBytes(record.Time.Format(time.RFC3339)))
		_, _ = f.Write([]byte(" "))
		_, _ = f.Write(util.UnsafeStringToBytes((record.Content)))
		_, _ = f.Write([]byte("\n\n"))
	}
}
