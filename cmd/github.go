package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v72/github"
	"log"
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
func createPrReview(c *github.Client, issuesMessage string) (err error) {
	owner, repo := splitOwnerRepo(repository)
	event := "REQUEST_CHANGES"
	if commentOnly {
		event = "COMMENT"
	}

	markdownIssues := strings.ReplaceAll(issuesMessage, "\t", "")
	reviewMessage := reviewTitle + reviewTitleText + markdownIssues + whatNowTitle + whatNowText

	review := github.PullRequestReviewRequest{
		Body:  &reviewMessage,
		Event: &event,
	}

	fmt.Println(owner + " " + repo)
	_, resp, err := c.PullRequests.CreateReview(context.TODO(), owner, repo, prNumber, &review)
	if err != nil {
		log.Fatal("failed to create pull request review: ", err)
	}

	if resp.StatusCode != 200 {
		statusMsg := fmt.Sprintf("%v, %v", resp.StatusCode, resp.Status)
		log.Fatalf("Failed to create pull request review, status: %v,", statusMsg)
	}

	return nil
}
