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
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Package is a package available in the APT system
type Package struct {
	Name         string
	Status       string
	Architecture string
	Version      string
}

// List returns a list of packages available in the system with their
// respective status.
func List() ([]*Package, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Package}\t${Architecture}\t${db:Status-Status}\t${Version}\n")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running dpkg-query: %s", err)
	}

	res := []*Package{}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), "\t")
		res = append(res, &Package{
			Name:         data[0],
			Architecture: data[1],
			Status:       data[2],
			Version:      data[3],
		})
	}

	return res, nil
}

// CheckForUpdates runs an apt update to retrieve new packages available
// from the repositories
func CheckForUpdates() (output []byte, err error) {
	cmd := exec.Command("apt-get", "update", "-q")
	return cmd.CombinedOutput()
}

// ListUpgradable return all the upgradable packages and the version that
// is going to be installed if an UpgradeAll is performed
func ListUpgradable() ([]*Package, error) {
	cmd := exec.Command("apt", "list", "--upgradable")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running apt list: %s", err)
	}
	re := regexp.MustCompile(`^([^ ]+) ([^ ]+) ([^ ]+)( \[upgradable from: [^\[\]]*\])?`)

	res := []*Package{}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		matches := re.FindAllStringSubmatch(scanner.Text(), -1)
		if len(matches) == 0 {
			continue
		}
		res = append(res, &Package{
			Name:         matches[0][1],
			Status:       "upgradable",
			Version:      matches[0][2],
			Architecture: matches[0][3],
		})
	}
	return res, nil
}

// Upgrade runs the upgrade for a set of packages
func Upgrade(packs ...*Package) (output []byte, err error) {
	args := []string{"upgrade", "-y"}
	for _, pack := range packs {
		if pack == nil || pack.Name == "" {
			return nil, fmt.Errorf("apt.Upgrade: Invalid package with empty Name")
		}
		args = append(args, pack.Name)
	}
	cmd := exec.Command("apt-get", args...)
	return cmd.CombinedOutput()
}

// UpgradeAll upgrade all upgradable packages
func UpgradeAll() (output []byte, err error) {
	cmd := exec.Command("apt-get", "upgrade", "-y")
	return cmd.CombinedOutput()
}

// Remove removes a set of packages
func Remove(packs ...*Package) (output []byte, err error) {
	args := []string{"remove", "-y"}
	for _, pack := range packs {
		if pack == nil || pack.Name == "" {
			return nil, fmt.Errorf("apt.Remove: Invalid package with empty Name")
		}
		args = append(args, pack.Name)
	}
	cmd := exec.Command("apt-get", args...)
	return cmd.CombinedOutput()
}
