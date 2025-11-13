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

package setting

type AttachmentSettingType struct {
	Storage      *Storage
	AllowedTypes string
	MaxSize      int64
	MaxFiles     int
	Enabled      bool
}

var Attachment AttachmentSettingType

func loadAttachmentFrom(rootCfg ConfigProvider) (err error) {
	Attachment = AttachmentSettingType{
		AllowedTypes: ".avif,.cpuprofile,.csv,.dmp,.docx,.fodg,.fodp,.fods,.fodt,.gif,.gz,.jpeg,.jpg,.json,.jsonc,.log,.md,.mov,.mp4,.odf,.odg,.odp,.ods,.odt,.patch,.pdf,.png,.pptx,.svg,.tgz,.txt,.webm,.webp,.xls,.xlsx,.zip",

		// FIXME: this size is used for both "issue attachment" and "release attachment"
		// The design is not right, these two should be different settings
		MaxSize: 2048,

		MaxFiles: 5,
		Enabled:  true,
	}
	sec, _ := rootCfg.GetSection("attachment")
	if sec == nil {
		Attachment.Storage, err = getStorage(rootCfg, "attachments", "", nil)
		return err
	}

	Attachment.AllowedTypes = sec.Key("ALLOWED_TYPES").MustString(Attachment.AllowedTypes)
	Attachment.MaxSize = sec.Key("MAX_SIZE").MustInt64(Attachment.MaxSize)
	Attachment.MaxFiles = sec.Key("MAX_FILES").MustInt(Attachment.MaxFiles)
	Attachment.Enabled = sec.Key("ENABLED").MustBool(Attachment.Enabled)
	Attachment.Storage, err = getStorage(rootCfg, "attachments", "", sec)
	return err
}
