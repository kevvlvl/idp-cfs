package client_git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type CfsGit interface {
	PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	Checkout(w *git.Worktree, opts *git.CheckoutOptions) error
	Head(r *git.Repository) (*plumbing.Reference, error)
	getRefForBranch(r *git.Repository, branchName string) *plumbing.Reference
	WorkTree(r *git.Repository) (*git.Worktree, error)
	Status(w *git.Worktree) (git.Status, error)
	AddGlob(w *git.Worktree, glob string) error
	Commit(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error)
	Push(r *git.Repository, o *git.PushOptions) error
}

type GitClient struct{}

type GitBasicAuth struct {
	user  string
	token string
}

const (
	// CodeGithub for the code repository of type Github (public/cloud)
	CodeGithub = "github"
	// CodeGitlab for the code repository of type Gitlab
	CodeGitlab = "gitlab"
	// CodeGitea for the code repository of type Gitea
	CodeGitea = "gitea"
)
