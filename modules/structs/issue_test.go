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

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestIssueTemplate_Type(t *testing.T) {
	tests := []struct {
		fileName string
		want     IssueTemplateType
	}{
		{
			fileName: ".kmup/ISSUE_TEMPLATE/bug_report.yaml",
			want:     IssueTemplateTypeYaml,
		},
		{
			fileName: ".kmup/ISSUE_TEMPLATE/bug_report.md",
			want:     IssueTemplateTypeMarkdown,
		},
		{
			fileName: ".kmup/ISSUE_TEMPLATE/bug_report.txt",
			want:     "",
		},
		{
			fileName: ".kmup/ISSUE_TEMPLATE/config.yaml",
			want:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			it := IssueTemplate{
				FileName: tt.fileName,
			}
			assert.Equal(t, tt.want, it.Type())
		})
	}
}

func TestIssueTemplateStringSlice_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
		tmpl    *IssueTemplate
		want    *IssueTemplate
		wantErr string
	}{
		{
			name:    "array",
			content: `labels: ["a", "b", "c"]`,
			tmpl: &IssueTemplate{
				Labels: []string{"should_be_overwrote"},
			},
			want: &IssueTemplate{
				Labels: []string{"a", "b", "c"},
			},
		},
		{
			name:    "string",
			content: `labels: "a,b,c"`,
			tmpl: &IssueTemplate{
				Labels: []string{"should_be_overwrote"},
			},
			want: &IssueTemplate{
				Labels: []string{"a", "b", "c"},
			},
		},
		{
			name:    "empty",
			content: `labels:`,
			tmpl: &IssueTemplate{
				Labels: []string{"should_be_overwrote"},
			},
			want: &IssueTemplate{
				Labels: nil,
			},
		},
		{
			name: "error",
			content: `
labels:
  a: aa
  b: bb
`,
			tmpl:    &IssueTemplate{},
			wantErr: "line 3: cannot unmarshal !!map into IssueTemplateStringSlice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := yaml.Unmarshal([]byte(tt.content), tt.tmpl)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.tmpl)
			}
		})
	}
}
