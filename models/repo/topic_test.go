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

package repo_test

import (
	"testing"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestAddTopic(t *testing.T) {
	totalNrOfTopics := 6
	repo1NrOfTopics := 3

	assert.NoError(t, unittest.PrepareTestDatabase())

	topics, err := db.Find[repo_model.Topic](t.Context(), &repo_model.FindTopicOptions{})
	assert.NoError(t, err)
	assert.Len(t, topics, totalNrOfTopics)

	topics, total, err := db.FindAndCount[repo_model.Topic](t.Context(), &repo_model.FindTopicOptions{
		ListOptions: db.ListOptions{Page: 1, PageSize: 2},
	})
	assert.NoError(t, err)
	assert.Len(t, topics, 2)
	assert.EqualValues(t, 6, total)

	topics, err = db.Find[repo_model.Topic](t.Context(), &repo_model.FindTopicOptions{
		RepoID: 1,
	})
	assert.NoError(t, err)
	assert.Len(t, topics, repo1NrOfTopics)

	assert.NoError(t, repo_model.SaveTopics(t.Context(), 2, "golang"))
	repo2NrOfTopics := 1
	topics, err = db.Find[repo_model.Topic](t.Context(), &repo_model.FindTopicOptions{})
	assert.NoError(t, err)
	assert.Len(t, topics, totalNrOfTopics)

	topics, err = db.Find[repo_model.Topic](t.Context(), &repo_model.FindTopicOptions{
		RepoID: 2,
	})
	assert.NoError(t, err)
	assert.Len(t, topics, repo2NrOfTopics)

	assert.NoError(t, repo_model.SaveTopics(t.Context(), 2, "golang", "kmup"))
	repo2NrOfTopics = 2
	totalNrOfTopics++
	topic, err := repo_model.GetTopicByName(t.Context(), "kmup")
	assert.NoError(t, err)
	assert.Equal(t, 1, topic.RepoCount)

	topics, err = db.Find[repo_model.Topic](t.Context(), &repo_model.FindTopicOptions{})
	assert.NoError(t, err)
	assert.Len(t, topics, totalNrOfTopics)

	topics, err = db.Find[repo_model.Topic](t.Context(), &repo_model.FindTopicOptions{
		RepoID: 2,
	})
	assert.NoError(t, err)
	assert.Len(t, topics, repo2NrOfTopics)
}

func TestTopicValidator(t *testing.T) {
	assert.True(t, repo_model.ValidateTopic("12345"))
	assert.True(t, repo_model.ValidateTopic("2-test"))
	assert.True(t, repo_model.ValidateTopic("foo.bar"))
	assert.True(t, repo_model.ValidateTopic("test-3"))
	assert.True(t, repo_model.ValidateTopic("first"))
	assert.True(t, repo_model.ValidateTopic("second-test-topic"))
	assert.True(t, repo_model.ValidateTopic("third-project-topic-with-max-length"))

	assert.False(t, repo_model.ValidateTopic("$fourth-test,topic"))
	assert.False(t, repo_model.ValidateTopic("-fifth-test-topic"))
	assert.False(t, repo_model.ValidateTopic("sixth-go-project-topic-with-excess-length"))
	assert.False(t, repo_model.ValidateTopic(".foo"))
}
