package platform_git

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"idp-cfs/platform_gp"
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

func (c *GitCode) PushFiles(url string, branch string, gpPath string, relativePath string) error {

	var files []string
	absolutePath := path.Join(gpPath, relativePath)

	err := filepath.Walk(absolutePath, func(filepath string, info os.FileInfo, err error) error {

		files = append(files, filepath)
		return nil
	})
	if err != nil {
		log.Error().Msgf("Failed to walk the directory %s. Error: %v", gpPath, err)
	}

	log.Info().Msgf("List of files to commit: %+v", files)

	var pat = ""
	var user = ""

	if c.githubCode != nil {
		pat = os.Getenv("CFS_GITHUB_PAT")
		user = os.Getenv("CFS_GITHUB_USER")
	}

	newCodeRepo, err := git.PlainClone(platform_gp.GetCheckinPath(), false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: user,
			Password: pat,
		},
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})

	if err != nil {
		log.Error().Msgf("Failed to clone the repo. Error: %v", err)
		return err
	}

	w, err := newCodeRepo.Worktree()

	if err != nil {
		log.Error().Msgf("Failed to return worktree of the repo. Error: %v", err)
		return err
	}

	for _, file := range files {
		_, err := w.Add(file)
		if err != nil {
			log.Error().Msgf("Failed to add file %s to commit. Error: %v", file, err)
		}
	}

	status, err := w.Status()

	if err != nil {
		log.Error().Msgf("Failed to get status. Error: %v", err)
		return err
	}

	log.Info().Msgf("Status of the repo before commit push: %+v", status)

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

	commitObj, err := newCodeRepo.CommitObject(commit)
	if err != nil {
		log.Error().Msgf("Failed for repo to commit. Error: %v", err)
	}

	log.Info().Msgf("Repo Commit. %v", commitObj)

	err = newCodeRepo.Push(&git.PushOptions{})
	if err != nil {
		log.Error().Msgf("Failed for push commit changes. Error: %v", err)
	}

	return errors.New("not implemented yet")
}
