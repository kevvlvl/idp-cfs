package git_api

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/git_client"
	"idp-cfs/global"
	"path"
)

func pushGoldenPath(tool, codeUrl, codeDefaultBranch, url, pathDir, branch, gpWorkdir, codeWorkDir string, tag *string) error {
	auth := getAuth(tool)
	git := git_client.GetGitClient()
	gitCodeAuth := git_client.GetAuth(auth.codeUser, auth.codeToken)
	gitGpAuth := git_client.GetAuth(auth.gpUser, auth.gpToken)

	if git != nil {

		// Clone the git repo
		codeRepo, err := git.CloneRepository(codeWorkDir, codeUrl, codeDefaultBranch, gitCodeAuth)
		if err != nil {
			return err
		}

		// Clone the gp repo
		_, err = git.CloneRepository(gpWorkdir, url, branch, gitGpAuth)
		gpPath := path.Join(gpWorkdir, pathDir)

		err = global.CopyFilesDeep(gpPath, codeWorkDir)
		if err != nil {
			log.Error().Msgf("Failed to copy files from goldenpath to code repo: %v", err)
			return err
		}

		err = git.PushFiles(codeRepo, gitCodeAuth)
		if err != nil {
			log.Error().Msgf("Failed to push files to the code repo: %v", err)
			return err
		}
	}

	return nil
}
