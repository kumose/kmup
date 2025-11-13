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
package repo

import (
	"reflect"
	"testing"
)

func Test_localizedExtensions(t *testing.T) {
	tests := []struct {
		name              string
		ext               string
		languageCode      string
		wantLocalizedExts []string
	}{
		{
			name:              "empty language",
			ext:               ".md",
			wantLocalizedExts: []string{".md"},
		},
		{
			name:              "No region - lowercase",
			languageCode:      "en",
			ext:               ".csv",
			wantLocalizedExts: []string{".en.csv", ".csv"},
		},
		{
			name:              "No region - uppercase",
			languageCode:      "FR",
			ext:               ".txt",
			wantLocalizedExts: []string{".fr.txt", ".txt"},
		},
		{
			name:              "With region - lowercase",
			languageCode:      "en-us",
			ext:               ".md",
			wantLocalizedExts: []string{".en-us.md", ".en_us.md", ".en.md", "_en.md", ".md"},
		},
		{
			name:              "With region - uppercase",
			languageCode:      "en-CA",
			ext:               ".MD",
			wantLocalizedExts: []string{".en-ca.MD", ".en_ca.MD", ".en.MD", "_en.MD", ".MD"},
		},
		{
			name:              "With region - all uppercase",
			languageCode:      "ZH-TW",
			ext:               ".md",
			wantLocalizedExts: []string{".zh-tw.md", ".zh_tw.md", ".zh.md", "_zh.md", ".md"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLocalizedExts := localizedExtensions(tt.ext, tt.languageCode); !reflect.DeepEqual(gotLocalizedExts, tt.wantLocalizedExts) {
				t.Errorf("localizedExtensions() = %v, want %v", gotLocalizedExts, tt.wantLocalizedExts)
			}
		})
	}
}
