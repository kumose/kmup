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

package setting

import (
	"github.com/kumose/kmup/modules/log"

	"github.com/42wim/httpsig"
)

// Federation settings
var (
	Federation = struct {
		Enabled             bool
		ShareUserStatistics bool
		MaxSize             int64
		Algorithms          []string
		DigestAlgorithm     string
		GetHeaders          []string
		PostHeaders         []string
	}{
		Enabled:             false,
		ShareUserStatistics: true,
		MaxSize:             4,
		Algorithms:          []string{"rsa-sha256", "rsa-sha512", "ed25519"},
		DigestAlgorithm:     "SHA-256",
		GetHeaders:          []string{"(request-target)", "Date"},
		PostHeaders:         []string{"(request-target)", "Date", "Digest"},
	}
)

// HttpsigAlgs is a constant slice of httpsig algorithm objects
var HttpsigAlgs []httpsig.Algorithm

func loadFederationFrom(rootCfg ConfigProvider) {
	if err := rootCfg.Section("federation").MapTo(&Federation); err != nil {
		log.Fatal("Failed to map Federation settings: %v", err)
	} else if !httpsig.IsSupportedDigestAlgorithm(Federation.DigestAlgorithm) {
		log.Fatal("unsupported digest algorithm: %s", Federation.DigestAlgorithm)
		return
	}

	// Get MaxSize in bytes instead of MiB
	Federation.MaxSize = 1 << 20 * Federation.MaxSize

	HttpsigAlgs = make([]httpsig.Algorithm, len(Federation.Algorithms))
	for i, alg := range Federation.Algorithms {
		HttpsigAlgs[i] = httpsig.Algorithm(alg)
	}
}
