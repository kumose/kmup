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

package markup_test

import (
	"io"
	"strings"
	"testing"

	"github.com/kumose/kmup/modules/emoji"
	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/markup/markdown"
	"github.com/kumose/kmup/modules/setting"
	testModule "github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
)

var (
	testRepoOwnerName = "user13"
	testRepoName      = "repo11"
	localMetas        = map[string]string{"user": testRepoOwnerName, "repo": testRepoName}
)

func TestRender_Commits(t *testing.T) {
	test := func(input, expected string) {
		rctx := markup.NewTestRenderContext(markup.TestAppURL, localMetas).WithRelativePath("a.md")
		buffer, err := markup.RenderString(rctx, input)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(buffer))
	}

	sha := "65f1bf27bc3bf70f64657658635e66094edbcb4d"
	repo := markup.TestAppURL + testRepoOwnerName + "/" + testRepoName + "/"
	commit := util.URLJoin(repo, "commit", sha)
	commitPath := "/user13/repo11/commit/" + sha
	tree := util.URLJoin(repo, "tree", sha, "src")

	file := util.URLJoin(repo, "commit", sha, "example.txt")
	fileWithExtra := file + ":"
	fileWithHash := file + "#L2"
	fileWithHasExtra := file + "#L2:"
	commitCompare := util.URLJoin(repo, "compare", sha+"..."+sha)
	commitCompareWithHash := commitCompare + "#L2"

	test(sha, `<p><a href="`+commitPath+`" rel="nofollow"><code>65f1bf27bc</code></a></p>`)
	test(sha[:7], `<p><a href="`+commitPath[:len(commitPath)-(40-7)]+`" rel="nofollow"><code>65f1bf2</code></a></p>`)
	test(sha[:39], `<p><a href="`+commitPath[:len(commitPath)-(40-39)]+`" rel="nofollow"><code>65f1bf27bc</code></a></p>`)
	test(commit, `<p><a href="`+commit+`" rel="nofollow"><code>65f1bf27bc</code></a></p>`)
	test(tree, `<p><a href="`+tree+`" rel="nofollow"><code>65f1bf27bc/src</code></a></p>`)

	test(file, `<p><a href="`+file+`" rel="nofollow"><code>65f1bf27bc/example.txt</code></a></p>`)
	test(fileWithExtra, `<p><a href="`+file+`" rel="nofollow"><code>65f1bf27bc/example.txt</code></a>:</p>`)
	test(fileWithHash, `<p><a href="`+fileWithHash+`" rel="nofollow"><code>65f1bf27bc/example.txt (L2)</code></a></p>`)
	test(fileWithHasExtra, `<p><a href="`+fileWithHash+`" rel="nofollow"><code>65f1bf27bc/example.txt (L2)</code></a>:</p>`)
	test(commitCompare, `<p><a href="`+commitCompare+`" rel="nofollow"><code>65f1bf27bc...65f1bf27bc</code></a></p>`)
	test(commitCompareWithHash, `<p><a href="`+commitCompareWithHash+`" rel="nofollow"><code>65f1bf27bc...65f1bf27bc (L2)</code></a></p>`)

	test("commit "+sha, `<p>commit <a href="`+commitPath+`" rel="nofollow"><code>65f1bf27bc</code></a></p>`)
	test("/home/kmup/"+sha, "<p>/home/kmup/"+sha+"</p>")
	test("deadbeef", `<p>deadbeef</p>`)
	test("d27ace93", `<p>d27ace93</p>`)
	test(sha[:14]+".x", `<p>`+sha[:14]+`.x</p>`)

	expected14 := `<a href="` + commitPath[:len(commitPath)-(40-14)] + `" rel="nofollow"><code>` + sha[:10] + `</code></a>`
	test(sha[:14]+".", `<p>`+expected14+`.</p>`)
	test(sha[:14]+",", `<p>`+expected14+`,</p>`)
	test("["+sha[:14]+"]", `<p>[`+expected14+`]</p>`)
}

