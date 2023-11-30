package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/rules"
)

func main() {

	// TODO: receive in Request Body (directly in a POST or through CLI)
	proc := rules.GetProcessor("platform-order.yaml")

	if proc != nil {
		success, _ := proc.DryRun()

		if success {
			log.Info().Msgf("Successfuly completed a dry-run without errors. Will execute real actions now.")
		}
	}
}
