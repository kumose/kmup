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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellEscape(t *testing.T) {
	tests := []struct {
		name     string
		toEscape string
		want     string
	}{
		{
			"Simplest case - nothing to escape",
			"a/b/c/d",
			"a/b/c/d",
		}, {
			"Prefixed tilde - with normal stuff - should not escape",
			"~/src/go/kumose/kmup",
			"~/src/go/kumose/kmup",
		}, {
			"Typical windows path with spaces - should get doublequote escaped",
			`C:\Program Files\Kmup v1.13 - I like lots of spaces\kmup`,
			`"C:\\Program Files\\Kmup v1.13 - I like lots of spaces\\kmup"`,
		}, {
			"Forward-slashed windows path with spaces - should get doublequote escaped",
			"C:/Program Files/Kmup v1.13 - I like lots of spaces/kmup",
			`"C:/Program Files/Kmup v1.13 - I like lots of spaces/kmup"`,
		}, {
			"Prefixed tilde - but then a space filled path",
			"~git/Kmup v1.13/kmup",
			`~git/"Kmup v1.13/kmup"`,
		}, {
			"Bangs are unfortunately not predictable so need to be singlequoted",
			"C:/Program Files/Kmup!/kmup",
			`'C:/Program Files/Kmup!/kmup'`,
		}, {
			"Newlines are just irritating",
			"/home/git/Kmup\n\nWHY-WOULD-YOU-DO-THIS\n\nKmup/kmup",
			"'/home/git/Kmup\n\nWHY-WOULD-YOU-DO-THIS\n\nKmup/kmup'",
		}, {
			"Similarly we should nicely handle multiple single quotes if we have to single-quote",
			"'!''!'''!''!'!'",
			`\''!'\'\''!'\'\'\''!'\'\''!'\''!'\'`,
		}, {
			"Double quote < ...",
			"~/<kmup",
			"~/\"<kmup\"",
		}, {
			"Double quote > ...",
			"~/kmup>",
			"~/\"kmup>\"",
		}, {
			"Double quote and escape $ ...",
			"~/$kmup",
			"~/\"\\$kmup\"",
		}, {
			"Double quote {...",
			"~/{kmup",
			"~/\"{kmup\"",
		}, {
			"Double quote }...",
			"~/kmup}",
			"~/\"kmup}\"",
		}, {
			"Double quote ()...",
			"~/(kmup)",
			"~/\"(kmup)\"",
		}, {
			"Double quote and escape `...",
			"~/kmup`",
			"~/\"kmup\\`\"",
		}, {
			"Double quotes can handle a number of things without having to escape them but not everything ...",
			"~/<kmup> ${kmup} `kmup` [kmup] (kmup) \"kmup\" \\kmup\\ 'kmup'",
			"~/\"<kmup> \\${kmup} \\`kmup\\` [kmup] (kmup) \\\"kmup\\\" \\\\kmup\\\\ 'kmup'\"",
		}, {
			"Single quotes don't need to escape except for '...",
			"~/<kmup> ${kmup} `kmup` (kmup) !kmup! \"kmup\" \\kmup\\ 'kmup'",
			"~/'<kmup> ${kmup} `kmup` (kmup) !kmup! \"kmup\" \\kmup\\ '\\''kmup'\\'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ShellEscape(tt.toEscape))
		})
	}
}