func TestRender_CrossReferences(t *testing.T) {
	defer testModule.MockVariableValue(&markup.RenderBehaviorForTesting.DisableAdditionalAttributes, true)()
	test := func(input, expected string) {
		rctx := markup.NewTestRenderContext(markup.TestAppURL, localMetas).WithRelativePath("a.md")
		buffer, err := markup.RenderString(rctx, input)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(buffer))
	}

	test(
		"test-owner/test-repo#12345",
		`<p><a href="/test-owner/test-repo/issues/12345" class="ref-issue" rel="nofollow">test-owner/test-repo#12345</a></p>`)
	test(
		"kumose/kmup#12345",
		`<p><a href="/kumose/kmup/issues/12345" class="ref-issue" rel="nofollow">kumose/kmup#12345</a></p>`)
	test(
		"/home/kmup/kumose/kmup#12345",
		`<p>/home/kmup/kumose/kmup#12345</p>`)
	test(
		util.URLJoin(markup.TestAppURL, "gokmup", "kmup", "issues", "12345"),
		`<p><a href="`+util.URLJoin(markup.TestAppURL, "gokmup", "kmup", "issues", "12345")+`" class="ref-issue" rel="nofollow">gokmup/kmup#12345</a></p>`)
	test(
		util.URLJoin(markup.TestAppURL, "go-kmup", "kmup", "issues", "12345"),
		`<p><a href="`+util.URLJoin(markup.TestAppURL, "go-kmup", "kmup", "issues", "12345")+`" class="ref-issue" rel="nofollow">kumose/kmup#12345</a></p>`)
	test(
		util.URLJoin(markup.TestAppURL, "gokmup", "some-repo-name", "issues", "12345"),
		`<p><a href="`+util.URLJoin(markup.TestAppURL, "gokmup", "some-repo-name", "issues", "12345")+`" class="ref-issue" rel="nofollow">gokmup/some-repo-name#12345</a></p>`)

	inputURL := "https://host/a/b/commit/0123456789012345678901234567890123456789/foo.txt?a=b#L2-L3"
	test(
		inputURL,
		`<p><a href="`+inputURL+`" rel="nofollow"><code>0123456789/foo.txt (L2-L3)</code></a></p>`)
}

