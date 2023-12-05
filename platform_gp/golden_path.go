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

func GetGoldenPath(url string, name string, branch string, path string, tag string) GoldenPath {

	// if github, set tool to github

	var tool string

	if strings.Contains(url, "github.com") {
		tool = GpGithub
	}

	return GoldenPath{
		Tool:   tool,
		Name:   name,
		URL:    url,
		Branch: branch,
		Path:   path,
		Tag:    tag,
	}
}

func (gp *GoldenPath) CloneGp() error {

	checkoutPath := getCheckoutPath()
	err := DeleteClonePathDir()
	if err != nil {

		log.Error().Msgf("Error cleaning up the folder %s. Error = %v: ", checkoutPath, err)
		return failedCloneGpError()
	}

	gitOptions := &git.CloneOptions{
		URL:          gp.URL,
		Progress:     os.Stdout,
		SingleBranch: false,
	}

	r, err := git.PlainClone(checkoutPath, false, gitOptions)

	if err != nil {
		log.Error().Msgf("Error trying to clone the gp URL %v - Error: %v", gp.URL, err)
	}

	headRef, err := r.Head()
	if err != nil {
		log.Error().Msgf("Unable to return reference of HEAD. Error: %v", err)
		return failedCloneGpError()
	}

	log.Info().Msgf("Cloned the golden path at %s. HEAD ref: %s", checkoutPath, headRef)
	gp.repository = r

	branchRef := getRefForBranchName(r, fmt.Sprintf("refs/remotes/origin/%s", gp.Branch))

	log.Info().Msgf("Found the branch with ref %+v", branchRef)

	worktree, err := r.Worktree()
	if err != nil {
		log.Error().Msgf("Error trying to get worktree for the repository. Error: %v", err)
		return failedCloneGpError()
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef.Name(),
		Create: false,
	})

	if err != nil {
		log.Error().Msgf("Error trying to checkout the branch. Error: %v", err)
		return failedCloneGpError()
	}

	if _, err := os.Stat(path.Join(checkoutPath, gp.Path)); !os.IsNotExist(err) {
		log.Info().Msgf("Succesfully verified that path %v exists in the cloned repo", gp.Path)
	} else {
		log.Error().Msgf("Failed to find the the path %v in the cloned repo. Error: %v", gp.Path, err)
	}

	return nil
}

func DeleteClonePathDir() error {

	checkoutPath := getCheckoutPath()
	return os.RemoveAll(checkoutPath)
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
		log.Error().Msgf("Error going through git refs. Error: %v", err)
	}

	return res
}

func getCheckoutPath() string {
	checkoutPath := os.Getenv("CFS_GP_CHECKOUT_PATH")

	if checkoutPath == "" {
		checkoutPath = "/tmp/gp"
	}

	return checkoutPath
}

func failedCloneGpError() error {
	return errors.New("failed to clone the GoldenPath")
}
