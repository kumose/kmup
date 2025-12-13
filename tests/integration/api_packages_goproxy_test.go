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

package integration

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/kumose/kmup/models/packages"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestPackageGo(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	packageName := "kmup.com/kumose/kmup"
	packageVersion := "v0.0.1"
	packageVersion2 := "v0.0.2"
	goModContent := `module "kmup.com/kumose/kmup"`

	createArchive := func(files map[string][]byte) []byte {
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		for name, content := range files {
			w, _ := zw.Create(name)
			w.Write(content)
		}
		zw.Close()
		return buf.Bytes()
	}

	url := fmt.Sprintf("/api/packages/%s/go", user.Name)

	t.Run("Upload", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		content := createArchive(nil)

		req := NewRequestWithBody(t, "PUT", url+"/upload", bytes.NewReader(content))
		MakeRequest(t, req, http.StatusUnauthorized)

		req = NewRequestWithBody(t, "PUT", url+"/upload", bytes.NewReader(content)).
			AddBasicAuth(user.Name)
		MakeRequest(t, req, http.StatusBadRequest)

		content = createArchive(map[string][]byte{
			packageName + "@" + packageVersion + "/go.mod": []byte(goModContent),
		})

		req = NewRequestWithBody(t, "PUT", url+"/upload", bytes.NewReader(content)).
			AddBasicAuth(user.Name)
		MakeRequest(t, req, http.StatusCreated)

		pvs, err := packages.GetVersionsByPackageType(t.Context(), user.ID, packages.TypeGo)
		assert.NoError(t, err)
		assert.Len(t, pvs, 1)

		pd, err := packages.GetPackageDescriptor(t.Context(), pvs[0])
		assert.NoError(t, err)
		assert.Nil(t, pd.Metadata)
		assert.Equal(t, packageName, pd.Package.Name)
		assert.Equal(t, packageVersion, pd.Version.Version)

		pfs, err := packages.GetFilesByVersionID(t.Context(), pvs[0].ID)
		assert.NoError(t, err)
		assert.Len(t, pfs, 1)
		assert.Equal(t, packageVersion+".zip", pfs[0].Name)
		assert.True(t, pfs[0].IsLead)

		pb, err := packages.GetBlobByID(t.Context(), pfs[0].BlobID)
		assert.NoError(t, err)
		assert.Equal(t, int64(len(content)), pb.Size)

		req = NewRequestWithBody(t, "PUT", url+"/upload", bytes.NewReader(content)).
			AddBasicAuth(user.Name)
		MakeRequest(t, req, http.StatusConflict)

		time.Sleep(time.Second) // Ensure the timestamp is different, then the "list" below can have stable order

		content = createArchive(map[string][]byte{
			packageName + "@" + packageVersion2 + "/go.mod": []byte(goModContent),
		})

		req = NewRequestWithBody(t, "PUT", url+"/upload", bytes.NewReader(content)).
			AddBasicAuth(user.Name)
		MakeRequest(t, req, http.StatusCreated)
	})

	t.Run("List", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequest(t, "GET", fmt.Sprintf("%s/%s/@v/list", url, packageName))
		resp := MakeRequest(t, req, http.StatusOK)

		assert.Equal(t, packageVersion+"\n"+packageVersion2+"\n", resp.Body.String())
	})

	t.Run("Info", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequest(t, "GET", fmt.Sprintf("%s/%s/@v/%s.info", url, packageName, packageVersion))
		resp := MakeRequest(t, req, http.StatusOK)

		type Info struct {
			Version string    `json:"Version"`
			Time    time.Time `json:"Time"`
		}

		info := &Info{}
		DecodeJSON(t, resp, &info)

		assert.Equal(t, packageVersion, info.Version)

		req = NewRequest(t, "GET", fmt.Sprintf("%s/%s/@v/latest.info", url, packageName))
		resp = MakeRequest(t, req, http.StatusOK)

		info = &Info{}
		DecodeJSON(t, resp, &info)

		assert.Equal(t, packageVersion2, info.Version)

		req = NewRequest(t, "GET", fmt.Sprintf("%s/%s/@latest", url, packageName))
		resp = MakeRequest(t, req, http.StatusOK)

		info = &Info{}
		DecodeJSON(t, resp, &info)

		assert.Equal(t, packageVersion2, info.Version)
	})

	t.Run("GoMod", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequest(t, "GET", fmt.Sprintf("%s/%s/@v/%s.mod", url, packageName, packageVersion))
		resp := MakeRequest(t, req, http.StatusOK)

		assert.Equal(t, goModContent, resp.Body.String())

		req = NewRequest(t, "GET", fmt.Sprintf("%s/%s/@v/latest.mod", url, packageName))
		resp = MakeRequest(t, req, http.StatusOK)

		assert.Equal(t, goModContent, resp.Body.String())
	})

	t.Run("Download", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequest(t, "GET", fmt.Sprintf("%s/%s/@v/%s.zip", url, packageName, packageVersion))
		MakeRequest(t, req, http.StatusOK)

		req = NewRequest(t, "GET", fmt.Sprintf("%s/%s/@v/latest.zip", url, packageName))
		MakeRequest(t, req, http.StatusOK)
	})
}
