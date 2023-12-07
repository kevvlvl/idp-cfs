package contract

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

// Load unmarshalls the YAML contract file into a struct
func Load(filePath string) (*Contract, error) {

	c := &Contract{}
	buf, err := os.ReadFile(filePath)

	if err != nil {

		errorMsg := fmt.Sprintf("error trying to read contract file: %v", err)

		log.Error().Msgf(errorMsg)
		return nil, errors.New(errorMsg)
	}

	err = yaml.Unmarshal(buf, c)

	valid := validate(c)

	if !valid {

		errorMsg := "the contract metadata is not valid"

		log.Error().Msg(errorMsg)
		return nil, errors.New(errorMsg)
	}

	return c, nil
}

// Validate returns true if the contract contains all valid values
func validate(contract *Contract) bool {

	var (
		validCode         = false
		validCodeValues   = false
		validGpValuesOmit = false
		validGpValues     = false
		validDeployment   = false
		codeTools         = [3]string{"github", "gitlab", "gitea"}
	)

	if contract != nil {

		// Validate Code section
		for _, v := range codeTools {
			if v == contract.Code.Tool {
				validCode = true
			}
		}

		validCodeValues = contract.Code.Repo != "" &&
			(contract.Code.Org == nil || *contract.Code.Org != "") &&
			contract.Code.Branch != ""

		// Validate Golden-Path section

		validGpValuesOmit = contract.GoldenPath.Url == nil &&
			contract.GoldenPath.Path == nil &&
			contract.GoldenPath.Branch == nil &&
			contract.GoldenPath.Tag == nil

		if !validGpValuesOmit {
			validGpValues = *contract.GoldenPath.Url != "" &&
				*contract.GoldenPath.Path != "" &&
				*contract.GoldenPath.Branch != ""

			// Tag field is optional even when the rest of fields are set. Skip validating tag
		}

		// Validate Deployment section

		validDeployment = contract.Deployment.Kubernetes.ClusterUrl != "" &&
			contract.Deployment.Kubernetes.Namespace != ""
	}

	log.Info().Msgf("Valid Contract Code Git Tool: %v - Values: %v", validCode, validCodeValues)
	log.Info().Msgf("Valid Contract Golden-Path: %v", validGpValues)
	log.Info().Msgf("Valid Contract Deployment: %v", validDeployment)

	return validCode && validCodeValues && validGpValues && validDeployment
}
