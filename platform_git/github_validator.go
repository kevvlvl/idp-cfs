package platform_git

import (
	"errors"
	"fmt"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
)

func validateApiResponse(resp *github.Response, e error, msg string) error {

	if e != nil {
		if resp.Response.StatusCode == 404 {

			errorMsg := fmt.Sprintf("HTTP404 - "+msg+" - response: %v - error: %v", resp, e)
			log.Warn().Msgf(errorMsg)
			return errors.New(errorMsg)

		} else if resp.Response.StatusCode >= 400 && resp.Response.StatusCode <= 499 {

			errorMsg := fmt.Sprintf("HTTP4xx - "+msg+" - response: %v - error: %v", resp, e)

			log.Warn().Msgf(errorMsg)
			return errors.New(errorMsg)
		} else if resp.Response.StatusCode >= 500 && resp.Response.StatusCode <= 599 {

			errorMsg := fmt.Sprintf("HTTP5xx - "+msg+" - response: %v - error: %v", resp, e)

			log.Error().Msgf(errorMsg)
			return errors.New(errorMsg)
		}
	}

	return nil
}
