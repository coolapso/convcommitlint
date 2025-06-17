package cmd

import (
	"slices"
	"testing"
)

func TestLintFooter(t *testing.T) {
	t.Run("valid trailer", func(t *testing.T) {
		footer := []string{"key: value"}
		got := lintFooter(footer)
		if got != nil {
			t.Fatal("got error expected nil")
		}
	})

	t.Run("valid trailer with multiple keys", func(t *testing.T) {
		footer := []string{"key: value", "key1: value1", "key2: value2"}
		got := lintFooter(footer)
		if got != nil {
			t.Fatal("got error expected nil")
		}
	})

	t.Run("invalid trailer without space", func(t *testing.T) {
		footer := []string{"key:value"}
		got := lintFooter(footer)
		if got == nil {
			t.Fatal("want error, got nil")
		}
	})

	t.Run("invalid trailer without colon", func(t *testing.T) {
		footer := []string{"key value"}
		got := lintFooter(footer)
		if got == nil {
			t.Fatal("want error, got nil")
		}
	})

	t.Run("invalid trailer with new line", func(t *testing.T) {
		footer := []string{"key: value", "\n", "key1: value1"}
		want := errInvalidFooterEmptyLine
		got := lintFooter(footer)
		if got == nil {
			t.Fatal("want error, got nil")
		}

		if !slices.Contains(got, want) {
			t.Fatalf("error message not in list, want: [ %v ], got %v", want, got)
		}
	})

	t.Run("invalid trailer with empty line", func(t *testing.T) {
		footer := []string{"key: value", "", "key1: value1"}
		want := errInvalidFooterEmptyLine
		got := lintFooter(footer)
		if got == nil {
			t.Fatal("want error, got nil")
		}

		if !slices.Contains(got, want) {
			t.Fatalf("error message not in list, want: [ %v ], got %v", want, got)
		}
	})

	t.Run("valid trailer with long key", func(t *testing.T) {
		footer := []string{"key: value", "key1: value1", " continuation value 1"}
		got := lintFooter(footer)
		if got != nil {
			t.Fatal("did not want errors, got: ", got)
		}
	})

	t.Run("valid trailer with new line at end", func(t *testing.T) {
		footer := []string{"key: value", "key1: value1", " continuation value 1", "key2: value2", "\n"}
		got := lintFooter(footer)
		if got != nil {
			t.Fatal("did not want errors, got: ", got)
		}
	})

	t.Run("invalid keys with spaces", func(t *testing.T) {
		trailer := []string{"key pair: value"}
		want := errInvalidFooterSpace
		got := lintFooter(trailer)
		if got == nil {
			t.Fatal("Want error got nil")
		}

		if !slices.Contains(got, want) {
			t.Fatalf("error not in list, want [ %v ], got %v", want, got)
		}
	})

	t.Run("invalid breaking change", func(t *testing.T) {
		trailer := []string{"breaking change: value"}
		got := lintFooter(trailer)
		if got == nil {
			t.Fatal("Want error got nil")
		}
	})

	t.Run("valid breaking change", func(t *testing.T) {
		trailer := []string{"BREAKING CHANGE: value"}
		got := lintFooter(trailer)
		if got != nil {
			t.Fatal("want nil got:", got)
		}
	})

	t.Run("typo breaking changes", func(t *testing.T) {
		trailer := []string{"BREAKING CHANGES: value"}
		want := errInvalidFooterBreakingTypo
		got := lintFooter(trailer)
		if got == nil {
			t.Fatal("want error got nil")
		}

		if !slices.Contains(got, want) {
			t.Fatalf("error not in list, want [ %v ], got: %v", want, got)
		}
	})
}
