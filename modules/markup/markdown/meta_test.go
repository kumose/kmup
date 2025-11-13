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

package markdown

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// IssueTemplate is a legacy to keep the unit tests working.
// Copied from structs.IssueTemplate, the original type has been changed a lot to support yaml template.
type IssueTemplate struct {
	Name   string   `json:"name" yaml:"name"`
	Title  string   `json:"title" yaml:"title"`
	About  string   `json:"about" yaml:"about"`
	Labels []string `json:"labels" yaml:"labels"`
	Ref    string   `json:"ref" yaml:"ref"`
}

func (it *IssueTemplate) Valid() bool {
	return strings.TrimSpace(it.Name) != "" && strings.TrimSpace(it.About) != ""
}

func TestExtractMetadata(t *testing.T) {
	t.Run("ValidFrontAndBody", func(t *testing.T) {
		var meta IssueTemplate
		body, err := ExtractMetadata(fmt.Sprintf("%s\n%s\n%s\n%s", sepTest, frontTest, sepTest, bodyTest), &meta)
		assert.NoError(t, err)
		assert.Equal(t, bodyTest, body)
		assert.Equal(t, metaTest, meta)
		assert.True(t, meta.Valid())
	})

	t.Run("NoFirstSeparator", func(t *testing.T) {
		var meta IssueTemplate
		_, err := ExtractMetadata(fmt.Sprintf("%s\n%s\n%s", frontTest, sepTest, bodyTest), &meta)
		assert.Error(t, err)
	})

	t.Run("NoLastSeparator", func(t *testing.T) {
		var meta IssueTemplate
		_, err := ExtractMetadata(fmt.Sprintf("%s\n%s\n%s", sepTest, frontTest, bodyTest), &meta)
		assert.Error(t, err)
	})

	t.Run("NoBody", func(t *testing.T) {
		var meta IssueTemplate
		body, err := ExtractMetadata(fmt.Sprintf("%s\n%s\n%s", sepTest, frontTest, sepTest), &meta)
		assert.NoError(t, err)
		assert.Empty(t, body)
		assert.Equal(t, metaTest, meta)
		assert.True(t, meta.Valid())
	})
}

func TestExtractMetadataBytes(t *testing.T) {
	t.Run("ValidFrontAndBody", func(t *testing.T) {
		var meta IssueTemplate
		body, err := ExtractMetadataBytes(fmt.Appendf(nil, "%s\n%s\n%s\n%s", sepTest, frontTest, sepTest, bodyTest), &meta)
		assert.NoError(t, err)
		assert.Equal(t, bodyTest, string(body))
		assert.Equal(t, metaTest, meta)
		assert.True(t, meta.Valid())
	})

	t.Run("NoFirstSeparator", func(t *testing.T) {
		var meta IssueTemplate
		_, err := ExtractMetadataBytes(fmt.Appendf(nil, "%s\n%s\n%s", frontTest, sepTest, bodyTest), &meta)
		assert.Error(t, err)
	})

	t.Run("NoLastSeparator", func(t *testing.T) {
		var meta IssueTemplate
		_, err := ExtractMetadataBytes(fmt.Appendf(nil, "%s\n%s\n%s", sepTest, frontTest, bodyTest), &meta)
		assert.Error(t, err)
	})

	t.Run("NoBody", func(t *testing.T) {
		var meta IssueTemplate
		body, err := ExtractMetadataBytes(fmt.Appendf(nil, "%s\n%s\n%s", sepTest, frontTest, sepTest), &meta)
		assert.NoError(t, err)
		assert.Empty(t, string(body))
		assert.Equal(t, metaTest, meta)
		assert.True(t, meta.Valid())
	})
}

var (
	sepTest   = "-----"
	frontTest = `name: Test
about: "A Test"
title: "Test Title"
labels:
  - bug
  - "test label"`
	bodyTest = "This is the body"
	metaTest = IssueTemplate{
		Name:   "Test",
		About:  "A Test",
		Title:  "Test Title",
		Labels: []string{"bug", "test label"},
	}
)
