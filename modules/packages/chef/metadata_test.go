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

package chef

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	packageName          = "kmup"
	packageVersion       = "1.0.1"
	packageAuthor        = "KN4CK3R"
	packageDescription   = "Package Description"
	packageRepositoryURL = "https://github.com/kumose/kmup"
)

func TestParsePackage(t *testing.T) {
	t.Run("MissingMetadataFile", func(t *testing.T) {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(zw)
		tw.Close()
		zw.Close()

		p, err := ParsePackage(&buf)
		assert.Nil(t, p)
		assert.ErrorIs(t, err, ErrMissingMetadataFile)
	})

	t.Run("Valid", func(t *testing.T) {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(zw)

		content := `{"name":"` + packageName + `","version":"` + packageVersion + `"}`

		hdr := &tar.Header{
			Name: packageName + "/metadata.json",
			Mode: 0o600,
			Size: int64(len(content)),
		}
		tw.WriteHeader(hdr)
		tw.Write([]byte(content))

		tw.Close()
		zw.Close()

		p, err := ParsePackage(&buf)
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, packageName, p.Name)
		assert.Equal(t, packageVersion, p.Version)
		assert.NotNil(t, p.Metadata)
	})
}

func TestParseChefMetadata(t *testing.T) {
	t.Run("InvalidName", func(t *testing.T) {
		for _, name := range []string{" test", "test "} {
			p, err := ParseChefMetadata(strings.NewReader(`{"name":"` + name + `","version":"1.0.0"}`))
			assert.Nil(t, p)
			assert.ErrorIs(t, err, ErrInvalidName)
		}
	})

	t.Run("InvalidVersion", func(t *testing.T) {
		for _, version := range []string{"1", "1.2.3.4", "1.0.0 "} {
			p, err := ParseChefMetadata(strings.NewReader(`{"name":"test","version":"` + version + `"}`))
			assert.Nil(t, p)
			assert.ErrorIs(t, err, ErrInvalidVersion)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		p, err := ParseChefMetadata(strings.NewReader(`{"name":"` + packageName + `","version":"` + packageVersion + `","description":"` + packageDescription + `","maintainer":"` + packageAuthor + `","source_url":"` + packageRepositoryURL + `"}`))
		assert.NotNil(t, p)
		assert.NoError(t, err)

		assert.Equal(t, packageName, p.Name)
		assert.Equal(t, packageVersion, p.Version)
		assert.Equal(t, packageDescription, p.Metadata.Description)
		assert.Equal(t, packageAuthor, p.Metadata.Author)
		assert.Equal(t, packageRepositoryURL, p.Metadata.RepositoryURL)
	})
}
