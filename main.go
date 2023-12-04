package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/rules"
)

func main() {

	// TODO: receive in Request Body (directly in a POST or through CLI)
	proc := rules.GetProcessor("platform-order.yaml")

	if proc != nil {
		dryRunResult, _ := proc.Execute(true)

		if dryRunResult == rules.Success {
			log.Info().Msgf("Successfuly completed a dry-run without errors. Will execute real actions now.")

			exec, _ := proc.Execute(false)

			if exec == rules.Success {
				log.Info().Msgf("Successfuly executed the idp-cfs contract. Your are now ready to code!")
			}
		}
	}
}
