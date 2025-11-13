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

package context

import (
	"fmt"
	"sync"

	"github.com/kumose/kmup/modules/cache"
	"github.com/kumose/kmup/modules/hcaptcha"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/mcaptcha"
	"github.com/kumose/kmup/modules/recaptcha"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/turnstile"

	"github.com/kumose-go/chi/captcha"
)

var (
	imageCaptchaOnce sync.Once
	cpt              *captcha.Captcha
)

// GetImageCaptcha returns global image captcha
func GetImageCaptcha() *captcha.Captcha {
	imageCaptchaOnce.Do(func() {
		cpt = captcha.NewCaptcha(captcha.Options{
			SubURL: setting.AppSubURL,
		})
		cpt.Store = cache.GetCache().ChiCache()
	})
	return cpt
}

// SetCaptchaData sets common captcha data
func SetCaptchaData(ctx *Context) {
	if !setting.Service.EnableCaptcha {
		return
	}
	ctx.Data["EnableCaptcha"] = setting.Service.EnableCaptcha
	ctx.Data["RecaptchaURL"] = setting.Service.RecaptchaURL
	ctx.Data["Captcha"] = GetImageCaptcha()
	ctx.Data["CaptchaType"] = setting.Service.CaptchaType
	ctx.Data["RecaptchaSitekey"] = setting.Service.RecaptchaSitekey
	ctx.Data["HcaptchaSitekey"] = setting.Service.HcaptchaSitekey
	ctx.Data["McaptchaSitekey"] = setting.Service.McaptchaSitekey
	ctx.Data["McaptchaURL"] = setting.Service.McaptchaURL
	ctx.Data["CfTurnstileSitekey"] = setting.Service.CfTurnstileSitekey
}

const (
	gRecaptchaResponseField  = "g-recaptcha-response"
	hCaptchaResponseField    = "h-captcha-response"
	mCaptchaResponseField    = "m-captcha-response"
	cfTurnstileResponseField = "cf-turnstile-response"
)

// VerifyCaptcha verifies Captcha data
// No-op if captchas are not enabled
func VerifyCaptcha(ctx *Context, tpl templates.TplName, form any) {
	if !setting.Service.EnableCaptcha {
		return
	}

	var valid bool
	var err error
	switch setting.Service.CaptchaType {
	case setting.ImageCaptcha:
		valid = GetImageCaptcha().VerifyReq(ctx.Req)
	case setting.ReCaptcha:
		valid, err = recaptcha.Verify(ctx, ctx.Req.Form.Get(gRecaptchaResponseField))
	case setting.HCaptcha:
		valid, err = hcaptcha.Verify(ctx, ctx.Req.Form.Get(hCaptchaResponseField))
	case setting.MCaptcha:
		valid, err = mcaptcha.Verify(ctx, ctx.Req.Form.Get(mCaptchaResponseField))
	case setting.CfTurnstile:
		valid, err = turnstile.Verify(ctx, ctx.Req.Form.Get(cfTurnstileResponseField))
	default:
		ctx.ServerError("Unknown Captcha Type", fmt.Errorf("unknown Captcha Type: %s", setting.Service.CaptchaType))
		return
	}
	if err != nil {
		log.Debug("Captcha Verify failed: %v", err)
	}

	if !valid {
		ctx.Data["Err_Captcha"] = true
		ctx.RenderWithErr(ctx.Tr("form.captcha_incorrect"), tpl, form)
	}
}