func TestRender_links(t *testing.T) {
	setting.AppURL = markup.TestAppURL
	defer testModule.MockVariableValue(&markup.RenderBehaviorForTesting.DisableAdditionalAttributes, true)()
	test := func(input, expected string) {
		buffer, err := markup.RenderString(markup.NewTestRenderContext().WithRelativePath("a.md"), input)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(buffer))
	}

	oldCustomURLSchemes := setting.Markdown.CustomURLSchemes
	markup.ResetDefaultSanitizerForTesting()
	defer func() {
		setting.Markdown.CustomURLSchemes = oldCustomURLSchemes
		markup.ResetDefaultSanitizerForTesting()
		markup.CustomLinkURLSchemes(oldCustomURLSchemes)
	}()
	setting.Markdown.CustomURLSchemes = []string{"ftp", "magnet"}
	markup.CustomLinkURLSchemes(setting.Markdown.CustomURLSchemes)

	// Text that should be turned into URL
	test(
		"https://www.example.com",
		`<p><a href="https://www.example.com" rel="nofollow">https://www.example.com</a></p>`)
	test(
		"http://www.example.com",
		`<p><a href="http://www.example.com" rel="nofollow">http://www.example.com</a></p>`)
	test(
		"https://example.com",
		`<p><a href="https://example.com" rel="nofollow">https://example.com</a></p>`)
	test(
		"http://example.com",
		`<p><a href="http://example.com" rel="nofollow">http://example.com</a></p>`)
	test(
		"http://foo.com/blah_blah",
		`<p><a href="http://foo.com/blah_blah" rel="nofollow">http://foo.com/blah_blah</a></p>`)
	test(
		"http://foo.com/blah_blah/",
		`<p><a href="http://foo.com/blah_blah/" rel="nofollow">http://foo.com/blah_blah/</a></p>`)
	test(
		"http://www.example.com/wpstyle/?p=364",
		`<p><a href="http://www.example.com/wpstyle/?p=364" rel="nofollow">http://www.example.com/wpstyle/?p=364</a></p>`)
	test(
		"https://www.example.com/foo/?bar=baz&inga=42&quux",
		`<p><a href="https://www.example.com/foo/?bar=baz&amp;inga=42&amp;quux" rel="nofollow">https://www.example.com/foo/?bar=baz&amp;inga=42&amp;quux</a></p>`)
	test(
		"http://142.42.1.1/",
		`<p><a href="http://142.42.1.1/" rel="nofollow">http://142.42.1.1/</a></p>`)
	test(
		"https://github.com/kumose/kmup/?p=aaa/bbb.html#ccc-ddd",
		`<p><a href="https://github.com/kumose/kmup/?p=aaa/bbb.html#ccc-ddd" rel="nofollow">https://github.com/kumose/kmup/?p=aaa/bbb.html#ccc-ddd</a></p>`)
	test(
		"https://en.wikipedia.org/wiki/URL_(disambiguation)",
		`<p><a href="https://en.wikipedia.org/wiki/URL_(disambiguation)" rel="nofollow">https://en.wikipedia.org/wiki/URL_(disambiguation)</a></p>`)
	test(
		"https://foo_bar.example.com/",
		`<p><a href="https://foo_bar.example.com/" rel="nofollow">https://foo_bar.example.com/</a></p>`)
	test(
		"https://stackoverflow.com/questions/2896191/what-is-go-used-fore",
		`<p><a href="https://stackoverflow.com/questions/2896191/what-is-go-used-fore" rel="nofollow">https://stackoverflow.com/questions/2896191/what-is-go-used-fore</a></p>`)
	test(
		"https://username:password@kmup.com",
		`<p><a href="https://username:password@kmup.com" rel="nofollow">https://username:password@kmup.com</a></p>`)
	test(
		"ftp://kmup.com/file.txt",
		`<p><a href="ftp://kmup.com/file.txt" rel="nofollow">ftp://kmup.com/file.txt</a></p>`)
	test(
		"magnet:?xt=urn:btih:5dee65101db281ac9c46344cd6b175cdcadabcde&dn=download",
		`<p><a href="magnet:?xt=urn:btih:5dee65101db281ac9c46344cd6b175cdcadabcde&amp;dn=download" rel="nofollow">magnet:?xt=urn:btih:5dee65101db281ac9c46344cd6b175cdcadabcde&amp;dn=download</a></p>`)
	test(
		`[link](https://example.com)`,
		`<p><a href="https://example.com" rel="nofollow">link</a></p>`)
	test(
		`[link](mailto:test@example.com)`,
		`<p><a href="mailto:test@example.com" rel="nofollow">link</a></p>`)
	test(
		`[link](javascript:xss)`,
		`<p>link</p>`)

	// Test that should *not* be turned into URL
	test(
		"www.example.com",
		`<p>www.example.com</p>`)
	test(
		"example.com",
		`<p>example.com</p>`)
	test(
		"test.example.com",
		`<p>test.example.com</p>`)
	test(
		"http://",
		`<p>http://</p>`)
	test(
		"https://",
		`<p>https://</p>`)
	test(
		"://",
		`<p>://</p>`)
	test(
		"www",
		`<p>www</p>`)
	test(
		"ftps://kmup.com",
		`<p>ftps://kmup.com</p>`)

	t.Run("LinkEllipsis", func(t *testing.T) {
		input := util.EllipsisDisplayString("http://10.1.2.3", 12)
		assert.Equal(t, "http://10…", input)
		test(input, "<p>http://10…</p>")

		input = util.EllipsisDisplayString("http://10.1.2.3", 13)
		assert.Equal(t, "http://10.…", input)
		test(input, "<p>http://10.…</p>")
	})
}

