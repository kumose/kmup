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

package foreachref_test

import (
	"testing"

	"github.com/kumose/kmup/modules/git/foreachref"

	"github.com/stretchr/testify/require"
)

func TestFormat_Flag(t *testing.T) {
	tests := []struct {
		name string

		givenFormat foreachref.Format

		wantFlag string
	}{
		{
			name: "references are delimited by dual null chars",

			// no reference fields requested
			givenFormat: foreachref.NewFormat(),

			// only a reference delimiter field in --format
			wantFlag: "%00%00",
		},

		{
			name: "a field is a space-separated key-value pair",

			givenFormat: foreachref.NewFormat("refname:short"),

			// only a reference delimiter field
			wantFlag: "refname:short %(refname:short)%00%00",
		},

		{
			name: "fields are separated by a null char field-delimiter",

			givenFormat: foreachref.NewFormat("refname:short", "author"),

			wantFlag: "refname:short %(refname:short)%00author %(author)%00%00",
		},

		{
			name: "multiple fields",

			givenFormat: foreachref.NewFormat("refname:lstrip=2", "objecttype", "objectname"),

			wantFlag: "refname:lstrip=2 %(refname:lstrip=2)%00objecttype %(objecttype)%00objectname %(objectname)%00%00",
		},
	}

	for _, test := range tests {
		tc := test // don't close over loop variable
		t.Run(tc.name, func(t *testing.T) {
			gotFlag := tc.givenFormat.Flag()

			require.Equal(t, tc.wantFlag, gotFlag, "unexpected for-each-ref --format string. wanted: '%s', got: '%s'", tc.wantFlag, gotFlag)
		})
	}
}
