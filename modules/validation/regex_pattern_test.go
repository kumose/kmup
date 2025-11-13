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
	"regexp"
	"testing"

	"github.com/kumose-go/chi/binding"
)

func getRegexPatternErrorString(pattern string) string {
	if _, err := regexp.Compile(pattern); err != nil {
		return err.Error()
	}
	return ""
}

func Test_RegexPatternValidation(t *testing.T) {
	AddBindingRules()

	regexValidationTestCases := []validationTestCase{
		{
			description: "Empty regex pattern",
			data: TestForm{
				RegexPattern: "",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "Valid regex",
			data: TestForm{
				RegexPattern: `(\d{1,3})+`,
			},
			expectedErrors: binding.Errors{},
		},

		{
			description: "Invalid regex",
			data: TestForm{
				RegexPattern: "[a-",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"RegexPattern"},
					Classification: ErrRegexPattern,
					Message:        getRegexPatternErrorString("[a-"),
				},
			},
		},
	}

	for _, testCase := range regexValidationTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			performValidationTest(t, testCase)
		})
	}
}
