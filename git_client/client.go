package git_client

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"idp-cfs2/util"
	"os"
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

	err := util.CreateFolder(path)
	if err != nil {
		return nil, err
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
		log.Error().Msgf("Failed to clone the repo: %v", err)
		return nil, err
	}

	headRef, err := g.headFunc(r)
	if err != nil {
		log.Error().Msgf("Unable to return reference of HEAD: %v", err)
		return nil, err
	}

	log.Info().Msgf("Cloned the git repo at %s. HEAD ref: %s", path, headRef)

	branchRef := g.getRefForBranch(r, branchRefStr)

	log.Info().Msgf("Found branch with ref %+v", branchRef)

	w, err := g.workTreeFunc(r)
	if err != nil {
		log.Error().Msgf("Error trying to get worktree for repository: %v", err)
		return nil, err
	}

	err = g.checkoutFunc(w, &git.CheckoutOptions{
		Branch: branchRef.Name(),
		Create: false,
	})

	if err != nil {
		log.Error().Msgf("Error trying to checkout the branch: %v", err)
		return nil, err
	}

	return r, nil
}

func (g *GitClient) PushFiles(repo *git.Repository, localDir string, auth *GitClientAuth) error {

	_, err := g.headFunc(repo)
	if err != nil {
		log.Error().Msgf("Failed to return HEAD: %v", err)
		return err
	}

	w, err := g.workTreeFunc(repo)
	if err != nil {
		log.Error().Msgf("Failed to return worktree: %v", err)
		return err
	}

	err = os.Chdir(localDir)
	if err != nil {
		log.Error().Msgf("Failed to change directory into %s: %v", localDir, err)
		return err
	}

	err = g.addGlobFunc(w, ".")
	if err != nil {
		log.Error().Msgf("Failed to add . to git: %v", err)
		return err
	}

	_, err = g.statusFunc(w)
	if err != nil {
		log.Error().Msgf("Failed to get status: %v", err)
		return err
	}

	commit, err := g.commitFunc(w, "Adding GP as per idp-cfs contract", &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.gitConfigUser,
			Email: g.gitConfigEmail,
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Error().Msgf("Failed to commit: %v", err)
		return err
	}

	log.Info().Msgf("Files Commit. %v", commit)

	err = g.pushFunc(repo, &git.PushOptions{
		Auth: &http.BasicAuth{
			Username: auth.User,
			Password: auth.Token,
		},
	})
	if err != nil {
		log.Error().Msgf("Failed for push commit changes: %v", err)
		return err
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
