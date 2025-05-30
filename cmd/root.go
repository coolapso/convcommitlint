package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type message struct {
	header header
	body   []string
	footer []string
}

var (
	path          = "./"
	baseBranch    = "main"
	lintAll       bool
	currentCommit bool
	createReview  bool
	prNumber      int
	repository    string
	commentOnly   bool
	Version       = "DEV"
)

var rootCmd = &cobra.Command{
	Use:   "convcommitlint",
	Short: "Lint conventional commits",
	Long:  `A simple, slightly opinionated, yet usable linter for conventional commits written in Go`,
	Run: func(cmd *cobra.Command, args []string) {
		lint()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print convcommitlint version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("convcommitlint %s\n", Version)
	},
}

func lint() {
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal("Failed to open repository")
	}

	headRef, err := r.Head()
	if err != nil {
		log.Fatal("Failed to get head")
	}

	var issuesFound bool
	var issuesMessage string
	if currentCommit {
		commit, err := r.CommitObject(headRef.Hash())
		if err != nil {
			log.Fatal(err)
		}

		message := parseCommitMessage(commit.Message)
		issues := lintCommitMessage(message)
		if issues != nil {
			issuesFound = true
			issuesMessage = issuesMessage + printIssues(issues, commit)
		}
	} else {
		cIter, err := r.Log(&git.LogOptions{From: headRef.Hash()})
		if err != nil {
			log.Fatal("Failed to get commit history")
		}

		var baseRef *plumbing.Reference
		if !lintAll {
			baseRef, err = getBaseRef(r, baseBranch)
			if err != nil {
				log.Fatalf("Failed to get base branch reference: %v", err)
			}
		}

		err = cIter.ForEach(func(c *object.Commit) error {
			if !lintAll {
				if c.Hash == baseRef.Hash() {
					return errStop
				}
			}

			message := parseCommitMessage(c.Message)
			issues := lintCommitMessage(message)
			if issues != nil {
				issuesFound = true
				issuesMessage = issuesMessage + printIssues(issues, c)

			}
			return nil
		})

		if err != nil && err != errStop {
			log.Fatalf("Error iterating commits: %v", err)
		}
	}

	if issuesFound {
		fmt.Print(issuesMessage)
		if createReview {
			err := createPrReview(issuesMessage)
			if err != nil {
				log.Println("Failed to create pull request review: ", err)
			}
		}
		os.Exit(1)
	}
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

func printIssues(issues []error, c *object.Commit) string {
	header := fmt.Sprintf("%s Commit %s by %s:\n", "⛔", c.Hash.String()[:7], c.Author.Name)
	subHeader := fmt.Sprintf("  Message: %s\n", strings.Split(strings.TrimSpace(c.Message), "\n")[0])
	issuesFound := "  Issues found:\n"
	var issuesMsg string
	for _, issue := range issues {
		issuesMsg = issuesMsg + fmt.Sprintf("\t- %s\n", issue.Error())
	}
	issuesMsg = issuesMsg + "\n"

	return fmt.Sprint(header + subHeader + issuesFound + issuesMsg)
}

func parseCommitMessage(commitMessage string) (msg message) {
	lines := strings.Split(commitMessage, "\n")
	header := strings.TrimSpace(lines[0])

	var bodyLines []string
	var footerLines []string
	section := "body"
	for _, line := range lines[1:] {
		if footerRegexp.MatchString(line) {
			section = "footer"
		}

		switch section {
		case "body":
			bodyLines = append(bodyLines, line)
		case "footer":
			footerLines = append(footerLines, strings.TrimRight(line, "\n"))
		}
	}

	msg.header = parseHeader(header)
	msg.body = bodyLines
	msg.footer = footerLines

	return msg
}

func lintCommitMessage(msg message) (errs []error) {

	errs = slices.Concat(errs, lintHeader(msg.header))

	if msg.footer != nil {
		errs = slices.Concat(errs, lintFooter(msg.footer))
	}

	// is  there empty line between header and footer?
	if msg.body == nil && msg.footer != nil {
		errs = append(errs, errInvalidBodyBlankLines)
	}

	// if only footer exists
	if msg.body != nil && msg.footer == nil {
		if !emptyLine(msg.body[0]) {
			errs = append(errs, errInvalidBodyBlankLines)
		}
	}

	//If both body & footer exist
	if msg.body != nil && msg.footer != nil {
		if !emptyLine(msg.body[0]) {
			errs = append(errs, errInvalidBodyBlankLines)
		}

		lastLine := len(msg.body) - 1
		if !emptyLine(msg.body[lastLine]) {
			errs = append(errs, errInvalidBodyBlankLines)
		}
	}

	return errs
}

func emptyLine(s string) bool {
	if s == "" || s == " " || s == "\n" {
		return true
	}

	return false
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.SetEnvPrefix("CONVCOMLINT")
	viper.AutomaticEnv()

	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolVarP(&createReview, "create-review", "r", false, "Creates review on github pull request")
	_ = viper.BindPFlag("create-review", rootCmd.Flags().Lookup("create-review"))
	createReview = viper.GetBool("create_review")

	rootCmd.PersistentFlags().BoolVar(&commentOnly, "comment-only", false, "Pull request reviews will only comment instead of requesting changes")
	_ = viper.BindPFlag("create-review", rootCmd.Flags().Lookup("comment-only"))
	commentOnly = viper.GetBool("comment_only")

	rootCmd.PersistentFlags().IntVar(&prNumber, "pr-number", 0, "The number of pull request to create the review")
	_ = viper.BindPFlag("pr-number", rootCmd.Flags().Lookup("pr-number"))
	prNumber = viper.GetInt("pr_number")

	viper.SetDefault("base_branch", baseBranch)
	rootCmd.PersistentFlags().StringVarP(&baseBranch, "base-branch", "b", baseBranch, "The base branch to check commits from")
	_ = viper.BindPFlag("base-branch", rootCmd.Flags().Lookup("base-branch"))
	baseBranch = viper.GetString("base_branch")

	rootCmd.PersistentFlags().BoolVarP(&lintAll, "lint-all", "a", false, "Lint all repository commits")
	_ = viper.BindPFlag("lint-all", rootCmd.Flags().Lookup("lint-all"))
	lintAll = viper.GetBool("lint_all")

	rootCmd.PersistentFlags().BoolVarP(&currentCommit, "current", "c", false, "Lint only the current commit")
	_ = viper.BindPFlag("current", rootCmd.Flags().Lookup("current"))
	currentCommit = viper.GetBool("current")

	viper.SetDefault("path", path)
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", path, "Git repository path")
	_ = viper.BindPFlag("path", rootCmd.Flags().Lookup("path"))
	path = viper.GetString("path")

	rootCmd.PersistentFlags().StringVar(&repository, "repository", "", "The github repository in owner/name format ex: coolapso/convcommitlint")
	_ = viper.BindPFlag("repository", rootCmd.Flags().Lookup("repository"))
	repository = viper.GetString("repository")

	rootCmd.MarkFlagsMutuallyExclusive("lint-all", "current", "base-branch")
}
