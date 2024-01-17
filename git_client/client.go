package git_client

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"idp-cfs/global"
	"time"
)

func GetGitClient() *GitClient {
	return &GitClient{
		plainCloneFunc: func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return git.PlainClone(path, isBare, o)
		},
		checkoutFunc: func(w *git.Worktree, opts *git.CheckoutOptions) error {
			return w.Checkout(opts)
		},
		headFunc: func(r *git.Repository) (*plumbing.Reference, error) {
			return r.Head()
		},
		workTreeFunc: func(r *git.Repository) (*git.Worktree, error) {
			return r.Worktree()
		},
		statusFunc: func(w *git.Worktree) (git.Status, error) {
			return w.Status()
		},
		addGlobFunc: func(w *git.Worktree, glob string) error {
			return w.AddGlob(glob)
		},
		commitFunc: func(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error) {
			return w.Commit(msg, opts)
		},
		pushFunc: func(r *git.Repository, o *git.PushOptions) error {
			return r.Push(o)
		},
		referencesFunc: func(r *git.Repository) (storer.ReferenceIter, error) {
			return r.References()
		},
	}
}

func (g *GitClient) CloneRepository(path, gitUrl string, branch string, auth *GitClientAuth) (*git.Repository, error) {

	err := global.CreateFolder(path)
	if err != nil {
		return nil, global.LogError(fmt.Sprintf("failed to create folder at path %s: %v", path, err))
	}

	var r *git.Repository
	branchRefStr := fmt.Sprintf("refs/heads/%s", branch)

	if auth != nil {

		log.Info().Msg("Cloning public repo with auth")

		r, err = g.plainCloneFunc(
			path,
			false,
			&git.CloneOptions{
				URL:           gitUrl,
				ReferenceName: plumbing.ReferenceName(branchRefStr),
				Progress:      log.Logger,
				Auth: &http.BasicAuth{
					Username: auth.User,
					Password: auth.Token,
				},
			})
	} else {

		log.Info().Msg("Cloning private repo without any auth")

		r, err = g.plainCloneFunc(path,
			false,
			&git.CloneOptions{
				URL:           gitUrl,
				ReferenceName: plumbing.ReferenceName(branchRefStr),
				Progress:      log.Logger,
			})
	}

	if err != nil {
		return nil, global.LogError(fmt.Sprintf("failed to clone repo: %v", err))
	}

	headRef, err := g.headFunc(r)
	if err != nil {
		return nil, global.LogError(fmt.Sprintf("failed to return HEAD reference: %v", err))
	}

	log.Info().Msgf("Cloned the git repo at %s. HEAD ref: %s", path, headRef)

	branchRef := g.getRefForBranch(r, branchRefStr)

	log.Info().Msgf("Found branch with ref %+v", branchRef)

	w, err := g.workTreeFunc(r)
	if err != nil {
		return nil, global.LogError(fmt.Sprintf("failed to get worktree for repo: %v", err))
	}

	err = g.checkoutFunc(w, &git.CheckoutOptions{
		Branch: branchRef.Name(),
		Create: false,
	})

	if err != nil {
		return nil, global.LogError(fmt.Sprintf("failed to checkout branch: %v", err))
	}

	return r, nil
}

func (g *GitClient) PushFiles(repo *git.Repository, auth *GitClientAuth) error {

	_, err := g.headFunc(repo)
	if err != nil {
		return global.LogError(fmt.Sprintf("failed to return HEAD: %v", err))
	}

	w, err := g.workTreeFunc(repo)
	if err != nil {
		return global.LogError(fmt.Sprintf("failed to return worktree: %v", err))
	}

	err = g.addGlobFunc(w, ".")
	if err != nil {
		return global.LogError(fmt.Sprintf("failed to add . to git: %v", err))
	}

	_, err = g.statusFunc(w)
	if err != nil {
		return global.LogError(fmt.Sprintf("failed to get status: %v", err))
	}

	commit, err := g.commitFunc(w, "Adding GP as per idp-cfs contract", &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.gitConfigUser,
			Email: g.gitConfigEmail,
			When:  time.Now(),
		},
	})

	if err != nil {
		return global.LogError(fmt.Sprintf("failed to commit: %v", err))
	}

	log.Info().Msgf("Files Commit. %v", commit)

	err = g.pushFunc(repo, &git.PushOptions{
		Auth: &http.BasicAuth{
			Username: auth.User,
			Password: auth.Token,
		},
	})
	if err != nil {
		return global.LogError(fmt.Sprintf("failed for push commit changes: %v", err))
	}

	return nil
}

func GetAuth(user, token string) *GitClientAuth {

	if user == "" && token == "" {
		log.Warn().Msg("Git getAuth - No basic auth env defined.")
		return nil
	} else {
		log.Info().Msg("Git getAuth - Basic auth env defined")
		return &GitClientAuth{
			User:  user,
			Token: token,
		}
	}
}

func (g *GitClient) getRefForBranch(r *git.Repository, branchName string) *plumbing.Reference {
	var res *plumbing.Reference

	if r == nil {
		log.Error().Msg("Repository is nil.")
		return nil
	}

	refs, _ := g.referencesFunc(r)
	err := refs.ForEach(func(ref *plumbing.Reference) error {

		if ref.Type() == plumbing.HashReference && ref.Name().String() == branchName {
			log.Info().Msgf(" - Ref Found for branch: %+v", ref)
			res = ref
		}

		return nil
	})
	if err != nil {
		log.Error().Msgf("Error going through git refs: %v", err)
	}

	return res
}
