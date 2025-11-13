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

package updatechecker

import (
	"context"
	"io"
	"net/http"

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/proxy"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/system"

	"github.com/hashicorp/go-version"
)

// CheckerState stores the remote version from the JSON endpoint
type CheckerState struct {
	LatestVersion string
}

// Name returns the name of the state item for update checker
func (r *CheckerState) Name() string {
	return "update-checker"
}

// KmupUpdateChecker returns error when new version of Kmup is available
func KmupUpdateChecker(httpEndpoint string) error {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: proxy.Proxy(),
		},
	}

	req, err := http.NewRequest(http.MethodGet, httpEndpoint, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type respType struct {
		Latest struct {
			Version string `json:"version"`
		} `json:"latest"`
	}
	respData := respType{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}

	return UpdateRemoteVersion(req.Context(), respData.Latest.Version)
}

// UpdateRemoteVersion updates the latest available version of Kmup
func UpdateRemoteVersion(ctx context.Context, version string) (err error) {
	return system.AppState.Set(ctx, &CheckerState{LatestVersion: version})
}

// GetRemoteVersion returns the current remote version (or currently installed version if fail to fetch from DB)
func GetRemoteVersion(ctx context.Context) string {
	item := new(CheckerState)
	if err := system.AppState.Get(ctx, item); err != nil {
		return ""
	}
	return item.LatestVersion
}

// GetNeedUpdate returns true whether a newer version of Kmup is available
func GetNeedUpdate(ctx context.Context) bool {
	curVer, err := version.NewVersion(setting.AppVer)
	if err != nil {
		// return false to fail silently
		return false
	}
	remoteVerStr := GetRemoteVersion(ctx)
	if remoteVerStr == "" {
		// no remote version is known
		return false
	}
	remoteVer, err := version.NewVersion(remoteVerStr)
	if err != nil {
		// return false to fail silently
		return false
	}
	return curVer.LessThan(remoteVer)
}
