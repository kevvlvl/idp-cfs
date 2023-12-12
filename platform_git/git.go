package platform_git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

func GetCode(tool string, codeClonePath string) *GitCode {

	var code *GitCode

	switch tool {

	case CodeGithub:
		code = &GitCode{
			githubCode:    GetGithubCode(),
			CodeClonePath: codeClonePath,
		}
	case CodeGitea:
		log.Warn().Msg("Not implemented yet!")
	case CodeGitlab:
		log.Warn().Msg("Not implemented yet!")
	default:
		log.Error().Msgf("Unexpected code tool system which somehow passed validation! Tool: %v", tool)

	}

	return code
}

func (c *GitCode) GetOrganization(org string) (*Organization, error) {

	if c.githubCode != nil {
		return c.githubCode.GetOrganization(org)
	} else {
		return nil, errors.New("not implemented yet")
	}
}

func (c *GitCode) GetRepository(repo string) (*Repository, error) {

	if c.githubCode != nil {

		r, err := c.githubCode.GetRepository(repo)
		c.Repository = r

		return r, err
	} else {
		return nil, errors.New("not implemented yet")
	}
}

func (c *GitCode) CreateRepository(repoName string, branch string) (*Repository, error) {
	if c.githubCode != nil {
		return c.githubCode.CreateRepository(repoName, branch)
	} else {
		return nil, errors.New("not implemented yet")
	}
}

func (c *GitCode) PushFiles(url string, branch string, relativePath string, gpCheckoutPath string) error {

	branchRef := fmt.Sprintf("refs/heads/%s", branch)

	err := c.DeleteCodePath()
	if err != nil {
		log.Error().Msgf("Failed to cleanup the code path: %v", err)
		return err
	}

	err = os.Mkdir(c.CodeClonePath, 0775)
	if err != nil {
		log.Error().Msgf("Failed to create the directory for %s: %v", c.CodeClonePath, err)
		return err
	}

	var pat = ""
	var user = ""

	if c.githubCode != nil {
		pat = GetPersonalAccessToken()
		user = GetUsername()
	}

	r, err := git.PlainClone(c.CodeClonePath, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: user,
			Password: pat,
		},
		ReferenceName: plumbing.ReferenceName(branchRef),
		Progress:      os.Stdout,
	})

	if err != nil && err.Error() != "remote repository is empty" {
		log.Error().Msgf("Failed to clone the repo: %v", err)
		return err
	}

	_, err = r.Head()
	if err != nil {
		log.Error().Msgf("Failed to return HEAD: %v", err)
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		log.Error().Msgf("Failed to return worktree: %v", err)
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRef),
	})

	if err != nil {
		log.Error().Msgf("Failed to checkout on a new branch: %v", err)
		return err
	}

	err = os.Chdir(c.CodeClonePath)
	if err != nil {
		log.Error().Msgf("Failed to change directory into %s: %v", c.CodeClonePath, err)
		return err
	}

	err = filepath.Walk(path.Join(gpCheckoutPath, relativePath), func(file string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			srcFile, _ := os.Open(file)
			defer func(srcFile *os.File) {
				err := srcFile.Close()
				if err != nil {
					log.Error().Msgf("Failed to close the src file %s: %v", srcFile.Name(), err)
				}
			}(srcFile)

			destFilePath := filepath.Join(c.CodeClonePath, info.Name())
			destFile, _ := os.Create(destFilePath)
			defer func(destFile *os.File) {
				err := destFile.Close()
				if err != nil {
					log.Error().Msgf("Failed to close the src file %s: %v", srcFile.Name(), err)
				}
			}(destFile)

			_, err := io.Copy(destFile, srcFile)
			if err != nil {
				log.Error().Msgf("Failed to copy the file from gp path to the new code path: %v", err)
				return err
			}

			_, err = w.Add(info.Name())
			if err != nil {
				log.Error().Msgf("Failed to add file %s to commit: %v", file, err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Msgf("Failed to walk the directory %s: %v", gpCheckoutPath, err)
		return err
	}

	_, err = w.Status()
	if err != nil {
		log.Error().Msgf("Failed to get status: %v", err)
		return err
	}

	commit, err := w.Commit("Adding GP as per idp-cfs contract", &git.CommitOptions{
		Author: &object.Signature{
			Name:  GitCommitAuthor,
			Email: GitCommitAuthorEmail,
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Error().Msgf("Failed to commit: %v", err)
		return err
	}

	log.Info().Msgf("Files Commit. %v", commit)

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("%s:%s", branchRef, branchRef)),
		},
		Auth: &http.BasicAuth{
			Username: user,
			Password: pat,
		},
	})
	if err != nil {
		log.Error().Msgf("Failed for push commit changes: %v", err)
		return err
	}

	return nil
}

func (c *GitCode) DeleteCodePath() error {
	return os.RemoveAll(c.CodeClonePath)
}
