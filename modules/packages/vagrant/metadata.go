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
	"compress/gzip"
	"io"
	"strings"

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/validation"
)

const (
	PropertyProvider = "vagrant.provider"
)

// Metadata represents the metadata of a Vagrant package
type Metadata struct {
	Author        string `json:"author,omitempty"`
	Description   string `json:"description,omitempty"`
	ProjectURL    string `json:"project_url,omitempty"`
	RepositoryURL string `json:"repository_url,omitempty"`
}

// ParseMetadataFromBox parses the metadata of a box file
func ParseMetadataFromBox(r io.Reader) (*Metadata, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hd, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if hd.Typeflag != tar.TypeReg {
			continue
		}

		if hd.Name == "info.json" {
			return ParseInfoFile(tr)
		}
	}

	return &Metadata{}, nil
}

// ParseInfoFile parses a info.json file to retrieve the metadata of a Vagrant package
func ParseInfoFile(r io.Reader) (*Metadata, error) {
	var values map[string]string
	if err := json.NewDecoder(r).Decode(&values); err != nil {
		return nil, err
	}

	m := &Metadata{}

	// There is no defined format for this file, just try the common keys
	for k, v := range values {
		switch strings.ToLower(k) {
		case "description":
			fallthrough
		case "short_description":
			m.Description = v
		case "website":
			fallthrough
		case "homepage":
			fallthrough
		case "url":
			if validation.IsValidURL(v) {
				m.ProjectURL = v
			}
		case "repository":
			fallthrough
		case "source":
			if validation.IsValidURL(v) {
				m.RepositoryURL = v
			}
		case "author":
			fallthrough
		case "authors":
			m.Author = v
		}
	}

	return m, nil
}
