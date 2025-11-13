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

package markup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescriptionSanitizer(t *testing.T) {
	testCases := []string{
		`<h1>Title</h1>`, `Title`,
		`<img src='img.png' alt='image'>`, ``,
		`<span class="emoji" aria-label="thumbs up">THUMBS UP</span>`, `<span class="emoji" aria-label="thumbs up">THUMBS UP</span>`,
		`<span style="color: red">Hello World</span>`, `<span>Hello World</span>`,
		`<br>`, ``,
		`<a href="https://example.com" target="_blank" rel="noopener noreferrer">https://example.com</a>`, `<a href="https://example.com" target="_blank" rel="noopener noreferrer nofollow">https://example.com</a>`,
		`<a href="data:1234">data</a>`, `data`,
		`<mark>Important!</mark>`, `Important!`,
		`<details>Click me! <summary>Nothing to see here.</summary></details>`, `Click me! Nothing to see here.`,
		`<input type="hidden">`, ``,
		`<b>I</b> have a <i>strong</i> <strong>opinion</strong> about <em>this</em>.`, `<b>I</b> have a <i>strong</i> <strong>opinion</strong> about <em>this</em>.`,
		`Provides alternative <code>wg(8)</code> tool`, `Provides alternative <code>wg(8)</code> tool`,
	}

	for i := 0; i < len(testCases); i += 2 {
		assert.Equal(t, testCases[i+1], SanitizeDescription(testCases[i]))
	}
}
