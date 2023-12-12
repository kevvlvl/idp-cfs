package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
	"idp-cfs/platform_git"
	"idp-cfs/platform_gp"
)

func main() {

	dryRunMode := flag.Bool("dryRunMode", true, "Enable or Disable dryrun Mode")
	contractFile := flag.String("contractFile", "platform-order.yaml", "Path to the contract file (in YAML format)")
	gpCheckoutPath := flag.String("gpCheckoutPath", platform_gp.DefaultCheckoutPath, "Path where the golden path is cloned/checked out")
	codeClonePath := flag.String("codeClonePath", platform_git.CodeClonePath, "Path where we copy the golden path to push into the code repo")

	flag.Parse()

	proc, err := contract.GetProcessor(*contractFile, *gpCheckoutPath, *codeClonePath)

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
