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

//go:build windows

package main

// Golang has the ability to load OS's timezone data from most UNIX systems (https://github.com/golang/go/blob/master/src/time/zoneinfo_unix.go)
// Even if the timezone data is missing, users could install the related packages to get it.
// But on Windows, although `zoneinfo_windows.go` tries to load the timezone data from Windows registry,
// some users still suffer from the issue that the timezone data is missing,
// So we import the tzdata package to make sure the timezone data is included in the binary.
//
// For non-Windows package builders, they could still use the "TAGS=timetzdata" to include the tzdata package in the binary.
// If we decided to add the tzdata for other platforms, modify the "go:build" directive above.
import _ "time/tzdata"
