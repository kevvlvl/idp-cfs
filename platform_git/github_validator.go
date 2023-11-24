package platform_git

import (
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
)

func validateApiResponse(resp *github.Response, e error, msg string) bool {

	if e != nil {
		if resp.Response.StatusCode >= 400 && resp.Response.StatusCode <= 499 {
			log.Warn().Msgf("Client error - "+msg+" - response: %v - error: %v", resp, e)
			return false
		}

		if resp.Response.StatusCode >= 500 && resp.Response.StatusCode <= 599 {
			log.Error().Msgf("SERVER EREROR - "+msg+" - response: %v - error: %v", resp, e)
			return false
		}
	}

	return true
}
