package global

import (
	"errors"
	"github.com/rs/zerolog/log"
)

func LogError(msg string) error {
	log.Error().Msg(msg)
	return errors.New(msg)
}
