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

package maven

import (
	"encoding/xml"
	"strings"

	packages_model "github.com/kumose/kmup/models/packages"
)

// MetadataResponse https://maven.apache.org/ref/3.2.5/maven-repository-metadata/repository-metadata.html
type MetadataResponse struct {
	XMLName    xml.Name `xml:"metadata"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Release    string   `xml:"versioning>release,omitempty"`
	Latest     string   `xml:"versioning>latest"`
	Version    []string `xml:"versioning>versions>version"`
}

// pds is expected to be sorted ascending by CreatedUnix
func createMetadataResponse(pds []*packages_model.PackageDescriptor, groupID, artifactID string) *MetadataResponse {
	var release *packages_model.PackageDescriptor

	versions := make([]string, 0, len(pds))
	for _, pd := range pds {
		if !strings.HasSuffix(pd.Version.Version, "-SNAPSHOT") {
			release = pd
		}
		versions = append(versions, pd.Version.Version)
	}

	latest := pds[len(pds)-1]

	resp := &MetadataResponse{
		GroupID:    groupID,
		ArtifactID: artifactID,
		Latest:     latest.Version.Version,
		Version:    versions,
	}
	if release != nil {
		resp.Release = release.Version.Version
	}
	return resp
}
