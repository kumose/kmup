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

package repo

import (
	"html/template"
	"net/http"
	"path"
	"strings"

	pull_model "github.com/kumose/kmup/models/pull"
	"github.com/kumose/kmup/modules/base"
	"github.com/kumose/kmup/modules/fileicon"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/gitdiff"
	files_service "github.com/kumose/kmup/services/repository/files"

	"github.com/go-enry/go-enry/v2"
)

// TreeList get all files' entries of a repository
func TreeList(ctx *context.Context) {
	tree, err := ctx.Repo.Commit.SubTree("/")
	if err != nil {
		ctx.ServerError("Repo.Commit.SubTree", err)
		return
	}

	entries, err := tree.ListEntriesRecursiveFast()
	if err != nil {
		ctx.ServerError("ListEntriesRecursiveFast", err)
		return
	}
	entries.CustomSort(base.NaturalSortCompare)

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !isExcludedEntry(entry) {
			files = append(files, entry.Name())
		}
	}
	ctx.JSON(http.StatusOK, files)
}

func isExcludedEntry(entry *git.TreeEntry) bool {
	if entry.IsDir() {
		return true
	}

	if entry.IsSubModule() {
		return true
	}

	if enry.IsVendor(entry.Name()) {
		return true
	}

	return false
}

// WebDiffFileItem is used by frontend, check the field names in frontend before changing
type WebDiffFileItem struct {
	FullName    string
	DisplayName string
	NameHash    string
	DiffStatus  string
	EntryMode   string
	IsViewed    bool
	Children    []*WebDiffFileItem
	FileIcon    template.HTML
}

// WebDiffFileTree is used by frontend, check the field names in frontend before changing
type WebDiffFileTree struct {
	TreeRoot WebDiffFileItem
}

// transformDiffTreeForWeb transforms a gitdiff.DiffTree into a WebDiffFileTree for Web UI rendering
// it also takes a map of file names to their viewed state, which is used to mark files as viewed
func transformDiffTreeForWeb(renderedIconPool *fileicon.RenderedIconPool, diffTree *gitdiff.DiffTree, filesViewedState map[string]pull_model.ViewedState) (dft WebDiffFileTree) {
	dirNodes := map[string]*WebDiffFileItem{"": &dft.TreeRoot}
	addItem := func(item *WebDiffFileItem) {
		var parentPath string
		pos := strings.LastIndexByte(item.FullName, '/')
		if pos == -1 {
			item.DisplayName = item.FullName
		} else {
			parentPath = item.FullName[:pos]
			item.DisplayName = item.FullName[pos+1:]
		}
		parentNode, parentExists := dirNodes[parentPath]
		if !parentExists {
			parentNode = &dft.TreeRoot
			fields := strings.Split(parentPath, "/")
			for idx, field := range fields {
				nodePath := strings.Join(fields[:idx+1], "/")
				node, ok := dirNodes[nodePath]
				if !ok {
					node = &WebDiffFileItem{EntryMode: "tree", DisplayName: field, FullName: nodePath}
					dirNodes[nodePath] = node
					parentNode.Children = append(parentNode.Children, node)
				}
				parentNode = node
			}
		}
		parentNode.Children = append(parentNode.Children, item)
	}

	for _, file := range diffTree.Files {
		item := &WebDiffFileItem{FullName: file.HeadPath, DiffStatus: file.Status}
		item.IsViewed = filesViewedState[item.FullName] == pull_model.Viewed
		item.NameHash = git.HashFilePathForWebUI(item.FullName)
		item.FileIcon = fileicon.RenderEntryIconHTML(renderedIconPool, &fileicon.EntryInfo{BaseName: path.Base(file.HeadPath), EntryMode: file.HeadMode})

		switch file.HeadMode {
		case git.EntryModeTree:
			item.EntryMode = "tree"
		case git.EntryModeCommit:
			item.EntryMode = "commit" // submodule
		default:
			// default to empty, and will be treated as "blob" file because there is no "symlink" support yet
		}
		addItem(item)
	}

	var mergeSingleDir func(node *WebDiffFileItem)
	mergeSingleDir = func(node *WebDiffFileItem) {
		if len(node.Children) == 1 {
			if child := node.Children[0]; child.EntryMode == "tree" {
				node.FullName = child.FullName
				node.DisplayName = node.DisplayName + "/" + child.DisplayName
				node.Children = child.Children
				mergeSingleDir(node)
			}
		}
	}
	for _, node := range dft.TreeRoot.Children {
		mergeSingleDir(node)
	}
	return dft
}

func TreeViewNodes(ctx *context.Context) {
	renderedIconPool := fileicon.NewRenderedIconPool()
	results, err := files_service.GetTreeViewNodes(ctx, ctx.Repo.RepoLink, renderedIconPool, ctx.Repo.Commit, ctx.Repo.TreePath, ctx.FormString("sub_path"))
	if err != nil {
		ctx.ServerError("GetTreeViewNodes", err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"fileTreeNodes": results, "renderedIconPool": renderedIconPool.IconSVGs})
}
