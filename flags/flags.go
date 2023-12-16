package flags

import (
	"flag"
	"idp-cfs/contract"
)

func GetCommandArgs() *CommandArgs {

	dryRunMode := flag.Bool("dryRunMode", true, "Enable or Disable Dryrun Mode")
	contractFile := flag.String("contractFile", "platform-order.yaml", "Path to the contract file (in YAML format)")
	gpClonePath := flag.String("gpClonePath", contract.GoldenPathClonePath, "Path where the golden path is cloned/checked out")
	codeClonePath := flag.String("codeClonePath", contract.CodeClonePath, "Path where we copy the golden path to push into the code repo")

	flag.Parse()

	return &CommandArgs{
		DryRunMode:    *dryRunMode,
		ContractFile:  *contractFile,
		GpClonePath:   *gpClonePath,
		CodeClonePath: *codeClonePath,
	}
}
