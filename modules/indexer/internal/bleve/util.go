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

package bleve

import (
	"errors"
	"os"
	"unicode"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"

	"github.com/blevesearch/bleve/v2"
	unicode_tokenizer "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/index/upsidedown"
	"github.com/ethantkoenig/rupture"
)

const (
	maxFuzziness = 2
)

// openIndexer open the index at the specified path, checking for metadata
// updates and bleve version updates.  If index needs to be created (or
// re-created), returns (nil, nil)
func openIndexer(path string, latestVersion int) (bleve.Index, int, error) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	metadata, err := rupture.ReadIndexMetadata(path)
	if err != nil {
		return nil, 0, err
	}
	if metadata.Version < latestVersion {
		// the indexer is using a previous version, so we should delete it and
		// re-populate
		return nil, metadata.Version, util.RemoveAll(path)
	}

	index, err := bleve.Open(path)
	if err != nil {
		if errors.Is(err, upsidedown.IncompatibleVersion) {
			log.Warn("Indexer was built with a previous version of bleve, deleting and rebuilding")
			return nil, 0, util.RemoveAll(path)
		}
		return nil, 0, err
	}

	return index, 0, nil
}

// GuessFuzzinessByKeyword guesses fuzziness based on the levenshtein distance and determines how many chars
// may be different on two string, and they still be considered equivalent.
// Given a phrase, its shortest word determines its fuzziness. If a phrase uses CJK (eg: `갃갃갃` `啊啊啊`), the fuzziness is zero.
func GuessFuzzinessByKeyword(s string) int {
	tokenizer := unicode_tokenizer.NewUnicodeTokenizer()
	tokens := tokenizer.Tokenize([]byte(s))

	if len(tokens) > 0 {
		fuzziness := maxFuzziness

		for _, token := range tokens {
			fuzziness = min(fuzziness, guessFuzzinessByKeyword(string(token.Term)))
		}

		return fuzziness
	}

	return 0
}

func guessFuzzinessByKeyword(s string) int {
	// according to https://github.com/blevesearch/bleve/issues/1563, the supported max fuzziness is 2
	// magic number 4 was chosen to determine the levenshtein distance per each character of a keyword
	// BUT, when using CJK (eg: `갃갃갃` `啊啊啊`), it mismatches a lot.
	// Likewise, queries whose terms contains characters that are *not* letters should not use fuzziness

	for _, r := range s {
		if r >= 128 || !unicode.IsLetter(r) {
			return 0
		}
	}
	return min(min(setting.Indexer.TypeBleveMaxFuzzniess, maxFuzziness), len(s)/4)
}
