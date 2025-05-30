package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v72/github"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	reviewTitle     = "## Conventional commit lint\n"
	reviewTitleText = "**This repository uses [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/), the following issues were found:**\n\n"
	whatNowTitle    = "## What can you do now?\n"
	whatNowText     = "There are more than one way of fixing this issues, using a interactive rebase (`git rebase -i <target branch>`) to reword/squash your commits and then force pushing towards your branch / fork is probably the most adequate."
)

// reviewGHPR creates a github pull request review
// If no PR Number is provided in flags or env variable, tries to guess it from github actions context
func createPrReview(issuesMessage string) (err error) {
	token, exists := os.LookupEnv("GITHUB_TOKEN")
	if !exists {
		return errMissingGHToken
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
		if repository, exists = os.LookupEnv("GITHUB_REPOSITORY"); !exists {
			return errMissingRepository
		}
	}

	owner, repo := splitOwnerRepo(repository)
	event := "REQUEST_CHANGES"
	if commentOnly {
		event = "COMMENT"
	}

	markdownIssues := strings.ReplaceAll(issuesMessage, "\t", "")
	reviewMessage := reviewTitle + reviewTitleText + markdownIssues + whatNowTitle + whatNowText

	client := github.NewClient(nil).WithAuthToken(token)
	review := github.PullRequestReviewRequest{
		Body:  &reviewMessage,
		Event: &event,
	}

	fmt.Println(owner + " " + repo)
	_, resp, err := client.PullRequests.CreateReview(context.TODO(), owner, repo, prNumber, &review)
	if err != nil {
		log.Fatal("failed to create pull request review: ", err)
	}

	if resp.StatusCode != 200 {
		statusMsg := fmt.Sprintf("%v, %v", resp.StatusCode, resp.Status)
		log.Fatalf("Failed to create pull request review, status: %v,", statusMsg)
	}

	return nil
}

func splitOwnerRepo(repository string) (owner, repo string) {
	owner = strings.Split(repository, "/")[0]
	repo = strings.Split(repository, "/")[1]

	return owner, repo
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
