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

package alpine

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	packageName        = "kmup"
	packageVersion     = "1.0.1"
	packageDescription = "Package Description"
	packageProjectURL  = "https://kmup.io"
	packageMaintainer  = "KN4CK3R <dummy@kmup.io>"
)

func createPKGINFOContent(name, version string) []byte {
	return []byte(`pkgname = ` + name + `
pkgver = ` + version + `
pkgdesc = ` + packageDescription + `
url = ` + packageProjectURL + `
# comment
builddate = 1678834800
packager = Kmup <pack@ag.er>
size = 123456
arch = aarch64
origin = origin
commit = 1111e709613fbc979651b09ac2bc27c6591a9999
maintainer = ` + packageMaintainer + `
license = MIT
depend = common
install_if = value
depend = kmup
provides = common
provides = kmup`)
}

func TestParsePackage(t *testing.T) {
	createPackage := func(name string, content []byte) io.Reader {
		names := []string{"first.stream", name}
		contents := [][]byte{{0}, content}

		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)

		for i := range names {
			if i != 0 {
				zw.Close()
				zw.Reset(&buf)
			}

			tw := tar.NewWriter(zw)
			hdr := &tar.Header{
				Name: names[i],
				Mode: 0o600,
				Size: int64(len(contents[i])),
			}
			tw.WriteHeader(hdr)
			tw.Write(contents[i])
			tw.Close()
		}

		zw.Close()

		return &buf
	}

	t.Run("MissingPKGINFOFile", func(t *testing.T) {
		data := createPackage("dummy.txt", []byte{})

		pp, err := ParsePackage(data)
		assert.Nil(t, pp)
		assert.ErrorIs(t, err, ErrMissingPKGINFOFile)
	})

	t.Run("InvalidPKGINFOFile", func(t *testing.T) {
		data := createPackage(".PKGINFO", []byte{})

		pp, err := ParsePackage(data)
		assert.Nil(t, pp)
		assert.ErrorIs(t, err, ErrInvalidName)
	})

	t.Run("Valid", func(t *testing.T) {
		data := createPackage(".PKGINFO", createPKGINFOContent(packageName, packageVersion))

		p, err := ParsePackage(data)
		assert.NoError(t, err)
		assert.NotNil(t, p)

		assert.Equal(t, "Q1SRYURM5+uQDqfHSwTnNIOIuuDVQ=", p.FileMetadata.Checksum)
	})
}

func TestParsePackageInfo(t *testing.T) {
	t.Run("InvalidName", func(t *testing.T) {
		data := createPKGINFOContent("", packageVersion)

		p, err := ParsePackageInfo(bytes.NewReader(data))
		assert.Nil(t, p)
		assert.ErrorIs(t, err, ErrInvalidName)
	})

	t.Run("InvalidVersion", func(t *testing.T) {
		data := createPKGINFOContent(packageName, "")

		p, err := ParsePackageInfo(bytes.NewReader(data))
		assert.Nil(t, p)
		assert.ErrorIs(t, err, ErrInvalidVersion)
	})

	t.Run("Valid", func(t *testing.T) {
		data := createPKGINFOContent(packageName, packageVersion)

		p, err := ParsePackageInfo(bytes.NewReader(data))
		assert.NoError(t, err)
		assert.NotNil(t, p)

		assert.Equal(t, packageName, p.Name)
		assert.Equal(t, packageVersion, p.Version)
		assert.Equal(t, packageDescription, p.VersionMetadata.Description)
		assert.Equal(t, packageMaintainer, p.VersionMetadata.Maintainer)
		assert.Equal(t, packageProjectURL, p.VersionMetadata.ProjectURL)
		assert.Equal(t, "MIT", p.VersionMetadata.License)
		assert.Empty(t, p.FileMetadata.Checksum)
		assert.Equal(t, "Kmup <pack@ag.er>", p.FileMetadata.Packager)
		assert.EqualValues(t, 1678834800, p.FileMetadata.BuildDate)
		assert.EqualValues(t, 123456, p.FileMetadata.Size)
		assert.Equal(t, "aarch64", p.FileMetadata.Architecture)
		assert.Equal(t, "origin", p.FileMetadata.Origin)
		assert.Equal(t, "1111e709613fbc979651b09ac2bc27c6591a9999", p.FileMetadata.CommitHash)
		assert.Equal(t, "value", p.FileMetadata.InstallIf)
		assert.ElementsMatch(t, []string{"common", "kmup"}, p.FileMetadata.Provides)
		assert.ElementsMatch(t, []string{"common", "kmup"}, p.FileMetadata.Dependencies)
	})
}
