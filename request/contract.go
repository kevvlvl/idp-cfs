package request

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

// Load unmarshalls the YAML contract file into a struct
func Load(filePath string) *Contract {

	order := &Contract{}
	buf, err := os.ReadFile(filePath)

	if err != nil {
		log.Error().Msgf("Error trying to read request file: %v", err)
		return nil
	}

	err = yaml.Unmarshal(buf, order)
	return order
}

// Validate returns true if the contract contains all valid values
func Validate(contract *Contract) bool {

	return true
}