func TestRender_email(t *testing.T) {
	setting.AppURL = markup.TestAppURL
	defer testModule.MockVariableValue(&markup.RenderBehaviorForTesting.DisableAdditionalAttributes, true)()
	test := func(input, expected string) {
		res, err := markup.RenderString(markup.NewTestRenderContext().WithRelativePath("a.md"), input)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(res), "input: %s", input)
	}

	// Text that should be turned into email link
	test(
		"info@kmup.com",
		`<p><a href="mailto:info@kmup.com" rel="nofollow">info@kmup.com</a></p>`)
	test(
		"(info@kmup.com)",
		`<p>(<a href="mailto:info@kmup.com" rel="nofollow">info@kmup.com</a>)</p>`)
	test(
		"[info@kmup.com]",
		`<p>[<a href="mailto:info@kmup.com" rel="nofollow">info@kmup.com</a>]</p>`)
	test(
		"info@kmup.com.",
		`<p><a href="mailto:info@kmup.com" rel="nofollow">info@kmup.com</a>.</p>`)
	test(
		"firstname+lastname@kmup.com",
		`<p><a href="mailto:firstname+lastname@kmup.com" rel="nofollow">firstname+lastname@kmup.com</a></p>`)
	test(
		"send email to info@kmup.co.uk.",
		`<p>send email to <a href="mailto:info@kmup.co.uk" rel="nofollow">info@kmup.co.uk</a>.</p>`)

	test(
		`j.doe@example.com,
	j.doe@example.com.
	j.doe@example.com;
	j.doe@example.com?
	j.doe@example.com!`,
		`<p><a href="mailto:j.doe@example.com" rel="nofollow">j.doe@example.com</a>,
<a href="mailto:j.doe@example.com" rel="nofollow">j.doe@example.com</a>.
<a href="mailto:j.doe@example.com" rel="nofollow">j.doe@example.com</a>;
<a href="mailto:j.doe@example.com" rel="nofollow">j.doe@example.com</a>?
<a href="mailto:j.doe@example.com" rel="nofollow">j.doe@example.com</a>!</p>`)

	// match GitHub behavior
	test("email@domain@domain.com", `<p>email@<a href="mailto:domain@domain.com" rel="nofollow">domain@domain.com</a></p>`)

	// match GitHub behavior
	test(`"info@kmup.com"`, `<p>&#34;<a href="mailto:info@kmup.com" rel="nofollow">info@kmup.com</a>&#34;</p>`)

	// Test that should *not* be turned into email links
	test(
		"/home/kmup/mailstore/info@kmup/com",
		`<p>/home/kmup/mailstore/info@kmup/com</p>`)
	test(
		"git@try.kmup.io:kumose/kmup.git",
		`<p>git@try.kmup.io:kumose/kmup.git</p>`)
	test(
		"https://foo:bar@kmup.io",
		`<p><a href="https://foo:bar@kmup.io" rel="nofollow">https://foo:bar@kmup.io</a></p>`)
	test(
		"kmup@3",
		`<p>kmup@3</p>`)
	test(
		"kmup@gmail.c",
		`<p>kmup@gmail.c</p>`)
	test(
		"email@domain..com",
		`<p>email@domain..com</p>`)

	cases := []struct {
		input, expected string
	}{
		// match GitHub behavior
		{"?a@d.zz", `<p>?<a href="mailto:a@d.zz" rel="nofollow">a@d.zz</a></p>`},
		{"*a@d.zz", `<p>*<a href="mailto:a@d.zz" rel="nofollow">a@d.zz</a></p>`},
		{"~a@d.zz", `<p>~<a href="mailto:a@d.zz" rel="nofollow">a@d.zz</a></p>`},

		// the following cases don't match GitHub behavior, but they are valid email addresses ...
		// maybe we should reduce the candidate characters for the "name" part in the future
		{"a*a@d.zz", `<p><a href="mailto:a*a@d.zz" rel="nofollow">a*a@d.zz</a></p>`},
		{"a~a@d.zz", `<p><a href="mailto:a~a@d.zz" rel="nofollow">a~a@d.zz</a></p>`},
	}
	for _, c := range cases {
		test(c.input, c.expected)
	}
}

