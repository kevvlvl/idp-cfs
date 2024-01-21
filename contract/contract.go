package contract

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"idp-cfs/global"
	"os"
)

// Load unmarshalls the YAML contract file into a struct
func Load(filePath string) (*Contract, error) {

	c := &Contract{}
	buf, err := os.ReadFile(filePath)

	if err != nil {
		log.Error().Msgf("Load() - error trying to read contract file: %v", err)
		return nil, err
	}

	err = yaml.Unmarshal(buf, c)

	if err != nil {
		log.Error().Msgf("Load() - Failed to unmarshal buffer: %v", err)
		return nil, err
	}

	valid := validate(c)

	if !valid {

		errorMsg := "contract metadata is not valid"

		log.Error().Msgf("Load() - %s", errorMsg)
		return nil, errors.New(errorMsg)
	}

	return c, nil
}

// Validate returns true if the contract contains all valid values
func validate(contract *Contract) bool {

	var (
		validAction     = false
		validCode       = false
		validCodeValues = false
		validGpValues   = false
		validDeployment = false
		codeTools       = [3]string{"github", "gitlab", "gitea"}
	)

	if contract != nil {

		if contract.Action == global.NewCode || contract.Action == global.UpdateCode {
			validAction = true
		}

		// Validate Code section
		for _, v := range codeTools {
			if v == contract.Code.Tool {
				validCode = true
			}
		}

		validCodeValues = contract.Code.Repo != "" &&
			contract.Code.Branch != "" &&
			(contract.Code.Url == nil || *contract.Code.Url != "")

		// Default value
		if contract.Code.Workdir == nil || *contract.Code.Workdir == "" {
			contract.Code.Workdir = global.StringPtr("/tmp/idp-cfs-code")
		}

		if contract.Code.Tool == global.ToolGithub && contract.Code.Url == nil {
			contract.Code.Url = global.StringPtr("github.com")
		}

		// Validate Golden-Path section

		if contract.GoldenPath == nil {
			validGpValues = true
		} else {
			validGpValues = contract.GoldenPath.Url != "" &&
				contract.GoldenPath.Path != "" &&
				contract.GoldenPath.Branch != ""

			// Tag field is optional even when the rest of fields are set. Skip validating tag

			// Default value
			if contract.GoldenPath.Workdir == nil || *contract.GoldenPath.Workdir == "" {
				contract.GoldenPath.Workdir = global.StringPtr("/tmp/idp-cfs-gp")
			}
		}

		// Validate Deployment section

		validDeployment = contract.Deployment.Kubernetes.ClusterUrl != "" &&
			contract.Deployment.Kubernetes.Namespace != ""
	}

	log.Info().Msgf("validate() - Valid Contract Action: %v", validAction)
	log.Info().Msgf("validate() - Valid Contract Code Git Tool: %v - Values: %v", validCode, validCodeValues)
	log.Info().Msgf("validate() - Valid Contract Golden-Path: %v", validGpValues)
	log.Info().Msgf("validate() - Valid Contract Deployment: %v", validDeployment)

	return validAction && validCode && validCodeValues && validGpValues && validDeployment
}
