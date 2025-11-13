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

package pprof

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/kumose/kmup/modules/log"
)

// DumpMemProfileForUsername dumps a memory profile at pprofDataPath as memprofile_<username>_<temporary id>
func DumpMemProfileForUsername(pprofDataPath, username string) error {
	f, err := os.CreateTemp(pprofDataPath, fmt.Sprintf("memprofile_%s_", username))
	if err != nil {
		return err
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	return pprof.WriteHeapProfile(f)
}

// DumpCPUProfileForUsername dumps a CPU profile at pprofDataPath as cpuprofile_<username>_<temporary id>
// the stop function it returns stops, writes and closes the CPU profile file
func DumpCPUProfileForUsername(pprofDataPath, username string) (func(), error) {
	f, err := os.CreateTemp(pprofDataPath, fmt.Sprintf("cpuprofile_%s_", username))
	if err != nil {
		return nil, err
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal("StartCPUProfile: %v", err)
	}
	return func() {
		pprof.StopCPUProfile()
		err = f.Close()
		if err != nil {
			log.Fatal("StopCPUProfile Close: %v", err)
		}
	}, nil
}
