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

//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/kumose/kmup/modules/assetfs"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: ./generate-bindata {local-directory} {bindata-filename}")
		os.Exit(1)
	}

	dir, filename := os.Args[1], os.Args[2]
	fmt.Printf("generating bindata for %s to %s\n", dir, filename)
	if err := assetfs.GenerateEmbedBindata(dir, filename); err != nil {
		fmt.Printf("failed: %s\n", err.Error())
		os.Exit(1)
	}
}
