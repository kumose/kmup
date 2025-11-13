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

package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getLicense(t *testing.T) {
	type args struct {
		name   string
		values *LicenseValues
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "regular",
			args: args{
				name:   "MIT",
				values: &LicenseValues{Owner: "Kmup", Year: "2023"},
			},
			want: `MIT License

Copyright (c) 2023 Kmup

Permission is hereby granted`,
			wantErr: assert.NoError,
		},
		{
			name: "license not found",
			args: args{
				name: "notfound",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLicense(tt.args.name, tt.args.values)
			if !tt.wantErr(t, err, fmt.Sprintf("GetLicense(%v, %v)", tt.args.name, tt.args.values)) {
				return
			}
			assert.Contains(t, string(got), tt.want, "GetLicense(%v, %v)", tt.args.name, tt.args.values)
		})
	}
}

func Test_fillLicensePlaceholder(t *testing.T) {
	type args struct {
		name   string
		values *LicenseValues
		origin string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "owner",
			args: args{
				name:   "regular",
				values: &LicenseValues{Year: "2023", Owner: "Kmup", Email: "teabot@kmup.io", Repo: "kmup"},
				origin: `
<name of author>
<owner>
[NAME]
[name of copyright owner]
[name of copyright holder]
<COPYRIGHT HOLDERS>
<copyright holders>
<AUTHOR>
<author's name or designee>
[one or more legally recognised persons or entities offering the Work under the terms and conditions of this Licence]
`,
			},
			want: `
Kmup
Kmup
Kmup
Kmup
Kmup
Kmup
Kmup
Kmup
Kmup
Kmup
`,
		},
		{
			name: "email",
			args: args{
				name:   "regular",
				values: &LicenseValues{Year: "2023", Owner: "Kmup", Email: "teabot@kmup.io", Repo: "kmup"},
				origin: `
[EMAIL]
`,
			},
			want: `
teabot@kmup.io
`,
		},
		{
			name: "repo",
			args: args{
				name:   "regular",
				values: &LicenseValues{Year: "2023", Owner: "Kmup", Email: "teabot@kmup.io", Repo: "kmup"},
				origin: `
<program>
<one line to give the program's name and a brief idea of what it does.>
`,
			},
			want: `
kmup
kmup
`,
		},
		{
			name: "year",
			args: args{
				name:   "regular",
				values: &LicenseValues{Year: "2023", Owner: "Kmup", Email: "teabot@kmup.io", Repo: "kmup"},
				origin: `
<year>
[YEAR]
{YEAR}
[yyyy]
[Year]
[year]
`,
			},
			want: `
2023
2023
2023
2023
2023
2023
`,
		},
		{
			name: "0BSD",
			args: args{
				name:   "0BSD",
				values: &LicenseValues{Year: "2023", Owner: "Kmup", Email: "teabot@kmup.io", Repo: "kmup"},
				origin: `
Copyright (C) YEAR by AUTHOR EMAIL

...

... THE AUTHOR BE LIABLE FOR ...
`,
			},
			want: `
Copyright (C) 2023 by Kmup teabot@kmup.io

...

... THE AUTHOR BE LIABLE FOR ...
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, string(fillLicensePlaceholder(tt.args.name, tt.args.values, []byte(tt.args.origin))), "fillLicensePlaceholder(%v, %v, %v)", tt.args.name, tt.args.values, tt.args.origin)
		})
	}
}
