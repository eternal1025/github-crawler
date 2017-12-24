package github

import (
	"encoding/json"
	"github.com/0xe8551ccb/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func SaveProjectItem(location string, project *ProjectItem) {
	projectLocation := filepath.Join(location, project.Name)
	err := os.MkdirAll(projectLocation, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join(projectLocation, "summary.json")
	jsonLine, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		log.Printf("failed to marshal item %s\n", project)
	}

	ioutil.WriteFile(filename, jsonLine, os.ModePerm)
}

func SaveIssueItem(location string, issue *IssueItem) {
	issueDir := filepath.Join(location, issue.ProjectName)
	err := os.MkdirAll(issueDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join(issueDir, "issues.jl")
	buf, err := json.Marshal(issue)
	if err != nil {
		log.Printf("failed to marshal item %s\n", issue)
	}

	err = utils.AppendStringToFile(filename, string(buf), true)
	if err != nil {
		log.Fatal(err)
	}
}
