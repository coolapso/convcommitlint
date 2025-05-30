package cmd

import (
	"os"
	"testing"
)

func TestReviewGHPr(t *testing.T) {
	t.Run("missing github token", func(t *testing.T) {
		_ = os.Unsetenv("GITHUB_TOKEN")
		want := errMissingGHToken
		got := createPrReview("")

		if got != want {
			t.Fatalf("Want [ %v ], got [ %v ]", want, got)
		}
	})

	t.Run("Misisng PR Number, local", func(t *testing.T) {
		_ = os.Setenv("GITHUB_TOKEN", "foo")
		_ = os.Unsetenv("GITHUB_ACTIONS")
		want := errMissingPRNum
		got := createPrReview("")

		if got != want {
			t.Fatalf("want [ %v ], got [ %v ]", want, got)
		}
	})

	t.Run("Missing PR Number githubAction", func(t *testing.T) {
		_ = os.Setenv("GITHUB_TOKEN", "foo")
		_ = os.Setenv("GITHUB_ACTIONS", "true")
		_ = os.Unsetenv("GITHUB_REF_NAME")
		got := createPrReview("")
		if got == nil {
			t.Fatalf("Want error, got nil")
		}
	})

	t.Run("Missing Repository, local", func(t *testing.T) {
		_ = os.Setenv("GITHUB_TOKEN", "foo")
		_ = os.Unsetenv("GITHUB_REPOSITORY")
		prNumber = 1
		want := errMissingRepository
		got := createPrReview("")

		if got != want {
			t.Fatalf("want [ %v ], got [ %v ]", want, got)
		}
	})

}

func TestGetPRNumber(t *testing.T) {
	t.Run("Valid ref", func(t *testing.T) {
		_ = os.Setenv("GITHUB_REF_NAME", "25/merge")
		want := 25
		got, err := getPRNumber()
		if err != nil {
			t.Fatalf("want nil, got error [ %v ]", err)
		}

		if got != want {
			t.Fatalf("Want [ %v ], got [ %v ]", want, got)
		}
	})

	t.Run("invalid ref", func(t *testing.T) {
		_ = os.Setenv("GITHUB_REF_NAME", "foo/merge")
		_, err := getPRNumber()
		if err == nil {
			t.Fatal("want error, got nil")
		}
	})

	t.Run("invalid ref", func(t *testing.T) {
		_ = os.Unsetenv("GITHUB_REF_NAME")
		_, err := getPRNumber()
		if err == nil {
			t.Fatal("want error, got nil")
		}
	})
}

func TestPullRequest(t *testing.T) {
	t.Run("Is pull request", func(t *testing.T) {
		_ = os.Setenv("GITHUB_EVENT_TYPE", "pull_request")
		want := true
		got := pullRequest()
		if got != want {
			t.Fatalf("Want [ %v ], got [ %v ]", want, got)
		}
	})

	t.Run("Is not pull request", func(t *testing.T) {
		_ = os.Setenv("GITHUB_EVENT_TYPE", "push")
		want := false
		got := pullRequest()
		if got != want {
			t.Fatalf("Want [ %v ], got [ %v ]", want, got)
		}
	})
}

func TestGithubAction(t *testing.T) {
	t.Run("Is github action", func(t *testing.T) {
		_ = os.Setenv("GITHUB_ACTIONS", "true")
		want := true
		got := githubAction()
		if got != want {
			t.Fatalf("Want [ %v ], got [ %v ]", want, got)
		}
		_ = os.Unsetenv("GITHUB_ACTIONS")
	})

	t.Run("Is not github action", func(t *testing.T) {
		want := false
		got := githubAction()
		if got != want {
			t.Fatalf("Want [ %v ], got [ %v ]", want, got)
		}
	})
}

func TestSplitOwnerRepo(t *testing.T) {
	repository = "coolapso/convcommitlint"
	wantOwner := "coolapso"
	wantRepo := "convcommitlint"
	gotOwner, gotRepo := splitOwnerRepo(repository)
	if gotOwner != wantOwner {
		t.Fatalf("Want owner [ %v ], got [ %v ]", gotOwner, wantOwner)
	}

	if gotRepo != wantRepo {
		t.Fatalf("Want repo [ %v ], got [ %v ]", gotRepo, wantRepo)
	}
}
