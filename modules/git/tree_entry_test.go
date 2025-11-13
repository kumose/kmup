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

package git

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntriesCustomSort(t *testing.T) {
	entries := Entries{
		&TreeEntry{name: "a-dir", entryMode: EntryModeTree},
		&TreeEntry{name: "a-submodule", entryMode: EntryModeCommit},
		&TreeEntry{name: "b-dir", entryMode: EntryModeTree},
		&TreeEntry{name: "b-submodule", entryMode: EntryModeCommit},
		&TreeEntry{name: "a-file", entryMode: EntryModeBlob},
		&TreeEntry{name: "b-file", entryMode: EntryModeBlob},
	}
	expected := slices.Clone(entries)
	rand.Shuffle(len(entries), func(i, j int) { entries[i], entries[j] = entries[j], entries[i] })
	assert.NotEqual(t, expected, entries)
	entries.CustomSort(strings.Compare)
	assert.Equal(t, expected, entries)
}
