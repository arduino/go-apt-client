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
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	out, err := os.ReadFile("testdata/dpkg-query-output-1.txt")
	require.NoError(t, err, "Reading test input data")
	list := parseDpkgQueryOutput(out)

	// Check list with expected output
	data, err := os.ReadFile("testdata/dpkg-query-output-1-result.json")
	require.NoError(t, err, "Reading test result data")
	var expected []*Package
	err = json.Unmarshal(data, &expected)
	require.NoError(t, err, "Unmarshaling test result data")
	require.Equal(t, len(expected), len(list), "Length of result")
	for i := range expected {
		require.Equal(t, expected[i], list[i], "Element", i, "of the result")
	}
}

func TestSearch(t *testing.T) {
	list, err := Search("nonexisting")
	require.NoError(t, err, "running Search command")
	require.Empty(t, list, "Search command result")

	list, err = Search("bash") // "bash" is almost always present on Linux systems
	require.NoError(t, err, "running Search command")
	require.NotEmpty(t, list, "Search command result")
}

func TestListUpgradable(t *testing.T) {
	list, err := ListUpgradable()
	for _, p := range list {
		fmt.Printf("%+v\n", p)
	}
	require.NoError(t, err, "running List command")
}

// func TestCheckForUpdates(t *testing.T) {
// 	out, err := CheckForUpdates()
// 	require.NoError(t, err, "running CheckForUpdate command")
// 	fmt.Printf(">>>\n%s\n<<<\n", string(out))
// 	fmt.Println("ERR:", err)
// }
