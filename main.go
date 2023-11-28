package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
	"idp-cfs/rules"
)

func main() {

	// TODO: receive in Request Body (directly in a POST or through CLI)
	p := rules.GetProcessor(contract.Load("platform-order.yaml"))
	success := p.DryRun()

	if success {
		log.Info().Msgf("Successfuly completed a dry-run without errors. Will execute real actions now.")
	}
}
