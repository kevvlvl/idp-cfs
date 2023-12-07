package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
)

func main() {

	// TODO: Receive parameters through CLI call + env varsh
	proc := contract.GetProcessor("platform-order.yaml")

	if proc != nil {
		dryRunResult, _ := proc.Execute(false)

		if dryRunResult == contract.IdpStatusSuccess {
			log.Info().Msgf("Successfuly executed the idp-cfs contract. Your are now ready to code!")
		}
	}
}
