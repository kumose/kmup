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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kumose-go/chi/binding"
	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testRoute = "/test"
)

type (
	validationTestCase struct {
		description    string
		data           any
		expectedErrors binding.Errors
	}

	TestForm struct {
		BranchName   string `form:"BranchName" binding:"GitRefName"`
		URL          string `form:"ValidUrl" binding:"ValidUrl"`
		URLs         string `form:"ValidUrls" binding:"ValidUrlList"`
		GlobPattern  string `form:"GlobPattern" binding:"GlobPattern"`
		RegexPattern string `form:"RegexPattern" binding:"RegexPattern"`
	}
)

func performValidationTest(t *testing.T, testCase validationTestCase) {
	httpRecorder := httptest.NewRecorder()
	m := chi.NewRouter()

	m.Post(testRoute, func(resp http.ResponseWriter, req *http.Request) {
		actual := binding.Validate(req, testCase.data)
		// see https://github.com/stretchr/testify/issues/435
		if actual == nil {
			actual = binding.Errors{}
		}

		assert.Equal(t, testCase.expectedErrors, actual)
	})

	req, err := http.NewRequest(http.MethodPost, testRoute, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "x-www-form-urlencoded")
	m.ServeHTTP(httpRecorder, req)

	switch httpRecorder.Code {
	case http.StatusNotFound:
		panic("Routing is messed up in test fixture (got 404): check methods and paths")
	case http.StatusInternalServerError:
		panic("Something bad happened on '" + testCase.description + "'")
	}
}
