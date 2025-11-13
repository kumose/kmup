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

package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmitReviewForm_IsEmpty(t *testing.T) {
	cases := []struct {
		form     SubmitReviewForm
		expected bool
	}{
		// Approved PR with a comment shouldn't count as empty
		{SubmitReviewForm{Type: "approve", Content: "Awesome"}, false},

		// Approved PR without a comment shouldn't count as empty
		{SubmitReviewForm{Type: "approve", Content: ""}, false},

		// Rejected PR without a comment should count as empty
		{SubmitReviewForm{Type: "reject", Content: ""}, true},

		// Rejected PR with a comment shouldn't count as empty
		{SubmitReviewForm{Type: "reject", Content: "Awesome"}, false},

		// Comment review on a PR with a comment shouldn't count as empty
		{SubmitReviewForm{Type: "comment", Content: "Awesome"}, false},

		// Comment review on a PR without a comment should count as empty
		{SubmitReviewForm{Type: "comment", Content: ""}, true},
	}

	for _, v := range cases {
		assert.Equal(t, v.expected, v.form.HasEmptyContent())
	}
}