func TestRender_emoji(t *testing.T) {
	setting.AppURL = markup.TestAppURL
	setting.StaticURLPrefix = markup.TestAppURL

	test := func(input, expected string) {
		expected = strings.ReplaceAll(expected, "&", "&amp;")
		buffer, err := markup.RenderString(markup.NewTestRenderContext().WithRelativePath("a.md"), input)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(buffer))
	}

	// Make sure we can successfully match every emoji in our dataset with regex
	for i := range emoji.GemojiData {
		test(
			emoji.GemojiData[i].Emoji,
			`<p><span class="emoji" aria-label="`+emoji.GemojiData[i].Description+`">`+emoji.GemojiData[i].Emoji+`</span></p>`)
	}
	for i := range emoji.GemojiData {
		test(
			":"+emoji.GemojiData[i].Aliases[0]+":",
			`<p><span class="emoji" aria-label="`+emoji.GemojiData[i].Description+`">`+emoji.GemojiData[i].Emoji+`</span></p>`)
	}

	// Text that should be turned into or recognized as emoji
	test(
		":kmup:",
		`<p><span class="emoji" aria-label="kmup"><img alt=":kmup:" src="`+setting.StaticURLPrefix+`/assets/img/emoji/kmup.png"/></span></p>`)
	test(
		":custom-emoji:",
		`<p>:custom-emoji:</p>`)
	setting.UI.CustomEmojisMap["custom-emoji"] = ":custom-emoji:"
	test(
		":custom-emoji:",
		`<p><span class="emoji" aria-label="custom-emoji"><img alt=":custom-emoji:" src="`+setting.StaticURLPrefix+`/assets/img/emoji/custom-emoji.png"/></span></p>`)
	test(
		"这是字符:1::+1: some🐊 \U0001f44d:custom-emoji: :kmup:",
		`<p>这是字符:1:<span class="emoji" aria-label="thumbs up">👍</span> some<span class="emoji" aria-label="crocodile">🐊</span> `+
			`<span class="emoji" aria-label="thumbs up">👍</span><span class="emoji" aria-label="custom-emoji"><img alt=":custom-emoji:" src="`+setting.StaticURLPrefix+`/assets/img/emoji/custom-emoji.png"/></span> `+
			`<span class="emoji" aria-label="kmup"><img alt=":kmup:" src="`+setting.StaticURLPrefix+`/assets/img/emoji/kmup.png"/></span></p>`)
	test(
		"Some text with 😄 in the middle",
		`<p>Some text with <span class="emoji" aria-label="grinning face with smiling eyes">😄</span> in the middle</p>`)
	test(
		"Some text with :smile: in the middle",
		`<p>Some text with <span class="emoji" aria-label="grinning face with smiling eyes">😄</span> in the middle</p>`)
	test(
		"Some text with 😄😄 2 emoji next to each other",
		`<p>Some text with <span class="emoji" aria-label="grinning face with smiling eyes">😄</span><span class="emoji" aria-label="grinning face with smiling eyes">😄</span> 2 emoji next to each other</p>`)
	test(
		"😎🤪🔐🤑❓",
		`<p><span class="emoji" aria-label="smiling face with sunglasses">😎</span><span class="emoji" aria-label="zany face">🤪</span><span class="emoji" aria-label="locked with key">🔐</span><span class="emoji" aria-label="money-mouth face">🤑</span><span class="emoji" aria-label="red question mark">❓</span></p>`)

	// should match nothing
	test(":100:200", `<p>:100:200</p>`)
	test("std::thread::something", `<p>std::thread::something</p>`)
	test(":not exist:", `<p>:not exist:</p>`)
}

