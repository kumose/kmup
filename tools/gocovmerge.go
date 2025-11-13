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

// gocovmerge takes the results from multiple `go test -coverprofile` runs and
// merges them into one profile

//go:build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"golang.org/x/tools/cover"
)

func mergeProfiles(p, merge *cover.Profile) {
	if p.Mode != merge.Mode {
		log.Fatalf("cannot merge profiles with different modes")
	}
	// Since the blocks are sorted, we can keep track of where the last block
	// was inserted and only look at the blocks after that as targets for merge
	startIndex := 0
	for _, b := range merge.Blocks {
		startIndex = mergeProfileBlock(p, b, startIndex)
	}
}

func mergeProfileBlock(p *cover.Profile, pb cover.ProfileBlock, startIndex int) int {
	sortFunc := func(i int) bool {
		pi := p.Blocks[i+startIndex]
		return pi.StartLine >= pb.StartLine && (pi.StartLine != pb.StartLine || pi.StartCol >= pb.StartCol)
	}

	i := 0
	if sortFunc(i) != true {
		i = sort.Search(len(p.Blocks)-startIndex, sortFunc)
	}
	i += startIndex
	if i < len(p.Blocks) && p.Blocks[i].StartLine == pb.StartLine && p.Blocks[i].StartCol == pb.StartCol {
		if p.Blocks[i].EndLine != pb.EndLine || p.Blocks[i].EndCol != pb.EndCol {
			log.Fatalf("OVERLAP MERGE: %v %v %v", p.FileName, p.Blocks[i], pb)
		}
		switch p.Mode {
		case "set":
			p.Blocks[i].Count |= pb.Count
		case "count", "atomic":
			p.Blocks[i].Count += pb.Count
		default:
			log.Fatalf("unsupported covermode: '%s'", p.Mode)
		}
	} else {
		if i > 0 {
			pa := p.Blocks[i-1]
			if pa.EndLine >= pb.EndLine && (pa.EndLine != pb.EndLine || pa.EndCol > pb.EndCol) {
				log.Fatalf("OVERLAP BEFORE: %v %v %v", p.FileName, pa, pb)
			}
		}
		if i < len(p.Blocks)-1 {
			pa := p.Blocks[i+1]
			if pa.StartLine <= pb.StartLine && (pa.StartLine != pb.StartLine || pa.StartCol < pb.StartCol) {
				log.Fatalf("OVERLAP AFTER: %v %v %v", p.FileName, pa, pb)
			}
		}
		p.Blocks = append(p.Blocks, cover.ProfileBlock{})
		copy(p.Blocks[i+1:], p.Blocks[i:])
		p.Blocks[i] = pb
	}
	return i + 1
}

func addProfile(profiles []*cover.Profile, p *cover.Profile) []*cover.Profile {
	i := sort.Search(len(profiles), func(i int) bool { return profiles[i].FileName >= p.FileName })
	if i < len(profiles) && profiles[i].FileName == p.FileName {
		mergeProfiles(profiles[i], p)
	} else {
		profiles = append(profiles, nil)
		copy(profiles[i+1:], profiles[i:])
		profiles[i] = p
	}
	return profiles
}

func dumpProfiles(profiles []*cover.Profile, out io.Writer) {
	if len(profiles) == 0 {
		return
	}
	fmt.Fprintf(out, "mode: %s\n", profiles[0].Mode)
	for _, p := range profiles {
		for _, b := range p.Blocks {
			fmt.Fprintf(out, "%s:%d.%d,%d.%d %d %d\n", p.FileName, b.StartLine, b.StartCol, b.EndLine, b.EndCol, b.NumStmt, b.Count)
		}
	}
}

func main() {
	flag.Parse()

	var merged []*cover.Profile

	for _, file := range flag.Args() {
		profiles, err := cover.ParseProfiles(file)
		if err != nil {
			log.Fatalf("failed to parse profile '%s': %v", file, err)
		}
		for _, p := range profiles {
			merged = addProfile(merged, p)
		}
	}

	dumpProfiles(merged, os.Stdout)
}
