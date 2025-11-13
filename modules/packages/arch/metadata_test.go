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

package arch

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/ulikunitz/xz"
)

const (
	packageName        = "kmup"
	packageVersion     = "1.0.1"
	packageDescription = "Package Description"
	packageProjectURL  = "https://kmup.com"
	packagePackager    = "KN4CK3R <packager@kmup.com>"
)

func createPKGINFOContent(name, version string) []byte {
	return []byte(`pkgname = ` + name + `
pkgbase = ` + name + `
pkgver = ` + version + `
pkgdesc = ` + packageDescription + `
url = ` + packageProjectURL + `
# comment
group=group
builddate = 1678834800
size = 123456
arch = x86_64
license = MIT
packager = ` + packagePackager + `
depend = common
xdata = value
depend = kmup
provides = common
provides = kmup
optdepend = hex
replaces = gogs
checkdepend = common
makedepend = cmake
conflict = ninja
backup = usr/bin/paket1`)
}

func TestParsePackage(t *testing.T) {
	createPackage := func(compression string, files map[string][]byte) io.Reader {
		var buf bytes.Buffer
		var cw io.WriteCloser
		switch compression {
		case "zst":
			cw, _ = zstd.NewWriter(&buf)
		case "xz":
			cw, _ = xz.NewWriter(&buf)
		case "gz":
			cw = gzip.NewWriter(&buf)
		}
		tw := tar.NewWriter(cw)

		for name, content := range files {
			hdr := &tar.Header{
				Name: name,
				Mode: 0o600,
				Size: int64(len(content)),
			}
			tw.WriteHeader(hdr)
			tw.Write(content)
		}

		tw.Close()
		cw.Close()

		return &buf
	}

	for _, c := range []string{"gz", "xz", "zst"} {
		t.Run(c, func(t *testing.T) {
			t.Run("MissingPKGINFOFile", func(t *testing.T) {
				data := createPackage(c, map[string][]byte{"dummy.txt": {}})

				pp, err := ParsePackage(data)
				assert.Nil(t, pp)
				assert.ErrorIs(t, err, ErrMissingPKGINFOFile)
			})

			t.Run("InvalidPKGINFOFile", func(t *testing.T) {
				data := createPackage(c, map[string][]byte{".PKGINFO": {}})

				pp, err := ParsePackage(data)
				assert.Nil(t, pp)
				assert.ErrorIs(t, err, ErrInvalidName)
			})

			t.Run("Valid", func(t *testing.T) {
				data := createPackage(c, map[string][]byte{
					".PKGINFO":        createPKGINFOContent(packageName, packageVersion),
					"/test/dummy.txt": {},
				})

				p, err := ParsePackage(data)
				assert.NoError(t, err)
				assert.NotNil(t, p)

				assert.ElementsMatch(t, []string{"/test/dummy.txt"}, p.FileMetadata.Files)
			})
		})
	}
}

func TestParsePackageInfo(t *testing.T) {
	t.Run("InvalidName", func(t *testing.T) {
		data := createPKGINFOContent("", packageVersion)

		p, err := ParsePackageInfo(bytes.NewReader(data))
		assert.Nil(t, p)
		assert.ErrorIs(t, err, ErrInvalidName)
	})

	t.Run("Regexp", func(t *testing.T) {
		assert.Regexp(t, versionPattern, "1.2_3~4+5")
		assert.Regexp(t, versionPattern, "1:2_3~4+5")
		assert.NotRegexp(t, versionPattern, "a:1.0.0-1")
		assert.NotRegexp(t, versionPattern, "0.0.1/1-1")
		assert.NotRegexp(t, versionPattern, "1.0.0 -1")
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
		assert.Equal(t, packageName, p.FileMetadata.Base)
		assert.Equal(t, packageVersion, p.Version)
		assert.Equal(t, packageDescription, p.VersionMetadata.Description)
		assert.Equal(t, packagePackager, p.FileMetadata.Packager)
		assert.Equal(t, packageProjectURL, p.VersionMetadata.ProjectURL)
		assert.ElementsMatch(t, []string{"MIT"}, p.VersionMetadata.Licenses)
		assert.EqualValues(t, 1678834800, p.FileMetadata.BuildDate)
		assert.EqualValues(t, 123456, p.FileMetadata.InstalledSize)
		assert.Equal(t, "x86_64", p.FileMetadata.Architecture)
		assert.ElementsMatch(t, []string{"value"}, p.FileMetadata.XData)
		assert.ElementsMatch(t, []string{"group"}, p.FileMetadata.Groups)
		assert.ElementsMatch(t, []string{"common", "kmup"}, p.FileMetadata.Provides)
		assert.ElementsMatch(t, []string{"common", "kmup"}, p.FileMetadata.Depends)
		assert.ElementsMatch(t, []string{"gogs"}, p.FileMetadata.Replaces)
		assert.ElementsMatch(t, []string{"hex"}, p.FileMetadata.OptDepends)
		assert.ElementsMatch(t, []string{"common"}, p.FileMetadata.CheckDepends)
		assert.ElementsMatch(t, []string{"ninja"}, p.FileMetadata.Conflicts)
		assert.ElementsMatch(t, []string{"cmake"}, p.FileMetadata.MakeDepends)
		assert.ElementsMatch(t, []string{"usr/bin/paket1"}, p.FileMetadata.Backup)
	})
}
