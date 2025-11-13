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

package alpine

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	packages_model "github.com/kumose/kmup/models/packages"
	"github.com/kumose/kmup/modules/json"
	packages_module "github.com/kumose/kmup/modules/packages"
	alpine_module "github.com/kumose/kmup/modules/packages/alpine"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/routers/api/packages/helper"
	"github.com/kumose/kmup/services/context"
	packages_service "github.com/kumose/kmup/services/packages"
	alpine_service "github.com/kumose/kmup/services/packages/alpine"
)

func apiError(ctx *context.Context, status int, obj any) {
	message := helper.ProcessErrorForUser(ctx, status, obj)
	ctx.PlainText(status, message)
}

func GetRepositoryKey(ctx *context.Context) {
	_, pub, err := alpine_service.GetOrCreateKeyPair(ctx, ctx.Package.Owner.ID)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	pubPem, _ := pem.Decode([]byte(pub))
	if pubPem == nil {
		apiError(ctx, http.StatusInternalServerError, "failed to decode private key pem")
		return
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	fingerprint, err := util.CreatePublicKeyFingerprint(pubKey)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.ServeContent(strings.NewReader(pub), &context.ServeHeaderOptions{
		ContentType: "application/x-pem-file",
		Filename:    fmt.Sprintf("%s@%s.rsa.pub", ctx.Package.Owner.LowerName, hex.EncodeToString(fingerprint)),
	})
}

func GetRepositoryFile(ctx *context.Context) {
	pv, err := alpine_service.GetOrCreateRepositoryVersion(ctx, ctx.Package.Owner.ID)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	s, u, pf, err := packages_service.OpenFileForDownloadByPackageVersion(
		ctx,
		pv,
		&packages_service.PackageFileInfo{
			Filename:     alpine_service.IndexArchiveFilename,
			CompositeKey: fmt.Sprintf("%s|%s|%s", ctx.PathParam("branch"), ctx.PathParam("repository"), ctx.PathParam("architecture")),
		},
		ctx.Req.Method,
	)
	if err != nil {
		if errors.Is(err, util.ErrNotExist) {
			apiError(ctx, http.StatusNotFound, err)
		} else {
			apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	helper.ServePackageFile(ctx, s, u, pf)
}

func UploadPackageFile(ctx *context.Context) {
	branch := strings.TrimSpace(ctx.PathParam("branch"))
	repository := strings.TrimSpace(ctx.PathParam("repository"))
	if branch == "" || repository == "" {
		apiError(ctx, http.StatusBadRequest, "invalid branch or repository")
		return
	}

	upload, needToClose, err := ctx.UploadStream()
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	if needToClose {
		defer upload.Close()
	}

	buf, err := packages_module.CreateHashedBufferFromReader(upload)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer buf.Close()

	pck, err := alpine_module.ParsePackage(buf)
	if err != nil {
		if errors.Is(err, util.ErrInvalidArgument) || errors.Is(err, io.EOF) {
			apiError(ctx, http.StatusBadRequest, err)
		} else {
			apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if _, err := buf.Seek(0, io.SeekStart); err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	fileMetadataRaw, err := json.Marshal(pck.FileMetadata)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	_, _, err = packages_service.CreatePackageOrAddFileToExisting(
		ctx,
		&packages_service.PackageCreationInfo{
			PackageInfo: packages_service.PackageInfo{
				Owner:       ctx.Package.Owner,
				PackageType: packages_model.TypeAlpine,
				Name:        pck.Name,
				Version:     pck.Version,
			},
			Creator:  ctx.Doer,
			Metadata: pck.VersionMetadata,
		},
		&packages_service.PackageFileCreationInfo{
			PackageFileInfo: packages_service.PackageFileInfo{
				Filename:     fmt.Sprintf("%s-%s.apk", pck.Name, pck.Version),
				CompositeKey: fmt.Sprintf("%s|%s|%s", branch, repository, pck.FileMetadata.Architecture),
			},
			Creator: ctx.Doer,
			Data:    buf,
			IsLead:  true,
			Properties: map[string]string{
				alpine_module.PropertyBranch:       branch,
				alpine_module.PropertyRepository:   repository,
				alpine_module.PropertyArchitecture: pck.FileMetadata.Architecture,
				alpine_module.PropertyMetadata:     string(fileMetadataRaw),
			},
		},
	)
	if err != nil {
		switch err {
		case packages_model.ErrDuplicatePackageVersion, packages_model.ErrDuplicatePackageFile:
			apiError(ctx, http.StatusConflict, err)
		case packages_service.ErrQuotaTotalCount, packages_service.ErrQuotaTypeSize, packages_service.ErrQuotaTotalSize:
			apiError(ctx, http.StatusForbidden, err)
		default:
			apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if err := alpine_service.BuildSpecificRepositoryFiles(ctx, ctx.Package.Owner.ID, branch, repository, pck.FileMetadata.Architecture); err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func DownloadPackageFile(ctx *context.Context) {
	branch := ctx.PathParam("branch")
	repository := ctx.PathParam("repository")
	architecture := ctx.PathParam("architecture")

	opts := &packages_model.PackageFileSearchOptions{
		OwnerID:      ctx.Package.Owner.ID,
		PackageType:  packages_model.TypeAlpine,
		Query:        ctx.PathParam("filename"),
		CompositeKey: fmt.Sprintf("%s|%s|%s", branch, repository, architecture),
	}
	pfs, _, err := packages_model.SearchFiles(ctx, opts)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	if len(pfs) == 0 {
		// Try again with architecture 'noarch'
		if architecture == alpine_module.NoArch {
			apiError(ctx, http.StatusNotFound, nil)
			return
		}

		opts.CompositeKey = fmt.Sprintf("%s|%s|%s", branch, repository, alpine_module.NoArch)
		if pfs, _, err = packages_model.SearchFiles(ctx, opts); err != nil {
			apiError(ctx, http.StatusInternalServerError, err)
			return
		}

		if len(pfs) == 0 {
			apiError(ctx, http.StatusNotFound, nil)
			return
		}
	}

	s, u, pf, err := packages_service.OpenFileForDownload(ctx, pfs[0], ctx.Req.Method)
	if err != nil {
		if errors.Is(err, util.ErrNotExist) {
			apiError(ctx, http.StatusNotFound, err)
		} else {
			apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	helper.ServePackageFile(ctx, s, u, pf)
}

func DeletePackageFile(ctx *context.Context) {
	branch, repository, architecture := ctx.PathParam("branch"), ctx.PathParam("repository"), ctx.PathParam("architecture")

	pfs, _, err := packages_model.SearchFiles(ctx, &packages_model.PackageFileSearchOptions{
		OwnerID:      ctx.Package.Owner.ID,
		PackageType:  packages_model.TypeAlpine,
		Query:        ctx.PathParam("filename"),
		CompositeKey: fmt.Sprintf("%s|%s|%s", branch, repository, architecture),
	})
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	if len(pfs) != 1 {
		apiError(ctx, http.StatusNotFound, nil)
		return
	}

	if err := packages_service.RemovePackageFileAndVersionIfUnreferenced(ctx, ctx.Doer, pfs[0]); err != nil {
		if errors.Is(err, util.ErrNotExist) {
			apiError(ctx, http.StatusNotFound, err)
		} else {
			apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if err := alpine_service.BuildSpecificRepositoryFiles(ctx, ctx.Package.Owner.ID, branch, repository, architecture); err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
