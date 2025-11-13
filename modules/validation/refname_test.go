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

func Test_GitRefNameValidation(t *testing.T) {
	AddBindingRules()
	gitRefNameValidationTestCases := []validationTestCase{
		{
			description: "Reference name contains only characters",
			data: TestForm{
				BranchName: "test",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "Reference name contains single slash",
			data: TestForm{
				BranchName: "feature/test",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "Reference name has allowed special characters",
			data: TestForm{
				BranchName: "debian/1%1.6.0-2",
			},
			expectedErrors: binding.Errors{},
		},
		{
			description: "Reference name contains backslash",
			data: TestForm{
				BranchName: "feature\\test",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name starts with dot",
			data: TestForm{
				BranchName: ".test",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name ends with dot",
			data: TestForm{
				BranchName: "test.",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name starts with slash",
			data: TestForm{
				BranchName: "/test",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name ends with slash",
			data: TestForm{
				BranchName: "test/",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name ends with .lock",
			data: TestForm{
				BranchName: "test.lock",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name contains multiple consecutive dots",
			data: TestForm{
				BranchName: "te..st",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name contains multiple consecutive slashes",
			data: TestForm{
				BranchName: "te//st",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name is single @",
			data: TestForm{
				BranchName: "@",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name has @{",
			data: TestForm{
				BranchName: "branch@{",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name has unallowed special character ~",
			data: TestForm{
				BranchName: "~debian/1%1.6.0-2",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name has unallowed special character *",
			data: TestForm{
				BranchName: "*debian/1%1.6.0-2",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name has unallowed special character ?",
			data: TestForm{
				BranchName: "?debian/1%1.6.0-2",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name has unallowed special character ^",
			data: TestForm{
				BranchName: "^debian/1%1.6.0-2",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name has unallowed special character :",
			data: TestForm{
				BranchName: "debian:jessie",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name has unallowed special character (whitespace)",
			data: TestForm{
				BranchName: "debian jessie",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
		{
			description: "Reference name has unallowed special character [",
			data: TestForm{
				BranchName: "debian[jessie",
			},
			expectedErrors: binding.Errors{
				binding.Error{
					FieldNames:     []string{"BranchName"},
					Classification: ErrGitRefName,
					Message:        "GitRefName",
				},
			},
		},
	}

	for _, testCase := range gitRefNameValidationTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			performValidationTest(t, testCase)
		})
	}
}
