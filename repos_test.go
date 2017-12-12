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
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAPTConfigFolder(t *testing.T) {
	repos, err := ParseAPTConfigFolder("testdata/apt")
	require.NoError(t, err, "running List command")

	expectedData, err := ioutil.ReadFile("testdata/TestParseAPTConfigFolder.json")
	require.NoError(t, err, "Reading test data")
	expected := []*Repository{}
	err = json.Unmarshal(expectedData, &expected)
	require.NoError(t, err, "Decoding expected data")

	for i, repo := range repos {
		require.EqualValues(t, expected[i], repo, "Comparing element %d", i)
	}
}
