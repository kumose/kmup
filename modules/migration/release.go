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

package migration

import (
	"io"
	"time"
)

// ReleaseAsset represents a release asset
type ReleaseAsset struct {
	ID            int64
	Name          string
	ContentType   *string `yaml:"content_type"`
	Size          *int
	DownloadCount *int `yaml:"download_count"`
	Created       time.Time
	Updated       time.Time

	DownloadURL *string `yaml:"download_url"` // SECURITY: It is the responsibility of downloader to make sure this is safe
	// if DownloadURL is nil, the function should be invoked
	DownloadFunc func() (io.ReadCloser, error) `yaml:"-"` // SECURITY: It is the responsibility of downloader to make sure this is safe
}

// Release represents a release
type Release struct {
	TagName         string `yaml:"tag_name"`         // SECURITY: This must pass git.IsValidRefPattern
	TargetCommitish string `yaml:"target_commitish"` // SECURITY: This must pass git.IsValidRefPattern
	Name            string
	Body            string
	Draft           bool
	Prerelease      bool
	PublisherID     int64  `yaml:"publisher_id"`
	PublisherName   string `yaml:"publisher_name"`
	PublisherEmail  string `yaml:"publisher_email"`
	Assets          []*ReleaseAsset
	Created         time.Time
	Published       time.Time
}

// GetExternalName ExternalUserMigrated interface
func (r *Release) GetExternalName() string { return r.PublisherName }

// GetExternalID ExternalUserMigrated interface
func (r *Release) GetExternalID() int64 { return r.PublisherID }
