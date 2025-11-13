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

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/doctor"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestDoctorRun(t *testing.T) {
	doctor.Register(&doctor.Check{
		Title: "Test Check",
		Name:  "test-check",
		Run:   func(ctx context.Context, logger log.Logger, autofix bool) error { return nil },

		SkipDatabaseInitialization: true,
	})
	app := &cli.Command{
		Commands: []*cli.Command{cmdDoctorCheck},
	}
	err := app.Run(t.Context(), []string{"./kmup", "check", "--run", "test-check"})
	assert.NoError(t, err)
	err = app.Run(t.Context(), []string{"./kmup", "check", "--run", "no-such"})
	assert.ErrorContains(t, err, `unknown checks: "no-such"`)
	err = app.Run(t.Context(), []string{"./kmup", "check", "--run", "test-check,no-such"})
	assert.ErrorContains(t, err, `unknown checks: "no-such"`)
}
