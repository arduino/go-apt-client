//
//  This file is part of go-apt-client library
//
//  Copyright (C) 2017  Arduino AG (http://www.arduino.cc/)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package apt

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestParseAPTConfigFolder(t *testing.T) {
	repos, err := ParseAPTConfigFolder("testdata/apt")
	require.NoError(t, err, "running List command")

	expectedData, err := os.ReadFile("testdata/TestParseAPTConfigFolder.json")
	require.NoError(t, err, "Reading test data")
	expected := []*Repository{}
	err = json.Unmarshal(expectedData, &expected)
	require.NoError(t, err, "Decoding expected data")

	for i, repo := range repos {
		assert.Empty(t, cmp.Diff(expected[i], repo, cmpopts.IgnoreFields(Repository{}, "configFile")))
	}
}

func TestAddAndRemoveRepository(t *testing.T) {
	// test cleanup
	defer os.Remove("testdata/apt3/sources.list.d/managed.list")      //nolint:errcheck
	defer os.Remove("testdata/apt3/sources.list.d/managed.list.save") //nolint:errcheck
	defer os.Remove("testdata/apt3/sources.list.d/managed.list.new")  //nolint:errcheck

	repo1 := &Repository{
		Enabled:      true,
		SourceRepo:   false,
		URI:          "http://ppa.launchpad.net/webupd8team/java/ubuntu",
		Distribution: "zesty",
		Components:   "main",
		Comment:      "",
	}
	repo2 := &Repository{
		Enabled:      false,
		SourceRepo:   true,
		URI:          "http://ppa.launchpad.net/webupd8team/java/ubuntu",
		Distribution: "zesty",
		Components:   "main",
		Comment:      "",
	}
	err := AddRepository(repo1, "testdata/apt3")
	require.NoError(t, err, "Adding repository")
	err = AddRepository(repo2, "testdata/apt3")
	require.NoError(t, err, "Adding repository")

	// check that we have repo1 and repo2 added
	repos, err := ParseAPTConfigFolder("testdata/apt3")
	require.NoError(t, err, "running List command")
	require.True(t, repos.Contains(repo1), "Configuration contains: %#v", repo1)
	require.True(t, repos.Contains(repo2), "Configuration contains: %#v", repo2)

	err = AddRepository(repo2, "testdata/apt3")
	require.Error(t, err, "Adding repository again")

	// no changes should have happened
	repos, err = ParseAPTConfigFolder("testdata/apt3")
	require.NoError(t, err, "running List command")
	require.True(t, repos.Contains(repo1), "Configuration contains: %#v", repo1)
	require.True(t, repos.Contains(repo2), "Configuration contains: %#v", repo2)

	err = RemoveRepository(repo2, "testdata/apt3")
	require.NoError(t, err, "Removing repository")

	// repo2 should be removed
	repos, err = ParseAPTConfigFolder("testdata/apt3")
	require.NoError(t, err, "running List command")
	require.True(t, repos.Contains(repo1), "Configuration contains: %#v", repo1)
	require.False(t, repos.Contains(repo2), "Configuration contains: %#v", repo2)

	err = RemoveRepository(repo2, "testdata/apt3")
	require.Error(t, err, "Removing repository again")

	// no changes should have happened
	repos, err = ParseAPTConfigFolder("testdata/apt3")
	require.NoError(t, err, "running List command")
	require.True(t, repos.Contains(repo1), "Configuration contains: %#v", repo1)
	require.False(t, repos.Contains(repo2), "Configuration contains: %#v", repo2)

	err = EditRepository(repo1, repo2, "testdata/apt3")
	require.NoError(t, err, "editing repository %#V -> %#V", repo1, repo2)

	// repo2 should be changed to repo1
	repos, err = ParseAPTConfigFolder("testdata/apt3")
	require.NoError(t, err, "running List command")
	require.False(t, repos.Contains(repo1), "Configuration contains: %#v", repo1)
	require.True(t, repos.Contains(repo2), "Configuration contains: %#v", repo2)

	err = EditRepository(repo1, repo2, "testdata/apt3")
	require.Error(t, err, "editing again repository %#v -> %#v", repo1, repo2)

	// no changes should have happened
	repos, err = ParseAPTConfigFolder("testdata/apt3")
	require.NoError(t, err, "running List command")
	require.False(t, repos.Contains(repo1), "Configuration contains: %#v", repo1)
	require.True(t, repos.Contains(repo2), "Configuration contains: %#v", repo2)
}
