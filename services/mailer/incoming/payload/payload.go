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

package payload

import (
	"context"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/util"
)

const replyPayloadVersion1 byte = 1

type payloadReferenceType byte

const (
	payloadReferenceIssue payloadReferenceType = iota
	payloadReferenceComment
)

// CreateReferencePayload creates data which GetReferenceFromPayload resolves to the reference again.
func CreateReferencePayload(reference any) ([]byte, error) {
	var refType payloadReferenceType
	var refID int64

	switch r := reference.(type) {
	case *issues_model.Issue:
		refType = payloadReferenceIssue
		refID = r.ID
	case *issues_model.Comment:
		refType = payloadReferenceComment
		refID = r.ID
	default:
		return nil, util.NewInvalidArgumentErrorf("unsupported reference type: %T", r)
	}

	payload, err := util.PackData(refType, refID)
	if err != nil {
		return nil, err
	}

	return append([]byte{replyPayloadVersion1}, payload...), nil
}

// GetReferenceFromPayload resolves the reference from the payload
func GetReferenceFromPayload(ctx context.Context, payload []byte) (any, error) {
	if len(payload) < 1 {
		return nil, util.NewInvalidArgumentErrorf("payload to small")
	}

	if payload[0] != replyPayloadVersion1 {
		return nil, util.NewInvalidArgumentErrorf("unsupported payload version")
	}

	var ref payloadReferenceType
	var id int64
	if err := util.UnpackData(payload[1:], &ref, &id); err != nil {
		return nil, err
	}

	switch ref {
	case payloadReferenceIssue:
		return issues_model.GetIssueByID(ctx, id)
	case payloadReferenceComment:
		return issues_model.GetCommentByID(ctx, id)
	default:
		return nil, util.NewInvalidArgumentErrorf("unsupported reference type: %T", ref)
	}
}
