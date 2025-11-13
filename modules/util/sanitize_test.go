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

package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeErrorCredentialURLs(t *testing.T) {
	err := errors.New("error with https://a@b.com")
	se := SanitizeErrorCredentialURLs(err)
	assert.Equal(t, "error with https://"+userPlaceholder+"@b.com", se.Error())
}

func TestSanitizeCredentialURLs(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			"https://github.com/go-kmup/test_repo.git",
			"https://github.com/go-kmup/test_repo.git",
		},
		{
			"https://mytoken@github.com/go-kmup/test_repo.git",
			"https://" + userPlaceholder + "@github.com/go-kmup/test_repo.git",
		},
		{
			"https://user:password@github.com/go-kmup/test_repo.git",
			"https://" + userPlaceholder + "@github.com/go-kmup/test_repo.git",
		},
		{
			"ftp://x@",
			"ftp://" + userPlaceholder + "@",
		},
		{
			"ftp://x/@",
			"ftp://x/@",
		},
		{
			"ftp://u@x/@", // test multiple @ chars
			"ftp://" + userPlaceholder + "@x/@",
		},
		{
			"😊ftp://u@x😊", // test unicode
			"😊ftp://" + userPlaceholder + "@x😊",
		},
		{
			"://@",
			"://@",
		},
		{
			"//u:p@h", // do not process URLs without explicit scheme, they are not treated as "valid" URLs because there is no scheme context in string
			"//u:p@h",
		},
		{
			"s://u@h", // the minimal pattern to be sanitized
			"s://" + userPlaceholder + "@h",
		},
		{
			"URLs in log https://u:b@h and https://u:b@h:80/, with https://h.com and u@h.com",
			"URLs in log https://" + userPlaceholder + "@h and https://" + userPlaceholder + "@h:80/, with https://h.com and u@h.com",
		},
	}

	for n, c := range cases {
		result := SanitizeCredentialURLs(c.input)
		assert.Equal(t, c.expected, result, "case %d: error should match", n)
	}
}
