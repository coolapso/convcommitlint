package cmd

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func TestCommitsSinceBaseExcludesSharedHistoryWhenBranchesDiverge(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		t.Fatal(err)
	}

	root := writeTestCommit(t, repo, "chore: root")
	base := writeTestCommit(t, repo, "invalid base commit", root)
	feature := writeTestCommit(t, repo, "fix: feature commit", root)
	head := writeTestCommit(t, repo, "feat: another feature commit", feature)

	commits, err := commitsSinceBase(repo, head, base)
	if err != nil {
		t.Fatal(err)
	}

	got := make(map[plumbing.Hash]bool)
	for _, commit := range commits {
		got[commit.Hash] = true
	}

	for _, hash := range []plumbing.Hash{head, feature} {
		if !got[hash] {
			t.Errorf("expected commit %s in range", hash)
		}
	}
	for _, hash := range []plumbing.Hash{base, root} {
		if got[hash] {
			t.Errorf("did not expect commit %s in range", hash)
		}
	}
}

func writeTestCommit(t *testing.T, repo *git.Repository, message string, parents ...plumbing.Hash) plumbing.Hash {
	t.Helper()

	commit := &object.Commit{
		Author:       object.Signature{Name: "Test", Email: "test@example.com", When: time.Unix(0, 0)},
		Committer:    object.Signature{Name: "Test", Email: "test@example.com", When: time.Unix(0, 0)},
		Message:      message,
		TreeHash:     plumbing.ZeroHash,
		ParentHashes: parents,
	}
	encoded := repo.Storer.NewEncodedObject()
	if err := commit.Encode(encoded); err != nil {
		t.Fatal(err)
	}
	hash, err := repo.Storer.SetEncodedObject(encoded)
	if err != nil {
		t.Fatal(err)
	}

	return hash
}
