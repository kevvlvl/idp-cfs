package main

import (
	"github.com/rs/zerolog/log"
	platform_git "idp-cfs/platform_git"
	"idp-cfs/request"
)

func main() {

	c := request.Load("platform-order.yaml")
	log.Info().Msgf("Contract loaded: %+v", c)

	ghc := platform_git.Login()

	ghc.GetOrganization("kevvlvl")
	ghc.GetRepositories("idp-cfs")
}
