package client_git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type MockGitClient struct {
	PlainCloneFunc      func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	CheckoutFunc        func(w *git.Worktree, opts *git.CheckoutOptions) error
	HeadFunc            func(r *git.Repository) (*plumbing.Reference, error)
	getRefForBranchFunc func(r *git.Repository, branchName string) *plumbing.Reference
	WorkTreeFunc        func(r *git.Repository) (*git.Worktree, error)
	StatusFunc          func(w *git.Worktree) (git.Status, error)
	AddGlobFunc         func(w *git.Worktree, glob string) error
	CommitFunc          func(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error)
	PushFunc            func(r *git.Repository, o *git.PushOptions) error
}

func (g *MockGitClient) PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	return g.PlainCloneFunc(path, isBare, o)
}

func (g *MockGitClient) Checkout(w *git.Worktree, opts *git.CheckoutOptions) error {
	return g.CheckoutFunc(w, opts)
}

func (g *MockGitClient) Head(r *git.Repository) (*plumbing.Reference, error) {
	return g.HeadFunc(r)
}

func (g *MockGitClient) getRefForBranch(r *git.Repository, branchName string) *plumbing.Reference {
	return g.getRefForBranchFunc(r, branchName)
}

func (g *MockGitClient) WorkTree(r *git.Repository) (*git.Worktree, error) {
	return g.WorkTreeFunc(r)
}

func (g *MockGitClient) Status(w *git.Worktree) (git.Status, error) {
	return g.StatusFunc(w)
}

func (g *MockGitClient) AddGlob(w *git.Worktree, glob string) error {
	return g.AddGlobFunc(w, glob)
}

func (g *MockGitClient) Commit(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error) {
	return g.CommitFunc(w, msg, opts)
}

func (g *MockGitClient) Push(r *git.Repository, o *git.PushOptions) error {
	return g.PushFunc(r, o)
}
