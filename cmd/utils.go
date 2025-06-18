package cmd

import (
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v72/github"
)

// This function could use a better name
func initGHSettings() (err error) {
	token, err = getGHToken()
	if err != nil {
		return err
	}

	if prNumber == 0 {
		if !githubAction() {
			return errMissingPRNum
		}

		prNumber, err = getPRNumber()
		if err != nil {
			return fmt.Errorf("failed to get pull request number from github context GITHUB_REF_NAME, %v", err)
		}
	}

	if repository == "" {
		repository, err = getRepository()
		if err != nil {
			return err
		}
	}

	return nil
}

func getGHToken() (string, error) {
	token, exists := os.LookupEnv("GITHUB_TOKEN")
	if !exists {
		return "", errMissingGHToken
	}

	return token, nil
}

func getRepository() (string, error) {
	repo, exists := os.LookupEnv("GITHUB_REPOSITORY")
	if !exists {
		return "", errMissingRepository
	}	

	return repo, nil
}

func prIsDraft(pr *github.PullRequest) bool { 
	return pr.GetDraft()
}

func getPRNumber() (n int, err error) {
	return strconv.Atoi(strings.Split(os.Getenv("GITHUB_REF_NAME"), "/")[0])
}

func githubAction() bool {
	if _, e := os.LookupEnv("GITHUB_ACTIONS"); e {
		return true
	}

	return false
}

func pullRequest() bool {
	if v := os.Getenv("GITHUB_EVENT_TYPE"); v == "pull_request" {
		return true
	}

	return false
}

func splitOwnerRepo(repository string) (owner, repo string) {
	owner = strings.Split(repository, "/")[0]
	repo = strings.Split(repository, "/")[1]

	return owner, repo
}

func getBaseRef(r *git.Repository, branchName string) (baseRef *plumbing.Reference, err error) {
	if !lintAll {
		baseRefName := plumbing.NewBranchReferenceName(branchName)
		baseRef, err = r.Reference(baseRefName, true)
		if err != nil {
			return nil, err
		}
	}

	return baseRef, nil
}

func emptyLine(s string) bool {
	if s == "" || s == " " || s == "\n" {
		return true
	}

	return false
}
