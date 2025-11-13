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

package issue

import (
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/issue/template"
	"github.com/kumose/kmup/modules/log"
	api "github.com/kumose/kmup/modules/structs"

	"gopkg.in/yaml.v3"
)

// templateDirCandidates issue templates directory
var templateDirCandidates = []string{
	"ISSUE_TEMPLATE",
	"issue_template",
	".kmup/ISSUE_TEMPLATE",
	".kmup/issue_template",
	".github/ISSUE_TEMPLATE",
	".github/issue_template",
	".gitlab/ISSUE_TEMPLATE",
	".gitlab/issue_template",
}

var templateConfigCandidates = []string{
	".kmup/ISSUE_TEMPLATE/config",
	".kmup/issue_template/config",
	".github/ISSUE_TEMPLATE/config",
	".github/issue_template/config",
}

func GetDefaultTemplateConfig() api.IssueConfig {
	return api.IssueConfig{
		BlankIssuesEnabled: true,
		ContactLinks:       make([]api.IssueConfigContactLink, 0),
	}
}

// GetTemplateConfig loads the given issue config file.
// It never returns a nil config.
func GetTemplateConfig(gitRepo *git.Repository, path string, commit *git.Commit) (api.IssueConfig, error) {
	if gitRepo == nil {
		return GetDefaultTemplateConfig(), nil
	}

	treeEntry, err := commit.GetTreeEntryByPath(path)
	if err != nil {
		return GetDefaultTemplateConfig(), err
	}

	reader, err := treeEntry.Blob().DataAsync()
	if err != nil {
		log.Debug("DataAsync: %v", err)
		return GetDefaultTemplateConfig(), nil
	}

	defer reader.Close()

	configContent, err := io.ReadAll(reader)
	if err != nil {
		return GetDefaultTemplateConfig(), err
	}

	issueConfig := GetDefaultTemplateConfig()
	if err := yaml.Unmarshal(configContent, &issueConfig); err != nil {
		return GetDefaultTemplateConfig(), err
	}

	for pos, link := range issueConfig.ContactLinks {
		if link.Name == "" {
			return GetDefaultTemplateConfig(), fmt.Errorf("contact_link at position %d is missing name key", pos+1)
		}

		if link.URL == "" {
			return GetDefaultTemplateConfig(), fmt.Errorf("contact_link at position %d is missing url key", pos+1)
		}

		if link.About == "" {
			return GetDefaultTemplateConfig(), fmt.Errorf("contact_link at position %d is missing about key", pos+1)
		}

		_, err = url.ParseRequestURI(link.URL)
		if err != nil {
			return GetDefaultTemplateConfig(), fmt.Errorf("%s is not a valid URL", link.URL)
		}
	}

	return issueConfig, nil
}

// IsTemplateConfig returns if the given path is a issue config file.
func IsTemplateConfig(path string) bool {
	for _, configName := range templateConfigCandidates {
		if path == configName+".yaml" || path == configName+".yml" {
			return true
		}
	}
	return false
}

// ParseTemplatesFromDefaultBranch parses the issue templates in the repo's default branch,
// returns valid templates and the errors of invalid template files (the errors map is guaranteed to be non-nil).
func ParseTemplatesFromDefaultBranch(repo *repo.Repository, gitRepo *git.Repository) (ret struct {
	IssueTemplates []*api.IssueTemplate
	TemplateErrors map[string]error
},
) {
	ret.TemplateErrors = map[string]error{}
	if repo.IsEmpty {
		return ret
	}

	commit, err := gitRepo.GetBranchCommit(repo.DefaultBranch)
	if err != nil {
		return ret
	}

	for _, dirName := range templateDirCandidates {
		tree, err := commit.SubTree(dirName)
		if err != nil {
			log.Debug("get sub tree of %s: %v", dirName, err)
			continue
		}
		entries, err := tree.ListEntries()
		if err != nil {
			log.Debug("list entries in %s: %v", dirName, err)
			return ret
		}
		for _, entry := range entries {
			if !template.CouldBe(entry.Name()) {
				continue
			}
			fullName := path.Join(dirName, entry.Name())
			if it, err := template.UnmarshalFromEntry(entry, dirName); err != nil {
				ret.TemplateErrors[fullName] = err
			} else {
				if !strings.HasPrefix(it.Ref, "refs/") { // Assume that the ref intended is always a branch - for tags users should use refs/tags/<ref>
					it.Ref = git.BranchPrefix + it.Ref
				}
				ret.IssueTemplates = append(ret.IssueTemplates, it)
			}
		}
	}
	return ret
}

// GetTemplateConfigFromDefaultBranch returns the issue config for this repo.
// It never returns a nil config.
func GetTemplateConfigFromDefaultBranch(repo *repo.Repository, gitRepo *git.Repository) (api.IssueConfig, error) {
	if repo.IsEmpty {
		return GetDefaultTemplateConfig(), nil
	}

	commit, err := gitRepo.GetBranchCommit(repo.DefaultBranch)
	if err != nil {
		return GetDefaultTemplateConfig(), err
	}

	for _, configName := range templateConfigCandidates {
		if _, err := commit.GetTreeEntryByPath(configName + ".yaml"); err == nil {
			return GetTemplateConfig(gitRepo, configName+".yaml", commit)
		}

		if _, err := commit.GetTreeEntryByPath(configName + ".yml"); err == nil {
			return GetTemplateConfig(gitRepo, configName+".yml", commit)
		}
	}

	return GetDefaultTemplateConfig(), nil
}

func HasTemplatesOrContactLinks(repo *repo.Repository, gitRepo *git.Repository) bool {
	ret := ParseTemplatesFromDefaultBranch(repo, gitRepo)
	if len(ret.IssueTemplates) > 0 {
		return true
	}

	issueConfig, _ := GetTemplateConfigFromDefaultBranch(repo, gitRepo)
	return len(issueConfig.ContactLinks) > 0
}
