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
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/kumose/kmup/modules/container"
)

// regexp is based on go-license, excluding README and NOTICE
// https://github.com/google/go-licenses/blob/master/licenses/find.go
var licenseRe = regexp.MustCompile(`^(?i)((UN)?LICEN(S|C)E|COPYING).*$`)

type LicenseEntry struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	LicenseText string `json:"licenseText"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run generate-go-licenses.go <base-dir> <out-json-file>")
		os.Exit(1)
	}

	base, out := os.Args[1], os.Args[2]

	// Add ext for excluded files because license_test.go will be included for some reason.
	// And there are more files that should be excluded, check with:
	//
	// go run github.com/google/go-licenses@v1.6.0 save . --force --save_path=.go-licenses 2>/dev/null
	// find .go-licenses -type f | while read FILE; do echo "${$(basename $FILE)##*.}"; done | sort -u
	//    AUTHORS
	//    COPYING
	//    LICENSE
	//    Makefile
	//    NOTICE
	//    gitignore
	//    go
	//    md
	//    mod
	//    sum
	//    toml
	//    txt
	//    yml
	//
	// It could be removed once we have a better regex.
	excludedExt := container.SetOf(".gitignore", ".go", ".mod", ".sum", ".toml", ".yml")

	var paths []string
	err := filepath.WalkDir(base, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !licenseRe.MatchString(entry.Name()) || excludedExt.Contains(filepath.Ext(entry.Name())) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	sort.Strings(paths)

	var entries []LicenseEntry
	for _, filePath := range paths {
		licenseText, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		pkgPath := filepath.ToSlash(filePath)
		pkgPath = strings.TrimPrefix(pkgPath, base+"/")
		pkgName := path.Dir(pkgPath)

		// There might be a bug somewhere in go-licenses that sometimes interprets the
		// root package as "." and sometimes as "github.com/kumose/kmup". Workaround by
		// removing both of them for the sake of stable output.
		if pkgName == "." || pkgName == "github.com/kumose/kmup" {
			continue
		}

		entries = append(entries, LicenseEntry{
			Name:        pkgName,
			Path:        pkgPath,
			LicenseText: string(licenseText),
		})
	}

	jsonBytes, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		panic(err)
	}

	// Ensure file has a final newline
	if jsonBytes[len(jsonBytes)-1] != '\n' {
		jsonBytes = append(jsonBytes, '\n')
	}

	err = os.WriteFile(out, jsonBytes, 0o644)
	if err != nil {
		panic(err)
	}
}
