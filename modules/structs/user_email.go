// Copyright 2015 The Gogs Authors. All rights reserved.
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

package structs

// Email an email address belonging to a user
type Email struct {
	// swagger:strfmt email
	// The email address
	Email string `json:"email"`
	// Whether the email address has been verified
	Verified bool `json:"verified"`
	// Whether this is the primary email address
	Primary bool `json:"primary"`
	// The unique identifier of the user who owns this email
	UserID int64 `json:"user_id"`
	// username of the user
	UserName string `json:"username"`
}

// CreateEmailOption options when creating email addresses
type CreateEmailOption struct {
	// email addresses to add
	Emails []string `json:"emails"`
}

// DeleteEmailOption options when deleting email addresses
type DeleteEmailOption struct {
	// email addresses to delete
	Emails []string `json:"emails"`
}
