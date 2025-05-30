package cmd

import (
	"slices"
	"testing"
)

func TestLintCommitMessage(t *testing.T) {
	t.Run("header only", func(t *testing.T) {
		commit := "type(scope): description"

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid single line body", func(t *testing.T) {
		commit := `type(scope): description

commit messasge body`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid single line body and footer", func(t *testing.T) {
		commit := `type(scope): description

commit messasge body

key: value`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid empty line end of footer", func(t *testing.T) {
		commit := `type(scope): description

commit messasge body

key: value
`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid header and trailer only", func(t *testing.T) {
		commit := `type(scope): description

key: value`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid header and multi line trailer only", func(t *testing.T) {
		commit := `type(scope): description

key: value
key1: value1`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid trailer with long lines", func(t *testing.T) {
		commit := `type(scope): description

commit messasge body

key: this is a very long
 value
`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid multi trailer with long lines", func(t *testing.T) {
		commit := `type(scope): description

commit messasge body

key: this is a very long
 value
key1: This is another footer
key2: This is another very
 long footer
`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid multi line body", func(t *testing.T) {
		commit := `type(scope): description

commit messasge body
another comit message`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid body with multi blank line separated", func(t *testing.T) {
		commit := `type(scope): description

commit messasge body

another comit message`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("valid body with multi blank line separated & footer", func(t *testing.T) {
		commit := `type(scope): description

commit messasge body

another comit message

key: value
key1: value with big
 content
key2: value`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got != nil {
			t.Fatal("Expected no issues got:", got)
		}
	})

	t.Run("invalid Missing separator", func(t *testing.T) {
		commit := `feat(scope) description`
		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got == nil {
			t.Fatal("Expected invalid hearror header got nil")
		}

		if !slices.Contains(got, errInvalidHeaderFormat) {
			t.Fatal("want invalid header error got none header, got: ", got)
		}
	})

	t.Run("invalid, missing empty lines", func(t *testing.T) {
		commits := []string{
			"feat(scope): description\nkey: value",
			"feat(scope): description\nkey: value\nkey1: value1",
		}
		want := errInvalidBodyBlankLines

		// Messages get really ugly when using the commit on the error message
		// Using the index instead
		for i, commit := range commits {
			m := parseCommitMessage(commit)
			got := lintCommitMessage(m)
			if got == nil {
				t.Fatalf("commit: [ %v ], Want at least one error, got nil", i)
			}

			if !slices.Contains(got, want) {
				t.Fatalf("commit [ %v ], want [ %v ], but got none", i, want)
			}
		}
	})

	t.Run("Invalid commit message", func(t *testing.T) {
		commit := `This is some commit

This is the commit description
this has some content: not a footer
but nothing on this commit follows
the convention`

		m := parseCommitMessage(commit)
		got := lintCommitMessage(m)
		if got == nil {
			t.Fatal("Expected errors, got nil")
		}

		want := []error{
			errInvalidHeaderFormat,
			errInvalidBodyBlankLines,
			errInvalidFooterSpace,
			errInvalidFooterFormat,
		}

		for _, err := range want {
			if !slices.Contains(got, err) {
				t.Fatalf("want [ %v ], got %v", want, got)
			}
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
