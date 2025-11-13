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

func Test_ValidURLListValidation(t *testing.T) {
	AddBindingRules()

	// This is a copy of all the URL tests cases, plus additional ones to
	// account for multiple URLs
	urlListValidationTestCases := []validationTestCase{
		{
			description: "Empty URL",
			data: TestForm{
				URLs: "",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "URL without port",
			data: TestForm{
				URLs: "http://test.lan/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "URL with port",
			data: TestForm{
				URLs: "http://test.lan:3000/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "URL with IPv6 address without port",
			data: TestForm{
				URLs: "http://[::1]/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "URL with IPv6 address with port",
			data: TestForm{
				URLs: "http://[::1]:3000/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "Invalid URL",
			data: TestForm{
				URLs: "http//test.lan/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URLs"},
					Classification: binding.ERR_URL,
					Message:        "http//test.lan/",
				},
			},
		},
		{
			description: "Invalid schema",
			data: TestForm{
				URLs: "ftp://test.lan/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URLs"},
					Classification: binding.ERR_URL,
					Message:        "ftp://test.lan/",
				},
			},
		},
		{
			description: "Invalid port",
			data: TestForm{
				URLs: "http://test.lan:3x4/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URLs"},
					Classification: binding.ERR_URL,
					Message:        "http://test.lan:3x4/",
				},
			},
		},
		{
			description: "Invalid port with IPv6 address",
			data: TestForm{
				URLs: "http://[::1]:3x4/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URLs"},
					Classification: binding.ERR_URL,
					Message:        "http://[::1]:3x4/",
				},
			},
		},
		{
			description: "Multi URLs",
			data: TestForm{
				URLs: "http://test.lan:3000/\nhttp://test.local/",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "Multi URLs with newline",
			data: TestForm{
				URLs: "http://test.lan:3000/\nhttp://test.local/\n",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "List with invalid entry",
			data: TestForm{
				URLs: "http://test.lan:3000/\nhttp://[::1]:3x4/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URLs"},
					Classification: binding.ERR_URL,
					Message:        "http://[::1]:3x4/",
				},
			},
		},
		{
			description: "List with two invalid entries",
			data: TestForm{
				URLs: "ftp://test.lan:3000/\nhttp://[::1]:3x4/\n",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"URLs"},
					Classification: binding.ERR_URL,
					Message:        "ftp://test.lan:3000/",
				},
				binding.Error{
					FieldNames:     []string{"URLs"},
					Classification: binding.ERR_URL,
					Message:        "http://[::1]:3x4/",
				},
			},
		},
	}

	for _, testCase := range urlListValidationTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			performValidationTest(t, testCase)
		})
	}
}
