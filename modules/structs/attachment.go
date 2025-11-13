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

package structs // import "github.com/kumose/kmup/modules/structs"

import (
	"time"
)

// Attachment a generic attachment
// swagger:model
type Attachment struct {
	// ID is the unique identifier for the attachment
	ID int64 `json:"id"`
	// Name is the filename of the attachment
	Name string `json:"name"`
	// Size is the file size in bytes
	Size int64 `json:"size"`
	// DownloadCount is the number of times the attachment has been downloaded
	DownloadCount int64 `json:"download_count"`
	// swagger:strfmt date-time
	// Created is the time when the attachment was uploaded
	Created time.Time `json:"created_at"`
	// UUID is the unique identifier for the attachment file
	UUID string `json:"uuid"`
	// DownloadURL is the URL to download the attachment
	DownloadURL string `json:"browser_download_url"`
}

// EditAttachmentOptions options for editing attachments
// swagger:model
type EditAttachmentOptions struct {
	// Name is the new filename for the attachment
	Name string `json:"name"`
}
