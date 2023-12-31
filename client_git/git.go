package client_git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"idp-cfs/util"
	"os"
	"time"
)

func GetGitClient() *GitClient {
	return &GitClient{}
}

func (g *GitClient) CloneRepository(path string, gitUrl string, branch *string, auth *GitBasicAuth) (*git.Repository, error) {

	err := util.CreateFolder(path)
	if err != nil {
		return nil, err
	}

	var r *git.Repository
	branchRefStr := fmt.Sprintf("refs/heads/%s", *branch)

	if auth != nil {

		log.Info().Msg("Cloning public repo with auth")

		r, err = g.PlainClone(
			path,
			false,
			&git.CloneOptions{
				URL:           gitUrl,
				ReferenceName: plumbing.ReferenceName(branchRefStr),
				Progress:      log.Logger,
				Auth: &http.BasicAuth{
					Username: auth.user,
					Password: auth.token,
				},
			})
	} else {

		log.Info().Msg("Cloning private repo without any auth")

		r, err = g.PlainClone(path,
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

	headRef, err := g.Head(r)
	if err != nil {
		log.Error().Msgf("Unable to return reference of HEAD: %v", err)
		return nil, err
	}

	if branch != nil {
		log.Info().Msgf("Cloned the git repo at %s. HEAD ref: %s", path, headRef)

		branchRef := g.getRefForBranch(r, branchRefStr)

		log.Info().Msgf("Found branch with ref %+v", branchRef)

		w, err := g.WorkTree(r)
		if err != nil {
			log.Error().Msgf("Error trying to get worktree for repository: %v", err)
			return nil, err
		}

		err = g.Checkout(w, &git.CheckoutOptions{
			Branch: branchRef.Name(),
			Create: false,
		})

		if err != nil {
			log.Error().Msgf("Error trying to checkout the branch: %v", err)
			return nil, err
		}
	}

	return r, nil
}

func (g *GitClient) PushFiles(repo *git.Repository, localDir string, auth *GitBasicAuth) error {

	_, err := g.Head(repo)
	if err != nil {
		log.Error().Msgf("Failed to return HEAD: %v", err)
		return err
	}

	w, err := g.WorkTree(repo)
	if err != nil {
		log.Error().Msgf("Failed to return worktree: %v", err)
		return err
	}

	err = os.Chdir(localDir)
	if err != nil {
		log.Error().Msgf("Failed to change directory into %s: %v", localDir, err)
		return err
	}

	err = g.AddGlob(w, ".")
	if err != nil {
		log.Error().Msgf("Failed to add . to git: %v", err)
		return err
	}

	_, err = g.Status(w)
	if err != nil {
		log.Error().Msgf("Failed to get status: %v", err)
		return err
	}

	// FIXME refactor parameters. externalize author into app config
	commit, err := g.Commit(w, "Adding GP as per idp-cfs contract", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "idp-cfs",
			Email: "idp-cfs@kevvlvl.github.noreply.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Error().Msgf("Failed to commit: %v", err)
		return err
	}

	log.Info().Msgf("Files Commit. %v", commit)

	err = g.Push(repo, &git.PushOptions{
		Auth: &http.BasicAuth{
			Username: auth.user,
			Password: auth.token,
		},
	})
	if err != nil {
		log.Error().Msgf("Failed for push commit changes: %v", err)
		return err
	}

	return nil
}

func GetAuth(user string, token string) *GitBasicAuth {

	if user == "" && token == "" {
		log.Debug().Msg("getAuth - No basic auth env defined.")
		return nil
	} else {
		log.Debug().Msg("getAuth - Basic auth env defined")
		return &GitBasicAuth{
			user:  user,
			token: token,
		}
	}
}

func (g *GitClient) getRefForBranch(r *git.Repository, branchName string) *plumbing.Reference {
	var res *plumbing.Reference

	refs, _ := r.References()
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

func (g *GitClient) PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, false, o)
}

func (g *GitClient) Checkout(w *git.Worktree, opts *git.CheckoutOptions) error {
	return w.Checkout(opts)
}

func (g *GitClient) Head(r *git.Repository) (*plumbing.Reference, error) {
	return r.Head()
}

func (g *GitClient) WorkTree(r *git.Repository) (*git.Worktree, error) {
	return r.Worktree()
}

func (g *GitClient) Status(w *git.Worktree) (git.Status, error) {
	return w.Status()
}

func (g *GitClient) AddGlob(w *git.Worktree, glob string) error {
	return w.AddGlob(glob)
}

func (g *GitClient) Commit(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error) {
	return w.Commit(msg, opts)
}

func (g *GitClient) Push(r *git.Repository, o *git.PushOptions) error {
	return r.Push(o)
}
