package git_client

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type GitClient struct {
	gitConfigUser       string
	gitConfigEmail      string
	plainCloneFunc      func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	checkoutFunc        func(w *git.Worktree, opts *git.CheckoutOptions) error
	headFunc            func(r *git.Repository) (*plumbing.Reference, error)
	getRefForBranchFunc func(r *git.Repository, branchName string) *plumbing.Reference
	workTreeFunc        func(r *git.Repository) (*git.Worktree, error)
	statusFunc          func(w *git.Worktree) (git.Status, error)
	addGlobFunc         func(w *git.Worktree, glob string) error
	commitFunc          func(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error)
	pushFunc            func(r *git.Repository, o *git.PushOptions) error
	referencesFunc      func(r *git.Repository) (storer.ReferenceIter, error)
}

type GitClientAuth struct {
	User  string
	Pass  string
	Token string
}
