package flags

type CommandArgs struct {
	DryRun       bool
	ContractFile string
}

const (
	ArgDryRun       = "dryRun"
	ArgContractFile = "contractFile"
)
