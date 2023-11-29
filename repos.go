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
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// RepositoryList is an array of Repository definitions
type RepositoryList []*Repository

// Contains checks if a repository definition is contained
// in the RepositoryList
func (r RepositoryList) Contains(repo *Repository) bool {
	return r.Find(repo) != nil
}

// Find search in the RepositoryList a repo that has the same
// metadata as the one passed as parameter
func (r RepositoryList) Find(repoToFind *Repository) *Repository {
	for _, repo := range r {
		if repoToFind.Equals(repo) {
			return repo
		}
	}
	return nil
}

// Repository contains metadata about a repository installed in the system
type Repository struct {
	Enabled      bool
	SourceRepo   bool
	Options      string
	URI          string
	Distribution string
	Components   string
	Comment      string

	configFile string
}

// Equals check if the Repository metadata are equivalent to the
// one provided as parameter. Two Repository are equivalent if all
// metadata matches with the exception of Enabled and Comment.
func (r *Repository) Equals(repo *Repository) bool {
	if r.Components != repo.Components {
		return false
	}
	if r.Distribution != repo.Distribution {
		return false
	}
	if r.URI != repo.URI {
		return false
	}
	if r.SourceRepo != repo.SourceRepo {
		return false
	}
	if r.Options != repo.Options {
		return false
	}
	return true
}

// APTConfigLine returns the "deb" or "deb-src" config line to put in
// source.list to install the Repository
func (r *Repository) APTConfigLine() string {
	res := ""
	if !r.Enabled {
		res = "# "
	}
	if r.SourceRepo {
		res += "deb-src "
	} else {
		res += "deb "
	}
	if strings.TrimSpace(r.Options) != "" {
		res += "[" + r.Options + "]"
	}
	res += r.URI + " " + r.Distribution + " " + r.Components
	if strings.TrimSpace(r.Comment) != "" {
		res += " # " + r.Comment
	}
	return res
}

var aptConfigLineRegexp = regexp.MustCompile(`^(# )?(deb|deb-src)(?: \[(.*)\])? ([^ ]+) ([^ ]+) ([^#\n]+)(?: +# *(.*))?$`)

func parseAPTConfigLine(line string) *Repository {
	match := aptConfigLineRegexp.FindAllStringSubmatch(line, -1)
	if len(match) == 0 || len(match[0]) < 6 {
		return nil
	}
	fields := match[0]
	//fmt.Printf("%+v\n", fields)
	return &Repository{
		Enabled:      fields[1] != "# ",
		SourceRepo:   fields[2] == "deb-src",
		Options:      fields[3],
		URI:          fields[4],
		Distribution: fields[5],
		Components:   fields[6],
		Comment:      fields[7],
	}
}

func parseAPTConfigFile(configPath string) (RepositoryList, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Reading %s: %s", configPath, err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))

	res := RepositoryList{}
	for scanner.Scan() {
		line := scanner.Text()
		repo := parseAPTConfigLine(line)
		//fmt.Printf("%+v\n", repo)
		if repo != nil {
			repo.configFile = configPath
			res = append(res, repo)
		}
	}
	return res, nil
}

// ParseAPTConfigFolder scans an APT config folder (usually /etc/apt) to
// get information about all configured repositories, it scans also
// "source.list.d" subfolder to find all the "*.list" files.
func ParseAPTConfigFolder(folderPath string) (RepositoryList, error) {
	sources := make([]string, 0)


	sourcesFile := filepath.Join(folderPath, "sources.list")
	if FileExists(sourcesFile) {
		sources = append(sources, sourcesFile)
	}

	sourcesFolder := filepath.Join(folderPath, "sources.list.d")
	list, err := ioutil.ReadDir(sourcesFolder)
	if err != nil {
		return nil, fmt.Errorf("Reading %s folder: %s", sourcesFolder, err)
	}
	for _, l := range list {
		if strings.HasSuffix(l.Name(), ".list") {
			sources = append(sources, filepath.Join(sourcesFolder, l.Name()))
		}
	}

	res := RepositoryList{}
	for _, source := range sources {
		repos, err := parseAPTConfigFile(source)
		if err != nil {
			return nil, fmt.Errorf("Parsing %s: %s", source, err)
		}
		res = append(res, repos...)
	}
	return res, nil
}

