package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
)

func main() {

	dryRunMode := flag.Bool("dryRunMode", true, "Enable or Disable dryrun Mode")
	contractFile := flag.String("contractFile", "platform-order.yaml", "Path to the contract file (in YAML format)")
	flag.Parse()

	proc, err := contract.GetProcessor(*contractFile)

	if err == nil && proc != nil {
		dryRunResult, err := proc.Execute(*dryRunMode)

		if err != nil {
			log.Error().Msgf("Error trying to execute: %v", err)
		}

		if dryRunResult == contract.IdpStatusSuccess {
			log.Info().Msgf("Successfuly executed the idp-cfs contract. Your are now ready to code!")
		}
	}
}
