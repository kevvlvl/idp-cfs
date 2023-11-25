package main

import (
	"github.com/rs/zerolog/log"
	platform_git "idp-cfs/platform_git"
	"idp-cfs/request"
)

func main() {

	c := request.Load("platform-order.yaml")
	log.Info().Msgf("Contract loaded: %+v", c)

	// Initiate validator:
	// valid git tool?
	// Connect to git
	// Can I create a repo?
	// valid golden path git path (using a generic git client)
	// does it contain code?
	// push golden path code into newly created repo at desired branch

	code := platform_git.GetGithubCode()
	code.Organization, code.OrgExists = code.GetOrganization("kevvlvl")
	code.Repository, code.RepoExists = code.GetRepository("idp-cfs")

}
