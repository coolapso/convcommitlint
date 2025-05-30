package cmd

import (
	"slices"
	"testing"
)

func TestParseHeader(t *testing.T) {
	t.Run("test valid header", func(t *testing.T) {
		want := header{
			text:            string("feat(scope)!: description"),
			commitType:      string("feat"),
			scope:           string("scope"),
			scopeParentesis: bool(true),
			bang:            bool(true),
			separator:       bool(true),
			description:     string("description"),
		}

		got := parseHeader("feat(scope)!: description")
		if got != want {
			t.Fatalf("want %v, got %v", want, got)
		}
	})
}

func TestLintHeader(t *testing.T) {
	t.Run("valid full header", func(t *testing.T) {
		h := header{
			text:            string("feat(scope)!: description"),
			commitType:      string("feat"),
			scope:           string("scope"),
			scopeParentesis: bool(true),
			bang:            bool(true),
			separator:       bool(true),
			description:     string("description"),
		}

		got := lintHeader(h)
		if got != nil {
			t.Fatal("want valid header, got: ", got)
		}
	})

	t.Run("valid header no scope", func(t *testing.T) {
		h := header{
			text:        string("feat: description"),
			commitType:  string("feat"),
			separator:   bool(true),
			description: string("description"),
		}
		got := lintHeader(h)
		if got != nil {
			t.Fatal("want valid header, got: ", got)
		}
	})

	t.Run("valid header no scope with bang", func(t *testing.T) {
		h := header{
			text:        string("feat!: description"),
			commitType:  string("feat"),
			bang:        bool(true),
			separator:   bool(true),
			description: string("description"),
		}
		got := lintHeader(h)
		if got != nil {
			t.Fatal("want valid header, got: ", got)
		}
	})

	t.Run("invalid header", func(t *testing.T) {
		headers := []header{
			{bang: true},
			{separator: true},
			{scopeParentesis: true},
			{bang: true, separator: true},
			{bang: true, scopeParentesis: true},
			{bang: true, scopeParentesis: true, separator: true},
			{
				text:            "feat(): description",
				commitType:      "feat",
				scopeParentesis: true,
				separator:       true,
				description:     "description",
			},
		}
		want := errInvalidHeaderFormat

		for _, h := range headers {
			got := lintHeader(h)
			if got == nil {
				t.Fatalf("header: %v, Expected invalid header error got nil", h)
			}

			if !slices.Contains(got, want) {
				t.Fatalf("header: [ %v ], error not in list want [ %v ], got: %v ", h, want, got)
			}
		}
	})

	t.Run("invalid long description", func(t *testing.T) {
		h := header{
			text:        string("feat(scope): this is a very long description, and error should be returned, otherwise somethings is wrong"),
			commitType:  "feat",
			scope:       "scope",
			separator:   true,
			description: string("this is a very long description, and error should be returned, otherwise somethings is wrong"),
		}
		got := lintHeader(h)
		if got == nil {
			t.Fatal("Expected invalid hearror header got nil")
		}

		if !slices.Contains(got, errLongHeader) {
			t.Fatalf("error not in list, want [ %v ], got %v", errLongHeader, got)
		}
	})

	t.Run("typo header feat", func(t *testing.T) {
		headers := []header{
			{text: "feta(scope): description", commitType: "feta", scope: "scope", separator: true, description: "description"},
			{text: "featur(scope): description", commitType: "featur", scope: "scope", separator: true, description: "description"},
		}
		want := errTypoHeaderFeat

		for _, h := range headers {
			got := lintHeader(h)
			if got == nil {
				t.Fatal("want error got nil")
			}

			if !slices.Contains(got, want) {
				t.Fatalf("header: %v, want error [ %v ], got %v", h, want, got)
			}
		}
	})

	t.Run("typo header fix", func(t *testing.T) {
		headers := []header{
			{text: "fxi(scope): description", commitType: "fxi", scope: "scope", separator: true, description: "description"},
			{text: "dix(scope): description", commitType: "dix", scope: "scope", separator: true, description: "description"},
		}
		want := errTypoHeaderFix

		for _, h := range headers {
			got := lintHeader(h)
			if got == nil {
				t.Fatal("want error got nil")
			}

			if !slices.Contains(got, want) {
				t.Fatalf("header: %v, want error: [ %v ], got: %v", h, want, got)
			}
		}
	})
}
