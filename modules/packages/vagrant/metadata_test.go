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

package vagrant

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/kumose/kmup/modules/json"

	"github.com/stretchr/testify/assert"
)

const (
	author        = "kmup"
	description   = "Package Description"
	projectURL    = "https://kmup.io"
	repositoryURL = "https://kmup.io/kmup/kmup"
)

func TestParseMetadataFromBox(t *testing.T) {
	createArchive := func(files map[string][]byte) io.Reader {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(zw)
		for filename, content := range files {
			hdr := &tar.Header{
				Name: filename,
				Mode: 0o600,
				Size: int64(len(content)),
			}
			tw.WriteHeader(hdr)
			tw.Write(content)
		}
		tw.Close()
		zw.Close()
		return &buf
	}

	t.Run("MissingInfoFile", func(t *testing.T) {
		data := createArchive(map[string][]byte{"dummy.txt": {}})

		metadata, err := ParseMetadataFromBox(data)
		assert.NotNil(t, metadata)
		assert.NoError(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		content, err := json.Marshal(map[string]string{
			"description": description,
			"author":      author,
			"website":     projectURL,
			"repository":  repositoryURL,
		})
		assert.NoError(t, err)

		data := createArchive(map[string][]byte{"info.json": content})

		metadata, err := ParseMetadataFromBox(data)
		assert.NotNil(t, metadata)
		assert.NoError(t, err)

		assert.Equal(t, author, metadata.Author)
		assert.Equal(t, description, metadata.Description)
		assert.Equal(t, projectURL, metadata.ProjectURL)
		assert.Equal(t, repositoryURL, metadata.RepositoryURL)
	})
}

func TestParseInfoFile(t *testing.T) {
	t.Run("UnknownKeys", func(t *testing.T) {
		content, err := json.Marshal(map[string]string{
			"package": "",
			"dummy":   "",
		})
		assert.NoError(t, err)

		metadata, err := ParseInfoFile(bytes.NewReader(content))
		assert.NotNil(t, metadata)
		assert.NoError(t, err)

		assert.Empty(t, metadata.Author)
		assert.Empty(t, metadata.Description)
		assert.Empty(t, metadata.ProjectURL)
		assert.Empty(t, metadata.RepositoryURL)
	})

	t.Run("Valid", func(t *testing.T) {
		content, err := json.Marshal(map[string]string{
			"description": description,
			"author":      author,
			"website":     projectURL,
			"repository":  repositoryURL,
		})
		assert.NoError(t, err)

		metadata, err := ParseInfoFile(bytes.NewReader(content))
		assert.NotNil(t, metadata)
		assert.NoError(t, err)

		assert.Equal(t, author, metadata.Author)
		assert.Equal(t, description, metadata.Description)
		assert.Equal(t, projectURL, metadata.ProjectURL)
		assert.Equal(t, repositoryURL, metadata.RepositoryURL)
	})
}
