package flags

import (
	"flag"
)

func GetCommandArgs() *CommandArgs {

	var (
		dryRunMode   = true
		contractFile = ""
	)

	flag.BoolVar(&dryRunMode, ArgDryRun, true, "Enable or Disable Dry-run")
	flag.StringVar(&contractFile, ArgContractFile, "platform-order.yaml", "Path to the contract file (in YAML format)")
	flag.Parse()

	return &CommandArgs{
		DryRun:       dryRunMode,
		ContractFile: contractFile,
	}
}
