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

package rpm

import (
	"bytes"
	"io"
	"strings"

	packages_module "github.com/kumose/kmup/modules/packages"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/sassoftware/go-rpmutils"
)

func SignPackage(buf *packages_module.HashedBuffer, privateKey string) (*packages_module.HashedBuffer, error) {
	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(privateKey))
	if err != nil {
		return nil, err
	}

	h, err := rpmutils.SignRpmStream(buf, keyring[0].PrivateKey, nil)
	if err != nil {
		return nil, err
	}

	signBlob, err := h.DumpSignatureHeader(false)
	if err != nil {
		return nil, err
	}

	if _, err := buf.Seek(int64(h.OriginalSignatureHeaderSize()), io.SeekStart); err != nil {
		return nil, err
	}

	// create new buf with signature prefix
	return packages_module.CreateHashedBufferFromReader(io.MultiReader(bytes.NewReader(signBlob), buf))
}
