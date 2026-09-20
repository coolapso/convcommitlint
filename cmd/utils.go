package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v72/github"
)

// This function could use a better name
// This function could use some tests, but ghactions doesn't let overide env variables
func initGHSettings() (err error) {
	if repository == "" {
		repository, err = getRepository()
		if err != nil {
			return err
		}
	}

	if !githubAction() && prNumber == 0 {
		return errMissingPRNum
	}

	if pullRequest() {
		prNumber, err = getPRNumber()
		if err != nil {
			return fmt.Errorf("failed to get pull request number from github context GITHUB_REF_NAME, %v", err)
		}
	}

	token, err = getGHToken()
	if err != nil {
		return err
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
	if v := os.Getenv("GITHUB_EVENT_NAME"); v == "pull_request" {
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

// commitsSinceBase returns the commits reachable from head but not from base,
// matching Git's base..head revision range even when the branches have
// diverged. Stopping when a walk first reaches base misses the common ancestor
// when base is not itself an ancestor of head.
func commitsSinceBase(r *git.Repository, head, base plumbing.Hash) ([]*object.Commit, error) {
	baseIter, err := r.Log(&git.LogOptions{From: base})
	if err != nil {
		return nil, err
	}
	defer baseIter.Close()

	baseCommits := make(map[plumbing.Hash]struct{})
	if err := baseIter.ForEach(func(c *object.Commit) error {
		baseCommits[c.Hash] = struct{}{}
		return nil
	}); err != nil {
		return nil, err
	}

	headIter, err := r.Log(&git.LogOptions{From: head})
	if err != nil {
		return nil, err
	}
	defer headIter.Close()

	var commits []*object.Commit
	if err := headIter.ForEach(func(c *object.Commit) error {
		if _, found := baseCommits[c.Hash]; !found {
			commits = append(commits, c)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return commits, nil
}

func emptyLine(s string) bool {
	if s == "" || s == " " || s == "\n" {
		return true
	}

	return false
}
