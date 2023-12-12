# Code Foundation Shop

First iteration of a minimal idp tool

## Setup

### Define the following env vars:

| Env Var                | Purpose                                                                                     |
|:-----------------------|:--------------------------------------------------------------------------------------------|
| CFS_GITHUB_USER        | Github Username for basic auth (in combination with the PAT. see line below)                |
| CFS_GITHUB_PAT         | personal access token authentication to Github if using Github as code target               |
| CFS_GP_CHECKOUT_PATH   | local directory to temporarily checkout the golden path. default = /tmp/gp                  |
| CFS_GP_CODE_CLONE_PATH | local directory to temporarily clone the newly created code repo and to push the gp into it |

### Contract

See _./docs/contract-examples_ for examples of valid platform requests from idp-cfs:

- `new-contract` will request a new code and k8s infra
- `update-contract` will request to update the existing code in a new branch and validate that the remainder of the platform exists to ensure no corruption

### Run tests

```shell
go test ./...
```

### Run idp-cfs

```shell
go run main.go --dryRunMode=true --contractFile=./platform-order.yaml
```