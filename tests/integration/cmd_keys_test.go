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

package integration

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/kumose/kmup/cmd"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func Test_CmdKeys(t *testing.T) {
	onKmupRun(t, func(*testing.T, *url.URL) {
		tests := []struct {
			name           string
			args           []string
			wantErr        bool
			expectedOutput string
		}{
			{"test_empty_1", []string{"keys", "--username=git", "--type=test", "--content=test"}, true, ""},
			{"test_empty_2", []string{"keys", "-e", "git", "-u", "git", "-t", "test", "-k", "test"}, true, ""},
			{
				"with_key",
				[]string{"keys", "-e", "git", "-u", "git", "-t", "ssh-rsa", "-k", "AAAAB3NzaC1yc2EAAAADAQABAAABgQDWVj0fQ5N8wNc0LVNA41wDLYJ89ZIbejrPfg/avyj3u/ZohAKsQclxG4Ju0VirduBFF9EOiuxoiFBRr3xRpqzpsZtnMPkWVWb+akZwBFAx8p+jKdy4QXR/SZqbVobrGwip2UjSrri1CtBxpJikojRIZfCnDaMOyd9Jp6KkujvniFzUWdLmCPxUE9zhTaPu0JsEP7MW0m6yx7ZUhHyfss+NtqmFTaDO+QlMR7L2QkDliN2Jl3Xa3PhuWnKJfWhdAq1Cw4oraKUOmIgXLkuiuxVQ6mD3AiFupkmfqdHq6h+uHHmyQqv3gU+/sD8GbGAhf6ftqhTsXjnv1Aj4R8NoDf9BS6KRkzkeun5UisSzgtfQzjOMEiJtmrep2ZQrMGahrXa+q4VKr0aKJfm+KlLfwm/JztfsBcqQWNcTURiCFqz+fgZw0Ey/de0eyMzldYTdXXNRYCKjs9bvBK+6SSXRM7AhftfQ0ZuoW5+gtinPrnmoOaSCEJbAiEiTO/BzOHgowiM="},
				false,
				"# kmup public key\ncommand=\"" + setting.AppPath + " --config=" + util.ShellEscape(setting.CustomConf) + " serv key-1\",no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty,no-user-rc,restrict ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDWVj0fQ5N8wNc0LVNA41wDLYJ89ZIbejrPfg/avyj3u/ZohAKsQclxG4Ju0VirduBFF9EOiuxoiFBRr3xRpqzpsZtnMPkWVWb+akZwBFAx8p+jKdy4QXR/SZqbVobrGwip2UjSrri1CtBxpJikojRIZfCnDaMOyd9Jp6KkujvniFzUWdLmCPxUE9zhTaPu0JsEP7MW0m6yx7ZUhHyfss+NtqmFTaDO+QlMR7L2QkDliN2Jl3Xa3PhuWnKJfWhdAq1Cw4oraKUOmIgXLkuiuxVQ6mD3AiFupkmfqdHq6h+uHHmyQqv3gU+/sD8GbGAhf6ftqhTsXjnv1Aj4R8NoDf9BS6KRkzkeun5UisSzgtfQzjOMEiJtmrep2ZQrMGahrXa+q4VKr0aKJfm+KlLfwm/JztfsBcqQWNcTURiCFqz+fgZw0Ey/de0eyMzldYTdXXNRYCKjs9bvBK+6SSXRM7AhftfQ0ZuoW5+gtinPrnmoOaSCEJbAiEiTO/BzOHgowiM= user-2\n",
			},
			{"invalid", []string{"keys", "--not-a-flag=git"}, true, "Incorrect Usage: flag provided but not defined: -not-a-flag\n\n"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// FIXME: this test is not quite right. Each "command run" always re-initializes settings
				defer test.MockVariableValue(&cmd.CmdKeys.Before, nil)() // don't re-initialize logger during the test

				var stdout, stderr bytes.Buffer
				app := &cli.Command{
					Writer:    &stdout,
					ErrWriter: &stderr,
					Commands:  []*cli.Command{cmd.CmdKeys},
				}
				err := app.Run(t.Context(), append([]string{"prog"}, tt.args...))
				if tt.wantErr {
					assert.Error(t, err)
					assert.Equal(t, tt.expectedOutput, stderr.String())
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedOutput, stdout.String())
				}
			})
		}
	})
}