func TestRender_ShortLinks(t *testing.T) {
	setting.AppURL = markup.TestAppURL
	tree := util.URLJoin(markup.TestRepoURL, "src", "master")

	test := func(input, expected string) {
		buffer, err := markdown.RenderString(markup.NewTestRenderContext(tree), input)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(buffer)))
	}

	url := util.URLJoin(tree, "Link")
	otherURL := util.URLJoin(tree, "Other-Link")
	encodedURL := util.URLJoin(tree, "Link%3F")
	imgurl := util.URLJoin(tree, "Link.jpg")
	otherImgurl := util.URLJoin(tree, "Link+Other.jpg")
	encodedImgurl := util.URLJoin(tree, "Link+%23.jpg")
	notencodedImgurl := util.URLJoin(tree, "some", "path", "Link+#.jpg")
	renderableFileURL := util.URLJoin(tree, "markdown_file.md")
	unrenderableFileURL := util.URLJoin(tree, "file.zip")
	favicon := "http://google.com/favicon.ico"

	test(
		"[[Link]]",
		`<p><a href="`+url+`" rel="nofollow">Link</a></p>`,
	)
	test(
		"[[Link.-]]",
		`<p><a href="http://localhost:3326/test-owner/test-repo/src/master/Link.-" rel="nofollow">Link.-</a></p>`,
	)
	test(
		"[[Link.jpg]]",
		`<p><a href="`+imgurl+`" rel="nofollow"><img src="`+imgurl+`" title="Link.jpg" alt="Link.jpg"/></a></p>`,
	)
	test(
		"[["+favicon+"]]",
		`<p><a href="`+favicon+`" rel="nofollow"><img src="`+favicon+`" title="favicon.ico" alt="`+favicon+`"/></a></p>`,
	)
	test(
		"[[Name|Link]]",
		`<p><a href="`+url+`" rel="nofollow">Name</a></p>`,
	)
	test(
		"[[Name|Link.jpg]]",
		`<p><a href="`+imgurl+`" rel="nofollow"><img src="`+imgurl+`" title="Name" alt="Name"/></a></p>`,
	)
	test(
		"[[Name|Link.jpg|alt=AltName]]",
		`<p><a href="`+imgurl+`" rel="nofollow"><img src="`+imgurl+`" title="AltName" alt="AltName"/></a></p>`,
	)
	test(
		"[[Name|Link.jpg|title=Title]]",
		`<p><a href="`+imgurl+`" rel="nofollow"><img src="`+imgurl+`" title="Title" alt="Title"/></a></p>`,
	)
	test(
		"[[Name|Link.jpg|alt=AltName|title=Title]]",
		`<p><a href="`+imgurl+`" rel="nofollow"><img src="`+imgurl+`" title="Title" alt="AltName"/></a></p>`,
	)
	test(
		"[[Name|Link.jpg|alt=\"AltName\"|title='Title']]",
		`<p><a href="`+imgurl+`" rel="nofollow"><img src="`+imgurl+`" title="Title" alt="AltName"/></a></p>`,
	)
	test(
		"[[Name|Link Other.jpg|alt=\"AltName\"|title='Title']]",
		`<p><a href="`+otherImgurl+`" rel="nofollow"><img src="`+otherImgurl+`" title="Title" alt="AltName"/></a></p>`,
	)
	test(
		"[[Link]] [[Other Link]]",
		`<p><a href="`+url+`" rel="nofollow">Link</a> <a href="`+otherURL+`" rel="nofollow">Other Link</a></p>`,
	)
	test(
		"[[Link?]]",
		`<p><a href="`+encodedURL+`" rel="nofollow">Link?</a></p>`,
	)
	test(
		"[[Link]] [[Other Link]] [[Link?]]",
		`<p><a href="`+url+`" rel="nofollow">Link</a> <a href="`+otherURL+`" rel="nofollow">Other Link</a> <a href="`+encodedURL+`" rel="nofollow">Link?</a></p>`,
	)
	test(
		"[[markdown_file.md]]",
		`<p><a href="`+renderableFileURL+`" rel="nofollow">markdown_file.md</a></p>`,
	)
	test(
		"[[file.zip]]",
		`<p><a href="`+unrenderableFileURL+`" rel="nofollow">file.zip</a></p>`,
	)
	test(
		"[[Link #.jpg]]",
		`<p><a href="`+encodedImgurl+`" rel="nofollow"><img src="`+encodedImgurl+`" title="Link #.jpg" alt="Link #.jpg"/></a></p>`,
	)
	test(
		"[[Name|Link #.jpg|alt=\"AltName\"|title='Title']]",
		`<p><a href="`+encodedImgurl+`" rel="nofollow"><img src="`+encodedImgurl+`" title="Title" alt="AltName"/></a></p>`,
	)
	test(
		"[[some/path/Link #.jpg]]",
		`<p><a href="`+notencodedImgurl+`" rel="nofollow"><img src="`+notencodedImgurl+`" title="Link #.jpg" alt="some/path/Link #.jpg"/></a></p>`,
	)
	test(
		"<p><a href=\"https://example.org\">[[foobar]]</a></p>",
		`<p><a href="https://example.org" rel="nofollow">[[foobar]]</a></p>`,
	)
}

func Test_ParseClusterFuzz(t *testing.T) {
	setting.AppURL = markup.TestAppURL

	localMetas := map[string]string{"user": "go-kmup", "repo": "kmup"}

	data := "<A><maTH><tr><MN><bodY ÿ><temPlate></template><tH><tr></A><tH><d<bodY "

	var res strings.Builder
	err := markup.PostProcessDefault(markup.NewTestRenderContext(localMetas), strings.NewReader(data), &res)
	assert.NoError(t, err)
	assert.NotContains(t, res.String(), "<html")

	data = "<!DOCTYPE html>\n<A><maTH><tr><MN><bodY ÿ><temPlate></template><tH><tr></A><tH><d<bodY "

	res.Reset()
	err = markup.PostProcessDefault(markup.NewTestRenderContext(localMetas), strings.NewReader(data), &res)

	assert.NoError(t, err)
	assert.NotContains(t, res.String(), "<html")
}

