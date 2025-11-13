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

package external

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/setting"

	"github.com/kballard/go-shellquote"
)

// RegisterRenderers registers all supported third part renderers according settings
func RegisterRenderers() {
	for _, renderer := range setting.ExternalMarkupRenderers {
		if renderer.Enabled && renderer.Command != "" && len(renderer.FileExtensions) > 0 {
			markup.RegisterRenderer(&Renderer{renderer})
		}
	}
}

// Renderer implements markup.Renderer for external tools
type Renderer struct {
	*setting.MarkupRenderer
}

var (
	_ markup.PostProcessRenderer = (*Renderer)(nil)
	_ markup.ExternalRenderer    = (*Renderer)(nil)
)

// Name returns the external tool name
func (p *Renderer) Name() string {
	return p.MarkupName
}

// NeedPostProcess implements markup.Renderer
func (p *Renderer) NeedPostProcess() bool {
	return p.MarkupRenderer.NeedPostProcess
}

// Extensions returns the supported extensions of the tool
func (p *Renderer) Extensions() []string {
	return p.FileExtensions
}

// SanitizerRules implements markup.Renderer
func (p *Renderer) SanitizerRules() []setting.MarkupSanitizerRule {
	return p.MarkupSanitizerRules
}

func (p *Renderer) GetExternalRendererOptions() (ret markup.ExternalRendererOptions) {
	ret.SanitizerDisabled = p.RenderContentMode == setting.RenderContentModeNoSanitizer || p.RenderContentMode == setting.RenderContentModeIframe
	ret.DisplayInIframe = p.RenderContentMode == setting.RenderContentModeIframe
	ret.ContentSandbox = p.RenderContentSandbox
	return ret
}

func envMark(envName string) string {
	if runtime.GOOS == "windows" {
		return "%" + envName + "%"
	}
	return "$" + envName
}

// Render renders the data of the document to HTML via the external tool.
func (p *Renderer) Render(ctx *markup.RenderContext, input io.Reader, output io.Writer) error {
	baseLinkSrc := ctx.RenderHelper.ResolveLink("", markup.LinkTypeDefault)
	baseLinkRaw := ctx.RenderHelper.ResolveLink("", markup.LinkTypeRaw)
	command := strings.NewReplacer(
		envMark("KMUP_PREFIX_SRC"), baseLinkSrc,
		envMark("KMUP_PREFIX_RAW"), baseLinkRaw,
	).Replace(p.Command)
	commands, err := shellquote.Split(command)
	if err != nil || len(commands) == 0 {
		return fmt.Errorf("%s invalid command %q: %w", p.Name(), p.Command, err)
	}
	args := commands[1:]

	if p.IsInputFile {
		// write to temp file
		f, cleanup, err := setting.AppDataTempDir("git-repo-content").CreateTempFileRandom("kmup_input")
		if err != nil {
			return fmt.Errorf("%s create temp file when rendering %s failed: %w", p.Name(), p.Command, err)
		}
		defer cleanup()

		_, err = io.Copy(f, input)
		if err != nil {
			_ = f.Close()
			return fmt.Errorf("%s write data to temp file when rendering %s failed: %w", p.Name(), p.Command, err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("%s close temp file when rendering %s failed: %w", p.Name(), p.Command, err)
		}
		args = append(args, f.Name())
	}

	processCtx, _, finished := process.GetManager().AddContext(ctx, fmt.Sprintf("Render [%s] for %s", commands[0], baseLinkSrc))
	defer finished()

	cmd := exec.CommandContext(processCtx, commands[0], args...)
	cmd.Env = append(
		os.Environ(),
		"KMUP_PREFIX_SRC="+baseLinkSrc,
		"KMUP_PREFIX_RAW="+baseLinkRaw,
	)
	if !p.IsInputFile {
		cmd.Stdin = input
	}
	var stderr bytes.Buffer
	cmd.Stdout = output
	cmd.Stderr = &stderr
	process.SetSysProcAttribute(cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s render run command %s %v failed: %w\nStderr: %s", p.Name(), commands[0], args, err, stderr.String())
	}
	return nil
}
