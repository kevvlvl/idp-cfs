package contract

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

// Load unmarshalls the YAML contract file into a struct
func Load(fr FileReader, filePath string) (*Contract, error) {

	c := &Contract{}
	buf, err := fr.ReadFile(filePath)

	if err != nil {
		log.Error().Msgf("error trying to read contract file: %v", err)
		return nil, err
	}

	err = yaml.Unmarshal(buf, c)

	if err != nil {
		log.Error().Msgf("Failed to unmarshal buffer: %v", err)
		return nil, err
	}

	valid := validate(c)

	if !valid {

		errorMsg := "contract metadata is not valid"

		log.Error().Msg(errorMsg)
		return nil, errors.New(errorMsg)
	}

	return c, nil
}

// Validate returns true if the contract contains all valid values
func validate(contract *Contract) bool {

	var (
		validAction       = false
		validCode         = false
		validCodeValues   = false
		validGpValuesOmit = false
		validGpValues     = false
		validDeployment   = false
		codeTools         = [3]string{"github", "gitlab", "gitea"}
	)

	if contract != nil {

		if contract.Action == NewContract || contract.Action == UpdateContract {
			validAction = true
		}

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

	log.Info().Msgf("Valid Contract Action: %v", validAction)
	log.Info().Msgf("Valid Contract Code Git Tool: %v - Values: %v", validCode, validCodeValues)
	log.Info().Msgf("Valid Contract Golden-Path: %v", validGpValues)
	log.Info().Msgf("Valid Contract Deployment: %v", validDeployment)

	return validAction && validCode && validCodeValues && validGpValues && validDeployment
}

func (f *CfsFileReader) ReadFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}