func TestPostProcess(t *testing.T) {
	setting.StaticURLPrefix = markup.TestAppURL // can't run standalone
	defer testModule.MockVariableValue(&markup.RenderBehaviorForTesting.DisableAdditionalAttributes, true)()

	test := func(input, expected string) {
		var res strings.Builder
		err := markup.PostProcessDefault(markup.NewTestRenderContext(markup.TestAppURL, map[string]string{"user": "go-kmup", "repo": "kmup"}), strings.NewReader(input), &res)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(res.String()))
	}

	// Issue index shouldn't be post-processing in a document.
	test(
		"#1",
		"#1")

	// But cross-referenced issue index should work.
	test(
		"kumose/kmup#12345",
		`<a href="/kumose/kmup/issues/12345" class="ref-issue">kumose/kmup#12345</a>`)

	// Test that other post-processing still works.
	test(
		":kmup:",
		`<span class="emoji" aria-label="kmup"><img alt=":kmup:" src="`+setting.StaticURLPrefix+`/assets/img/emoji/kmup.png"/></span>`)
	test(
		"Some text with 😄 in the middle",
		`Some text with <span class="emoji" aria-label="grinning face with smiling eyes">😄</span> in the middle`)
	test("http://localhost:3326/person/repo/issues/4#issuecomment-1234",
		`<a href="http://localhost:3326/person/repo/issues/4#issuecomment-1234" class="ref-issue">person/repo#4 (comment)</a>`)

	// special tags, GitHub's behavior, and for unclosed tags, output as text content as much as possible
	test("<script>a", `&lt;script&gt;a`)
	test("<script>a</script>", `&lt;script&gt;a&lt;/script&gt;`)
	test("<STYLE>a", `&lt;STYLE&gt;a`)
	test("<style>a</STYLE>", `&lt;style&gt;a&lt;/STYLE&gt;`)

	// other special tags, our special behavior
	test("<?php\nfoo", "&lt;?php\nfoo")
	test("<%asp\nfoo", "&lt;%asp\nfoo")
}

func TestIssue16020(t *testing.T) {
	setting.AppURL = markup.TestAppURL

	localMetas := map[string]string{
		"user": "go-kmup",
		"repo": "kmup",
	}

	data := `<img src="data:image/png;base64,i//V"/>`

	var res strings.Builder
	err := markup.PostProcessDefault(markup.NewTestRenderContext(localMetas), strings.NewReader(data), &res)
	assert.NoError(t, err)
	assert.Equal(t, data, res.String())
}

func BenchmarkEmojiPostprocess(b *testing.B) {
	data := "🥰 "
	for len(data) < 1<<16 {
		data += data
	}
	b.ResetTimer()
	for b.Loop() {
		var res strings.Builder
		err := markup.PostProcessDefault(markup.NewTestRenderContext(localMetas), strings.NewReader(data), &res)
		assert.NoError(b, err)
	}
}

func TestFuzz(t *testing.T) {
	s := "t/l/issues/8#/../../a"
	renderContext := markup.NewTestRenderContext()
	err := markup.PostProcessDefault(renderContext, strings.NewReader(s), io.Discard)
	assert.NoError(t, err)
}

func TestIssue18471(t *testing.T) {
	data := `http://domain/org/repo/compare/783b039...da951ce`

	var res strings.Builder
	err := markup.PostProcessDefault(markup.NewTestRenderContext(localMetas), strings.NewReader(data), &res)

	assert.NoError(t, err)
	assert.Equal(t, `<a href="http://domain/org/repo/compare/783b039...da951ce" class="compare"><code>783b039...da951ce</code></a>`, res.String())
}

func TestIsFullURL(t *testing.T) {
	assert.True(t, markup.IsFullURLString("https://example.com"))
	assert.True(t, markup.IsFullURLString("mailto:test@example.com"))
	assert.True(t, markup.IsFullURLString("data:image/11111"))
	assert.False(t, markup.IsFullURLString("/foo:bar"))
}