// AddRepository adds the specified repository by changing the specified APT
// config folder (usually /etc/apt). The new repository is saved into
// a file named "managed.list"
func AddRepository(repo *Repository, configFolderPath string) error {
	repos, err := ParseAPTConfigFolder(configFolderPath)
	if err != nil {
		return fmt.Errorf("parsing APT config: %s", err)
	}
	if repos.Contains(repo) {
		return fmt.Errorf("The repository is already configured")
	}

	// Add to the "managed.list" file
	managedPath := filepath.Join(configFolderPath, "sources.list.d", "managed.list")
	f, err := os.OpenFile(managedPath, os.O_APPEND|os.O_WRONLY, 0644)
	if os.IsNotExist(err) {
		f, err = os.OpenFile(managedPath, os.O_CREATE|os.O_WRONLY, 0644)
	}
	if err != nil {
		return fmt.Errorf("Opening %s: %s", managedPath, err)
	}
	defer f.Close()
	if _, err = f.WriteString(repo.APTConfigLine() + "\n"); err != nil {
		return fmt.Errorf("Writing repo data to config file %s: %s", managedPath, err)
	}
	return nil
}

// RemoveRepository removes a repository from the repository list files
// found in the specified APT config folder (usually /etc/apt)
func RemoveRepository(repo *Repository, configFolderPath string) error {
	// Read all repos configurations
	repos, err := ParseAPTConfigFolder(configFolderPath)
	if err != nil {
		return fmt.Errorf("parsing APT config: %s", err)
	}

	// Find the repo to remove
	repoToRemove := repos.Find(repo)
	if repoToRemove == nil {
		return fmt.Errorf("Repository already removed")
	}

	// Read the config file that contains the repo config to remove
	fileToFilter := repoToRemove.configFile
	data, err := ioutil.ReadFile(fileToFilter)
	if err != nil {
		return fmt.Errorf("Reading config file %s: %s", fileToFilter, err)
	}

	// Create the new version of the file
	scanner := bufio.NewScanner(bytes.NewReader(data))
	newContent := ""
	for scanner.Scan() {
		line := scanner.Text()
		r := parseAPTConfigLine(line)
		if r!= nil && r.Equals(repo) {
			// Filter repo configs that match the repo to be removed
			continue
		}
		newContent += line + "\n"
	}

	err = replaceFile(fileToFilter, []byte(newContent))
	if err != nil {
		return fmt.Errorf("Writing of new config: %s", err)
	}

	return nil
}

// EditRepository replace an old repo configuration with a new repo
// configuration in the specified APT config folder (usually /etc/apt).
func EditRepository(old *Repository, new *Repository, configFolderPath string) error {
	// Read all repos configurations
	repos, err := ParseAPTConfigFolder(configFolderPath)
	if err != nil {
		return fmt.Errorf("parsing APT config: %s", err)
	}

	// Find the repo to edit
	repoToEdit := repos.Find(old)
	if repoToEdit == nil {
		return fmt.Errorf("Repository doesn't exist")
	}

	// Read the config file that contains the repo configuration to edit
	fileToEdit := repoToEdit.configFile
	data, err := ioutil.ReadFile(fileToEdit)
	if err != nil {
		return fmt.Errorf("Reading config file %s: %s", fileToEdit, err)
	}

	// Create the new version of the file
	scanner := bufio.NewScanner(bytes.NewReader(data))
	newContent := ""
	for scanner.Scan() {
		line := scanner.Text()
		r := parseAPTConfigLine(line)
		if r.Equals(old) {
			// Write the new config to replace the old one
			newContent += new.APTConfigLine() + "\n"
			continue
		}
		newContent += line + "\n"
	}

	err = replaceFile(fileToEdit, []byte(newContent))
	if err != nil {
		return fmt.Errorf("Writing of new config: %s", err)
	}

	return nil
}

func replaceFile(path string, newContent []byte) error {
	newPath := path + ".new"
	backupPath := path + ".save"

	// Create the new version of the file
	err := ioutil.WriteFile(newPath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("Creating replacement file for %s: %s", newPath, err)
	}

	// Only in case of error clean-up the new copy (otherwise ignore the error...)
	defer os.Remove(newPath)

	// Make a backup copy
	err = os.Rename(path, backupPath)
	if err != nil {
		return fmt.Errorf("Making backup copy of %s: %s", path, err)
	}

	// Rename the new copy to the final path
	err = os.Rename(newPath, path)
	if err != nil {
		// Something went wrong... try to rollback the backup
		os.Rename(backupPath, path)
		return fmt.Errorf("Renaming %s to %s: %s", newPath, path, err)
	}

	return nil
}
