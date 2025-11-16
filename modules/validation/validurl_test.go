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

package validation

import (
	"testing"

	"github.com/kumose-go/chi/binding"
)

func Test_ValidURLValidation(t *testing.T) {
	AddBindingRules()

	urlValidationTestCases := []validationTestCase{
		{
			description: "Empty URL",
			data: TestForm{
				URL: "",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "URL without port",
			data: TestForm{
				URL: "http://test.lan/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "URL with port",
			data: TestForm{
				URL: "http://test.lan:3326/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "URL with IPv6 address without port",
			data: TestForm{
				URL: "http://[::1]/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "URL with IPv6 address with port",
			data: TestForm{
				URL: "http://[::1]:3326/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "Invalid URL",
			data: TestForm{
				URL: "http//test.lan/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URL"},
					Classification: binding.ERR_URL,
					Message:        "Url",
				},
			},
		},
		{
			description: "Invalid schema",
			data: TestForm{
				URL: "ftp://test.lan/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URL"},
					Classification: binding.ERR_URL,
					Message:        "Url",
				},
			},
		},
		{
			description: "Invalid port",
			data: TestForm{
				URL: "http://test.lan:3x4/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URL"},
					Classification: binding.ERR_URL,
					Message:        "Url",
				},
			},
		},
		{
			description: "Invalid port with IPv6 address",
			data: TestForm{
				URL: "http://[::1]:3x4/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URL"},
					Classification: binding.ERR_URL,
					Message:        "Url",
				},
			},
		},
	}

	for _, testCase := range urlValidationTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			performValidationTest(t, testCase)
		})
	}
}
