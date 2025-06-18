package cmd

import (
	"github.com/google/go-github/v72/github"
	"os"
	"testing"
)

func TestGetGHToken(t *testing.T) {
	t.Run("missing token", func(t *testing.T) {
		want := errMissingGHToken
		_, err := getGHToken()
		if err != want {
			t.Fatalf("Expected error [ %v ], got [ %v ]", want, err)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		want := "valid-token"
		_ = os.Setenv("GITHUB_TOKEN", want)
		defer os.Unsetenv("GITHUB_TOKEN")
		token, err := getGHToken()
		if err != nil {
			t.Fatal("Expected nil, got:", err)
		}

		if token != want {
			t.Fatalf("Expected token [ %v ], got [ %v ]", want, token)
		}
	})
}

func TestPRisDraft(t *testing.T) {
	want := true
	ghpr := &github.PullRequest{
		Draft: &want,
	}
	got := prIsDraft(ghpr)
	if got != want {
		t.Fatalf("Expected [ %v ], got [ %v ]", want, got)
	}
}

func TestGetRepository(t *testing.T) {
	t.Run("missing repository", func(t *testing.T) {
		_ = os.Unsetenv("GITHUB_REPOSITORY")
		want := errMissingRepository
		_, err := getRepository()
		if err != want {
			t.Fatalf("Expected error [ %v ], got [ %v ]", want, err)
		}
	})

	t.Run("valid repository", func(t *testing.T) {
		want := "owner/repo"
		_ = os.Setenv("GITHUB_REPOSITORY", want)
		defer os.Unsetenv("GITHUB_REPOSITORY")
		repo, err := getRepository()
		if err != nil {
			t.Fatal("Expected nil, got:", err)
		}
		if repo != want {
			t.Fatalf("Expected repository [ %v ], got [ %v ]", want, repo)
		}
	})
}

func TestGetPRNumber(t *testing.T) {
	t.Run("valid ref", func(t *testing.T) {
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

func TestEmptyLine(t *testing.T) {
	testsCases := []struct {
		want   bool
		string string
	}{
		{want: true, string: ""},
		{want: true, string: " "},
		{want: true, string: "\n"},
		{want: false, string: "foo"},
	}

	for i, testCase := range testsCases {
		got := emptyLine(testCase.string)
		if testCase.want != got {
			t.Fatalf("case [ %v ]: want [ %v ], got [ %v ]", i, testCase.want, got)
		}
	}
}
