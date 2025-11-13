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

package composer

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/kumose/kmup/modules/json"

	"github.com/stretchr/testify/assert"
)

const (
	name        = "kmup/composer-package"
	description = "Package Description"
	readme      = "Package Readme"
	comments    = "Package Comment"
	packageType = "composer-plugin"
	author      = "Kmup Authors"
	email       = "no.reply@kmup.io"
	homepage    = "https://kmup.io"
	license     = "MIT"
)

const composerContent = `{
    "name": "` + name + `",
    "description": "` + description + `",
    "type": "` + packageType + `",
    "license": "` + license + `",
    "authors": [
        {
            "name": "` + author + `",
            "email": "` + email + `"
        }
    ],
    "homepage": "` + homepage + `",
    "autoload": {
        "psr-4": {"Kmup\\ComposerPackage\\": "src/"}
    },
    "require": {
        "php": ">=7.2 || ^8.0"
    },
    "_comments": "` + comments + `"
}`

func TestLicenseUnmarshal(t *testing.T) {
	var l Licenses
	assert.NoError(t, json.NewDecoder(strings.NewReader(`["MIT"]`)).Decode(&l))
	assert.Len(t, l, 1)
	assert.Equal(t, "MIT", l[0])
	assert.NoError(t, json.NewDecoder(strings.NewReader(`"MIT"`)).Decode(&l))
	assert.Len(t, l, 1)
	assert.Equal(t, "MIT", l[0])
}

func TestCommentsUnmarshal(t *testing.T) {
	var c Comments
	assert.NoError(t, json.NewDecoder(strings.NewReader(`["comment"]`)).Decode(&c))
	assert.Len(t, c, 1)
	assert.Equal(t, "comment", c[0])
	assert.NoError(t, json.NewDecoder(strings.NewReader(`"comment"`)).Decode(&c))
	assert.Len(t, c, 1)
	assert.Equal(t, "comment", c[0])
}

func TestParsePackage(t *testing.T) {
	createArchive := func(files map[string]string) []byte {
		var buf bytes.Buffer
		archive := zip.NewWriter(&buf)
		for name, content := range files {
			w, _ := archive.Create(name)
			w.Write([]byte(content))
		}
		archive.Close()
		return buf.Bytes()
	}

	t.Run("MissingComposerFile", func(t *testing.T) {
		data := createArchive(map[string]string{"dummy.txt": ""})

		cp, err := ParsePackage(bytes.NewReader(data), int64(len(data)))
		assert.Nil(t, cp)
		assert.ErrorIs(t, err, ErrMissingComposerFile)
	})

	t.Run("MissingComposerFileInRoot", func(t *testing.T) {
		data := createArchive(map[string]string{"sub/sub/composer.json": ""})

		cp, err := ParsePackage(bytes.NewReader(data), int64(len(data)))
		assert.Nil(t, cp)
		assert.ErrorIs(t, err, ErrMissingComposerFile)
	})

	t.Run("InvalidComposerFile", func(t *testing.T) {
		data := createArchive(map[string]string{"composer.json": ""})

		cp, err := ParsePackage(bytes.NewReader(data), int64(len(data)))
		assert.Nil(t, cp)
		assert.Error(t, err)
	})

	t.Run("InvalidPackageName", func(t *testing.T) {
		data := createArchive(map[string]string{"composer.json": "{}"})

		cp, err := ParsePackage(bytes.NewReader(data), int64(len(data)))
		assert.Nil(t, cp)
		assert.ErrorIs(t, err, ErrInvalidName)
	})

	t.Run("InvalidPackageVersion", func(t *testing.T) {
		data := createArchive(map[string]string{"composer.json": `{"name": "kmup/composer-package", "version": "1.a.3"}`})

		cp, err := ParsePackage(bytes.NewReader(data), int64(len(data)))
		assert.Nil(t, cp)
		assert.ErrorIs(t, err, ErrInvalidVersion)
	})

	t.Run("InvalidReadmePath", func(t *testing.T) {
		data := createArchive(map[string]string{"composer.json": `{"name": "kmup/composer-package", "readme": "sub/README.md"}`})

		cp, err := ParsePackage(bytes.NewReader(data), int64(len(data)))
		assert.NoError(t, err)
		assert.NotNil(t, cp)

		assert.Empty(t, cp.Metadata.Readme)
	})

	t.Run("Valid", func(t *testing.T) {
		data := createArchive(map[string]string{"composer.json": composerContent, "README.md": readme})

		cp, err := ParsePackage(bytes.NewReader(data), int64(len(data)))
		assert.NoError(t, err)
		assert.NotNil(t, cp)

		assert.Equal(t, name, cp.Name)
		assert.Empty(t, cp.Version)
		assert.Equal(t, description, cp.Metadata.Description)
		assert.Equal(t, readme, cp.Metadata.Readme)
		assert.Len(t, cp.Metadata.Comments, 1)
		assert.Equal(t, comments, cp.Metadata.Comments[0])
		assert.Len(t, cp.Metadata.Authors, 1)
		assert.Equal(t, author, cp.Metadata.Authors[0].Name)
		assert.Equal(t, email, cp.Metadata.Authors[0].Email)
		assert.Equal(t, homepage, cp.Metadata.Homepage)
		assert.Equal(t, packageType, cp.Type)
		assert.Len(t, cp.Metadata.License, 1)
		assert.Equal(t, license, cp.Metadata.License[0])
	})
}
