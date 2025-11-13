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

package user

const (
	// SettingsKeyHiddenCommentTypes is the setting key for hidden comment types
	SettingsKeyHiddenCommentTypes = "issue.hidden_comment_types"
	// SettingsKeyDiffWhitespaceBehavior is the setting key for whitespace behavior of diff
	SettingsKeyDiffWhitespaceBehavior = "diff.whitespace_behaviour"
	// SettingsKeyShowOutdatedComments is the setting key whether or not to show outdated comments in PRs
	SettingsKeyShowOutdatedComments = "comment_code.show_outdated"

	// UserActivityPubPrivPem is user's private key
	UserActivityPubPrivPem = "activitypub.priv_pem"
	// UserActivityPubPubPem is user's public key
	UserActivityPubPubPem = "activitypub.pub_pem"
	// SignupIP is the IP address that the user signed up with
	SignupIP = "signup.ip"
	// SignupUserAgent is the user agent that the user signed up with
	SignupUserAgent = "signup.user_agent"

	SettingsKeyCodeViewShowFileTree = "code_view.show_file_tree"

	SettingsKeyEmailNotificationKmupActions        = "email_notification.kmup_actions"
	SettingEmailNotificationKmupActionsAll         = "all"
	SettingEmailNotificationKmupActionsFailureOnly = "failure-only" // Default for actions email preference
	SettingEmailNotificationKmupActionsDisabled    = "disabled"
)
