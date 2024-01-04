package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs2/contract"
	"idp-cfs2/flags"
)

func main() {

	log.Info().Msg("START IDP-CFS2")

	args := flags.GetCommandArgs()
	state := contract.GetState(args.DryRun, args.ContractFile)

	if state != nil {
		_, err := state.Deploy()
		if err != nil {
			log.Error().Msgf("Error when trying to deploy: %v", err)
			return
		}
	}

	log.Info().Msg("COMPLETED IDP-CFS2")
}
