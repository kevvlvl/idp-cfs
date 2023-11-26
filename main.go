package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
	platform_git "idp-cfs/platform_git"
	"idp-cfs/rules"
)

func main() {

	// TODO: receive in Request Body (directly or through CLI)
	contractFile := "platform-order.yaml"

	c := contract.Load(contractFile)
	log.Info().Msgf("Contract loaded: %+v", c)

	code := platform_git.GetGithubCode()

	p := rules.GetProcessor(c, code, nil)
	success := p.DryRun()

	if success {
		log.Info().Msgf("Successfuly completed a dry-run without errors. Will execute real actions now.")
	}

	//code.Organization, code.OrgExists = code.GetOrganization("kevvlvl")
	//code.Repository, code.RepoExists = code.GetRepository("idp-cfs")
}
