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

package gitgrep

import (
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestIndexSettingToGitGrepPathspecList(t *testing.T) {
	defer test.MockVariableValue(&setting.Indexer.IncludePatterns, setting.IndexerGlobFromString("a"))()
	defer test.MockVariableValue(&setting.Indexer.ExcludePatterns, setting.IndexerGlobFromString("b"))()
	assert.Equal(t, []string{":(glob)a", ":(glob,exclude)b"}, indexSettingToGitGrepPathspecList())
}
