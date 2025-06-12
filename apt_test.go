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
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDpkgQueryOutput(t *testing.T) {
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

func TestCheckForUpdates(t *testing.T) {
	out, err := CheckForUpdates()
	require.NoError(t, err, "running CheckForUpdate command")
	fmt.Printf(">>>\n%s\n<<<\n", string(out))
	fmt.Println("ERR:", err)
}

func TestParseListUpgradableOutput(t *testing.T) {
	t.Run("edges cases", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected []*Package
		}{
			{
				name:     "empty input",
				input:    "",
				expected: []*Package{},
			},
			{
				name:     "line not matching regex",
				input:    "this-is-not a-valid-line\n",
				expected: []*Package{},
			},
			{
				name:  "upgradable package without [upgradable from]",
				input: "nano/bionic-updates 2.9.3-2 amd64\n",
				expected: []*Package{
					{
						Name:         "nano",
						Status:       "upgradable",
						Version:      "2.9.3-2",
						Architecture: "amd64",
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := parseListUpgradableOutput(strings.NewReader(tt.input))
				require.Equal(t, tt.expected, res)
			})
		}
	})

	t.Run("golden file: list-upgradable.golden", func(t *testing.T) {
		data, err := os.ReadFile("testdata/apt-list-upgradable.golden")
		require.NoError(t, err, "Reading golden file")
		result := parseListUpgradableOutput(strings.NewReader(string(data)))

		want := []*Package{
			{Name: "apt-transport-https", Status: "upgradable", Version: "2.0.11", Architecture: "all"},
			{Name: "apt-utils", Status: "upgradable", Version: "2.0.11", Architecture: "amd64"},
			{Name: "apt", Status: "upgradable", Version: "2.0.11", Architecture: "amd64"},
			{Name: "code-insiders", Status: "upgradable", Version: "1.101.0-1749657374", Architecture: "amd64"},
			{Name: "code", Status: "upgradable", Version: "1.100.3-1748872405", Architecture: "amd64"},
			{Name: "containerd.io", Status: "upgradable", Version: "1.7.27-1", Architecture: "amd64"},
			{Name: "distro-info-data", Status: "upgradable", Version: "0.43ubuntu1.18", Architecture: "all"},
			{Name: "docker-ce-cli", Status: "upgradable", Version: "5:28.1.1-1~ubuntu.20.04~focal", Architecture: "amd64"},
			{Name: "python3.12", Status: "upgradable", Version: "3.12.11-1+focal1", Architecture: "amd64"},
			{Name: "xdg-desktop-portal", Status: "upgradable", Version: "1.14.3-1~flatpak1~20.04", Architecture: "amd64"},
			{Name: "xserver-common", Status: "upgradable", Version: "2:1.20.13-1ubuntu1~20.04.20", Architecture: "all"},
			{Name: "xserver-xephyr", Status: "upgradable", Version: "2:1.20.13-1ubuntu1~20.04.20", Architecture: "amd64"},
			{Name: "xserver-xorg-core", Status: "upgradable", Version: "2:1.20.13-1ubuntu1~20.04.20", Architecture: "amd64"},
			{Name: "xserver-xorg-legacy", Status: "upgradable", Version: "2:1.20.13-1ubuntu1~20.04.20", Architecture: "amd64"},
			{Name: "xwayland", Status: "upgradable", Version: "2:1.20.13-1ubuntu1~20.04.20", Architecture: "amd64"},
		}
		require.NotNil(t, result)
		require.Equal(t, want, result, "Parsed result should match expected from golden file")
	})
}
