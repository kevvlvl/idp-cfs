package flags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCommandArgs_MissingArgs_DefaultValues_NoErrors(t *testing.T) {

	args := GetCommandArgs()

	assert.NotNil(t, args)
	assert.Equal(t, true, args.DryRun)
	assert.Equal(t, "platform-order.yaml", args.ContractFile)
}
