package util

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ValidateApiResponse(resp *http.Response, e error, msg string) error {

	if e != nil {
		if resp.StatusCode == 404 {

			errorMsg := fmt.Sprintf("HTTP404 - "+msg+" - response: %v - error: %v", resp, e)
			log.Warn().Msgf(errorMsg)
			return errors.New(errorMsg)

		} else if resp.StatusCode >= 400 && resp.StatusCode <= 499 {

			errorMsg := fmt.Sprintf("HTTP4xx - "+msg+" - response: %v - error: %v", resp, e)

			log.Warn().Msgf(errorMsg)
			return errors.New(errorMsg)
		} else if resp.StatusCode >= 500 && resp.StatusCode <= 599 {

			errorMsg := fmt.Sprintf("HTTP5xx - "+msg+" - response: %v - error: %v", resp, e)

			log.Error().Msgf(errorMsg)
			return errors.New(errorMsg)
		}
	}

	return nil
}
