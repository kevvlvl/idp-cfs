package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
	"idp-cfs/platform_git"
	"idp-cfs/rules"
)

func main() {

	// TODO: receive in Request Body (directly in a POST or through CLI)
	// action = new-contract = not idempotent. action = update-contract = PUT/idempotent
	c := contract.Load("platform-order.yaml")
	log.Info().Msgf("Contract loaded: %+v", c)

	code := platform_git.GetGithubCode()

	p := rules.GetProcessor(c, code, nil)
	success := p.DryRun()

	if success {
		log.Info().Msgf("Successfuly completed a dry-run without errors. Will execute real actions now.")
	}
}
