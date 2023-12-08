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
	"idp-cfs/platform_gp"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

func GetCode(tool string) *GitCode {

	var code *GitCode

	switch tool {

	case CodeGithub:
		code = &GitCode{githubCode: GetGithubCode()}
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

func (c *GitCode) PushFiles(url string, branch string, relativePath string) error {

	// FIXME refactor all the disgusting error handling.

	branchRef := fmt.Sprintf("refs/heads/%s", branch)
	gpPath := platform_gp.GetCheckoutPath()
	err := deleteCodePath()

	// FIXME refactor error handling
	os.Mkdir(getCodePath(), 0775)
	if err != nil {
		log.Error().Msgf("Failed to cleanup the code path. Error: %v", err)
	}

	var pat = ""
	var user = ""

	if c.githubCode != nil {
		pat = os.Getenv("CFS_GITHUB_PAT")
		user = os.Getenv("CFS_GITHUB_USER")
	}

	r, err := git.PlainClone(getCodePath(), false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: user,
			Password: pat,
		},
		ReferenceName: plumbing.ReferenceName(branchRef),
		Progress:      os.Stdout,
	})

	if err != nil && err.Error() != "remote repository is empty" {
		log.Error().Msgf("Failed to clone the repo. Error: %v", err)
		return err
	}

	_, err = r.Head()
	if err != nil {
		log.Error().Msgf("Failed to return HEAD. Error: %v", err)
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		log.Error().Msgf("Failed to return worktree. Error: %v", err)
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRef),
	})

	if err != nil {
		log.Error().Msgf("Failed to checkout on a new branch. Error: %v", err)
		return err
	}

	absolutePath := path.Join(gpPath, relativePath)

	err = filepath.Walk(absolutePath, func(file string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// FIXME refactor
			srcFile, _ := os.Open(file)
			defer srcFile.Close()

			destFilePath := filepath.Join(getCodePath(), info.Name())
			destFile, _ := os.Create(destFilePath)
			defer destFile.Close()

			_, err := io.Copy(destFile, srcFile)
			if err != nil {
				log.Error().Msgf("Failed to copy the file from gp path to the new code path. Error: %v", err)
			}

			// FIXME refactor - make sure working dir is code path!
			_, err = w.Add(info.Name())
			if err != nil {
				log.Error().Msgf("Failed to add file %s to commit. Error: %v", file, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Msgf("Failed to walk the directory %s. Error: %v", gpPath, err)
	}

	_, err = w.Status()
	if err != nil {
		log.Error().Msgf("Failed to get status. Error: %v", err)
		return err
	}

	// FIXME refactor parameters. externalize author into app config
	commit, err := w.Commit("Adding GP as per idp-cfs contract", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "idp-cfs",
			Email: "idp-cfs@kevvlvl.github.noreply.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Error().Msgf("Failed to commit. Error: %v", err)
		return err
	}

	log.Info().Msgf("Files Commit. %v", commit)

	// FIXME refactor Refspec
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)),
		},
		Auth: &http.BasicAuth{
			Username: user,
			Password: pat,
		},
	})
	if err != nil {
		log.Error().Msgf("Failed for push commit changes. Error: %v", err)
	}

	return errors.New("not implemented yet")
}

func getCodePath() string {

	checkoutPath := os.Getenv("CFS_GP_CODE_CLONE_PATH")

	if checkoutPath == "" {
		checkoutPath = "/tmp/idp-cfs-code"
	}

	return checkoutPath

}

func deleteCodePath() error {
	return os.RemoveAll(getCodePath())
}
