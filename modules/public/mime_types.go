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

package public

import "strings"

// wellKnownMimeTypesLower comes from Golang's builtin mime package: `builtinTypesLower`, see the comment of detectWellKnownMimeType
var wellKnownMimeTypesLower = map[string]string{
	".avif": "image/avif",
	".css":  "text/css; charset=utf-8",
	".gif":  "image/gif",
	".htm":  "text/html; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".js":   "text/javascript; charset=utf-8",
	".json": "application/json",
	".mjs":  "text/javascript; charset=utf-8",
	".pdf":  "application/pdf",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".wasm": "application/wasm",
	".webp": "image/webp",
	".xml":  "text/xml; charset=utf-8",

	// well, there are some types missing from the builtin list
	".txt": "text/plain; charset=utf-8",
}

// detectWellKnownMimeType will return the mime-type for a well-known file ext name
// The purpose of this function is to bypass the unstable behavior of Golang's mime.TypeByExtension
// mime.TypeByExtension would use OS's mime-type config to overwrite the well-known types (see its document).
// If the user's OS has incorrect mime-type config, it would make Kmup can not respond a correct Content-Type to browsers.
// For example, if Kmup returns `text/plain` for a `.js` file, the browser couldn't run the JS due to security reasons.
// detectWellKnownMimeType makes the Content-Type for well-known files stable.
func detectWellKnownMimeType(ext string) string {
	ext = strings.ToLower(ext)
	return wellKnownMimeTypesLower[ext]
}
