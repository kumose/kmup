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

	"github.com/kumose/kmup/modules/glob"

	"github.com/kumose-go/chi/binding"
)

func getGlobPatternErrorString(pattern string) string {
	// It would be unwise to rely on that glob
	// compilation errors don't ever change.
	if _, err := glob.Compile(pattern); err != nil {
		return err.Error()
	}
	return ""
}

func Test_GlobPatternValidation(t *testing.T) {
	AddBindingRules()

	globValidationTestCases := []validationTestCase{
		{
			description: "Empty glob pattern",
			data: TestForm{
				GlobPattern: "",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "Valid glob",
			data: TestForm{
				GlobPattern: "{master,release*}",
			},
			expectedErrors: binding.Errors{},
		},

		{
			description: "Invalid glob",
			data: TestForm{
				GlobPattern: "[a-",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"GlobPattern"},
					Classification: ErrGlobPattern,
					Message:        getGlobPatternErrorString("[a-"),
				},
			},
		},
	}

	for _, testCase := range globValidationTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			performValidationTest(t, testCase)
		})
	}
}
