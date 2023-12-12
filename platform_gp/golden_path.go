package platform_gp

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"
)

func GetGoldenPath(url string, branch string, path string, tag string, gpCheckoutPath string) GoldenPath {

	var tool string

	// if github, set tool to github
	if strings.Contains(url, "github.com") {
		tool = GpGithub
	}

	return GoldenPath{
		Tool:           tool,
		URL:            url,
		Branch:         branch,
		Path:           path,
		Tag:            tag,
		GpCheckoutPath: gpCheckoutPath,
	}
}

func (gp *GoldenPath) CloneGp() error {

	err := gp.DeleteClonePathDir()
	if err != nil {

		log.Error().Msgf("Error cleaning up the folder %s. Error = %v: ", gp.GpCheckoutPath, err)
		return failedCloneGpError()
	}

	r, err := git.PlainClone(gp.GpCheckoutPath, false, &git.CloneOptions{
		URL:      gp.URL,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Error().Msgf("Error trying to clone the gp URL %v: %v", gp.URL, err)
	}

	headRef, err := r.Head()
	if err != nil {
		log.Error().Msgf("Unable to return reference of HEAD: %v", err)
		return failedCloneGpError()
	}

	log.Info().Msgf("Cloned the golden path at %s. HEAD ref: %s", gp.GpCheckoutPath, headRef)
	gp.repository = r

	branchRef := getRefForBranchName(r, fmt.Sprintf("refs/remotes/origin/%s", gp.Branch))

	log.Info().Msgf("Found the branch with ref %+v", branchRef)

	worktree, err := r.Worktree()
	if err != nil {
		log.Error().Msgf("Error trying to get worktree for the repository: %v", err)
		return failedCloneGpError()
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef.Name(),
		Create: false,
	})

	if err != nil {
		log.Error().Msgf("Error trying to checkout the branch: %v", err)
		return failedCloneGpError()
	}

	if _, err := os.Stat(path.Join(gp.GpCheckoutPath, gp.Path)); !os.IsNotExist(err) {
		log.Info().Msgf("Succesfully verified that path %v exists in the cloned repo", gp.Path)
	} else {
		log.Error().Msgf("Failed to find the the path %v in the cloned repo: %v", gp.Path, err)
	}

	return nil
}

func (gp *GoldenPath) DeleteClonePathDir() error {
	return os.RemoveAll(gp.GpCheckoutPath)
}

// showRefsFound outputs all found Refs for the git repository in input
func getRefForBranchName(r *git.Repository, branchName string) *plumbing.Reference {
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

func failedCloneGpError() error {
	return errors.New("failed to clone the GoldenPath")
}
