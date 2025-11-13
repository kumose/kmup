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

package label

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kumose/kmup/modules/options"

	"gopkg.in/yaml.v3"
)

type labelFile struct {
	Labels []*Label `yaml:"labels"`
}

// ErrTemplateLoad represents a "ErrTemplateLoad" kind of error.
type ErrTemplateLoad struct {
	TemplateFile  string
	OriginalError error
}

// IsErrTemplateLoad checks if an error is a ErrTemplateLoad.
func IsErrTemplateLoad(err error) bool {
	_, ok := err.(ErrTemplateLoad)
	return ok
}

func (err ErrTemplateLoad) Error() string {
	return fmt.Sprintf("failed to load label template file %q: %v", err.TemplateFile, err.OriginalError)
}

// LoadTemplateFile loads the label template file by given file name, returns a slice of Label structs.
func LoadTemplateFile(fileName string) ([]*Label, error) {
	data, err := options.Labels(fileName)
	if err != nil {
		return nil, ErrTemplateLoad{fileName, fmt.Errorf("LoadTemplateFile: %w", err)}
	}

	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		return parseYamlFormat(fileName, data)
	}
	return parseLegacyFormat(fileName, data)
}

func parseYamlFormat(fileName string, data []byte) ([]*Label, error) {
	lf := &labelFile{}

	if err := yaml.Unmarshal(data, lf); err != nil {
		return nil, err
	}

	// Validate label data and fix colors
	for _, l := range lf.Labels {
		l.Color = strings.TrimSpace(l.Color)
		if len(l.Name) == 0 || len(l.Color) == 0 {
			return nil, ErrTemplateLoad{fileName, errors.New("label name and color are required fields")}
		}
		color, err := NormalizeColor(l.Color)
		if err != nil {
			return nil, ErrTemplateLoad{fileName, fmt.Errorf("bad HTML color code '%s' in label: %s", l.Color, l.Name)}
		}
		l.Color = color
	}

	return lf.Labels, nil
}

func parseLegacyFormat(fileName string, data []byte) ([]*Label, error) {
	lines := strings.Split(string(data), "\n")
	list := make([]*Label, 0, len(lines))
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		if len(line) == 0 {
			continue
		}

		parts, description, _ := strings.Cut(line, ";")

		color, labelName, ok := strings.Cut(parts, " ")
		if !ok {
			return nil, ErrTemplateLoad{fileName, fmt.Errorf("line is malformed: %s", line)}
		}

		color, err := NormalizeColor(color)
		if err != nil {
			return nil, ErrTemplateLoad{fileName, fmt.Errorf("bad HTML color code '%s' in line: %s", color, line)}
		}

		list = append(list, &Label{
			Name:        strings.TrimSpace(labelName),
			Color:       color,
			Description: strings.TrimSpace(description),
		})
	}

	return list, nil
}

// LoadTemplateDescription loads the labels from a template file, returns a description string by joining each Label.Name with comma
func LoadTemplateDescription(fileName string) (string, error) {
	var buf strings.Builder
	list, err := LoadTemplateFile(fileName)
	if err != nil {
		return "", err
	}

	for i := range list {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(list[i].Name)
	}
	return buf.String(), nil
}
