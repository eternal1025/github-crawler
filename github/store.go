package github

import (
	"path"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/0xe8551ccb/utils"
)

func SaveProjectItem(location string, project *ProjectItem)  {
	projectLocation := path.Join(location, project.Name)
	err := os.MkdirAll(projectLocation, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := path.Join(projectLocation, "summary.json")
	jsonLine, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		log.Printf("failed to marshal item %s\n", project)
	}

	ioutil.WriteFile(filename, jsonLine, os.ModePerm)
}

func SaveIssueItem(location string, issue *IssueItem)  {
	issueDir := path.Join(location, issue.ProjectName)
	err := os.MkdirAll(issueDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := path.Join(issueDir, "issues.jl")
	buf, err := json.Marshal(issue)
	if err != nil {
		log.Printf("failed to marshal item %s\n", issue)
	}

	err = utils.AppendStringToFile(filename, string(buf), true)
	if err != nil {
		log.Fatal(err)
	}
}