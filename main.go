package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
	"idp-cfs/flags"
)

func main() {

	args := flags.GetCommandArgs()
	proc, err := contract.GetProcessor(args.ContractFile, args.GpClonePath, args.CodeClonePath)

	if err == nil && proc != nil {
		dryRunResult, err := proc.Execute(args.DryRunMode)

		if err != nil {
			log.Error().Msgf("Error trying to execute: %v", err)
		}

		if dryRunResult == contract.IdpStatusSuccess {
			log.Info().Msgf("Successfuly executed the idp-cfs contract. Your are now ready to code!")
		}
	}
}
