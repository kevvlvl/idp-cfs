package main

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/request"
)

func main() {

	c := request.Load("platform-order.yaml")

	log.Info().Msgf("Contract loaded: %+v", c)
}
