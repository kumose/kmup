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

package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestDefaultCommand(t *testing.T) {
	test := func(t *testing.T, args []string, expectedRetName string, expectedRetValid bool) {
		called := false
		cmd := &cli.Command{
			DefaultCommand: "test",
			Commands: []*cli.Command{
				{
					Name: "test",
					Action: func(ctx context.Context, command *cli.Command) error {
						retName, retValid := isValidDefaultSubCommand(command)
						assert.Equal(t, expectedRetName, retName)
						assert.Equal(t, expectedRetValid, retValid)
						called = true
						return nil
					},
				},
			},
		}
		assert.NoError(t, cmd.Run(t.Context(), args))
		assert.True(t, called)
	}
	test(t, []string{"./kmup"}, "", true)
	test(t, []string{"./kmup", "test"}, "", true)
	test(t, []string{"./kmup", "other"}, "other", false)
}
