package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
)

func main() {

	/*-----------------------------------
	TODO parameters of CLI:
	- contract file (string path)
	- dryRunMode (bool)
	-----------------------------------*/
	dryRunMode := false
	contractFile := "platform-order.yaml"

	proc, err := contract.GetProcessor(contractFile)

	if err == nil && proc != nil {
		dryRunResult, err := proc.Execute(dryRunMode)

		if err != nil {
			log.Error().Msgf("Error trying to execute: %v", err)
		}

		if dryRunResult == contract.IdpStatusSuccess {
			log.Info().Msgf("Successfuly executed the idp-cfs contract. Your are now ready to code!")
		}
	}
}
