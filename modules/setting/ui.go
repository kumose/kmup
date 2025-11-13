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
	"time"

	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/log"
)

// UI settings
var UI = struct {
	ExplorePagingNum        int
	SitemapPagingNum        int
	IssuePagingNum          int
	RepoSearchPagingNum     int
	MembersPagingNum        int
	FeedMaxCommitNum        int
	FeedPagingNum           int
	PackagesPagingNum       int
	GraphMaxCommitNum       int
	CodeCommentLines        int
	ReactionMaxUserNum      int
	MaxDisplayFileSize      int64
	ShowUserEmail           bool
	DefaultShowFullName     bool
	DefaultTheme            string
	Themes                  []string
	FileIconTheme           string
	Reactions               []string
	ReactionsLookup         container.Set[string] `ini:"-"`
	CustomEmojis            []string
	CustomEmojisMap         map[string]string `ini:"-"`
	EnabledEmojis           []string
	EnabledEmojisSet        container.Set[string] `ini:"-"`
	SearchRepoDescription   bool
	OnlyShowRelevantRepos   bool
	ExploreDefaultSort      string `ini:"EXPLORE_PAGING_DEFAULT_SORT"`
	PreferredTimestampTense string

	AmbiguousUnicodeDetection bool

	Notification struct {
		MinTimeout            time.Duration
		TimeoutStep           time.Duration
		MaxTimeout            time.Duration
		EventSourceUpdateTime time.Duration
	} `ini:"ui.notification"`

	SVG struct {
		Enabled bool `ini:"ENABLE_RENDER"`
	} `ini:"ui.svg"`

	CSV struct {
		MaxFileSize int64
		MaxRows     int
	} `ini:"ui.csv"`

	Admin struct {
		UserPagingNum   int
		RepoPagingNum   int
		NoticePagingNum int
		OrgPagingNum    int
	} `ini:"ui.admin"`
	User struct {
		RepoPagingNum int
		OrgPagingNum  int
	} `ini:"ui.user"`
	Meta struct {
		Author      string
		Description string
		Keywords    string
	} `ini:"ui.meta"`
}{
	ExplorePagingNum:        20,
	SitemapPagingNum:        20,
	IssuePagingNum:          20,
	RepoSearchPagingNum:     20,
	MembersPagingNum:        20,
	FeedMaxCommitNum:        5,
	FeedPagingNum:           20,
	PackagesPagingNum:       20,
	GraphMaxCommitNum:       100,
	CodeCommentLines:        4,
	ReactionMaxUserNum:      10,
	MaxDisplayFileSize:      8388608,
	DefaultTheme:            `kmup-auto`,
	FileIconTheme:           `material`,
	Reactions:               []string{`+1`, `-1`, `laugh`, `hooray`, `confused`, `heart`, `rocket`, `eyes`},
	CustomEmojis:            []string{`git`, `kmup`, `codeberg`, `gitlab`, `github`, `gogs`},
	CustomEmojisMap:         map[string]string{"git": ":git:", "kmup": ":kmup:", "codeberg": ":codeberg:", "gitlab": ":gitlab:", "github": ":github:", "gogs": ":gogs:"},
	ExploreDefaultSort:      "recentupdate",
	PreferredTimestampTense: "mixed",

	AmbiguousUnicodeDetection: true,

	Notification: struct {
		MinTimeout            time.Duration
		TimeoutStep           time.Duration
		MaxTimeout            time.Duration
		EventSourceUpdateTime time.Duration
	}{
		MinTimeout:            10 * time.Second,
		TimeoutStep:           10 * time.Second,
		MaxTimeout:            60 * time.Second,
		EventSourceUpdateTime: 10 * time.Second,
	},
	SVG: struct {
		Enabled bool `ini:"ENABLE_RENDER"`
	}{
		Enabled: true,
	},
	CSV: struct {
		MaxFileSize int64
		MaxRows     int
	}{
		MaxFileSize: 524288,
		MaxRows:     2500,
	},
	Admin: struct {
		UserPagingNum   int
		RepoPagingNum   int
		NoticePagingNum int
		OrgPagingNum    int
	}{
		UserPagingNum:   50,
		RepoPagingNum:   50,
		NoticePagingNum: 25,
		OrgPagingNum:    50,
	},
	User: struct {
		RepoPagingNum int
		OrgPagingNum  int
	}{
		RepoPagingNum: 15,
		OrgPagingNum:  15,
	},
	Meta: struct {
		Author      string
		Description string
		Keywords    string
	}{
		Author:      "Kmup - Working with a cup of tea",
		Description: "Kmup (Working with a cup of tea) is a painless self-hosted Git service written in Go",
		Keywords:    "go,git,self-hosted,kmup",
	},
}

func loadUIFrom(rootCfg ConfigProvider) {
	mustMapSetting(rootCfg, "ui", &UI)
	sec := rootCfg.Section("ui")
	UI.ShowUserEmail = sec.Key("SHOW_USER_EMAIL").MustBool(true)
	UI.DefaultShowFullName = sec.Key("DEFAULT_SHOW_FULL_NAME").MustBool(false)
	UI.SearchRepoDescription = sec.Key("SEARCH_REPO_DESCRIPTION").MustBool(true)

	if UI.PreferredTimestampTense != "mixed" && UI.PreferredTimestampTense != "absolute" {
		log.Fatal("ui.PREFERRED_TIMESTAMP_TENSE must be either 'mixed' or 'absolute'")
	}

	// OnlyShowRelevantRepos=false is important for many private/enterprise instances,
	// because many private repositories do not have "description/topic", users just want to search by their names.
	UI.OnlyShowRelevantRepos = sec.Key("ONLY_SHOW_RELEVANT_REPOS").MustBool(false)

	UI.ReactionsLookup = make(container.Set[string])
	for _, reaction := range UI.Reactions {
		UI.ReactionsLookup.Add(reaction)
	}
	UI.CustomEmojisMap = make(map[string]string)
	for _, emoji := range UI.CustomEmojis {
		UI.CustomEmojisMap[emoji] = ":" + emoji + ":"
	}
	UI.EnabledEmojisSet = container.SetOf(UI.EnabledEmojis...)
}
