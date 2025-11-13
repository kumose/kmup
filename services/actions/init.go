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

package actions

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/queue"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
	notify_service "github.com/kumose/kmup/services/notify"
)

func initGlobalRunnerToken(ctx context.Context) error {
	// use the same env name as the runner, for consistency
	token := os.Getenv("KMUP_RUNNER_REGISTRATION_TOKEN")
	tokenFile := os.Getenv("KMUP_RUNNER_REGISTRATION_TOKEN_FILE")
	if token != "" && tokenFile != "" {
		return errors.New("both KMUP_RUNNER_REGISTRATION_TOKEN and KMUP_RUNNER_REGISTRATION_TOKEN_FILE are set, only one can be used")
	}
	if tokenFile != "" {
		file, err := os.ReadFile(tokenFile)
		if err != nil {
			return fmt.Errorf("unable to read KMUP_RUNNER_REGISTRATION_TOKEN_FILE: %w", err)
		}
		token = strings.TrimSpace(string(file))
	}
	if token == "" {
		return nil
	}

	if len(token) < 32 {
		return errors.New("KMUP_RUNNER_REGISTRATION_TOKEN must be at least 32 random characters")
	}

	existing, err := actions_model.GetRunnerToken(ctx, token)
	if err != nil && !errors.Is(err, util.ErrNotExist) {
		return fmt.Errorf("unable to check existing token: %w", err)
	}
	if existing != nil {
		if !existing.IsActive {
			log.Warn("The token defined by KMUP_RUNNER_REGISTRATION_TOKEN is already invalidated, please use the latest one from web UI")
		}
		return nil
	}
	_, err = actions_model.NewRunnerTokenWithValue(ctx, 0, 0, token)
	return err
}

func Init(ctx context.Context) error {
	if !setting.Actions.Enabled {
		return nil
	}

	jobEmitterQueue = queue.CreateUniqueQueue(graceful.GetManager().ShutdownContext(), "actions_ready_job", jobEmitterQueueHandler)
	if jobEmitterQueue == nil {
		return errors.New("unable to create actions_ready_job queue")
	}
	go graceful.GetManager().RunWithCancel(jobEmitterQueue)

	notify_service.RegisterNotifier(NewNotifier())
	return initGlobalRunnerToken(ctx)
}
