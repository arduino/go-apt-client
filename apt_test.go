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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	list, err := List()
	require.NoError(t, err, "running List command")
	require.NotEmpty(t, list, "List command result")

	var dpkg *Package
	statusCount := map[string]int{}
	for _, p := range list {
		if p.Name == "dpkg" {
			dpkg = p
			continue
		}
		statusCount[p.Status]++
		// fmt.Printf("%+v\n", p)
	}

	// fmt.Println("Summary:")
	// for k, v := range statusCount {
	// 	fmt.Printf("  %s: %d\n", k, v)
	// }

	require.NotNil(t, dpkg, "search package 'dpkg'")
	require.Equal(t, "installed", dpkg.Status, "'dpkg' status")
}

func TestSearch(t *testing.T) {
	list, err := Search("nonexisting")
	require.NoError(t, err, "running Search command")
	require.Empty(t, list, "Search command result")

	list, err = Search("header")
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
