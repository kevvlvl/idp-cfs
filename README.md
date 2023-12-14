# Code Foundation Shop

First iteration of a minimal idp tool

## Setup

### Define the following env vars:

| Env Var           | Purpose                                                                                              |
|:------------------|:-----------------------------------------------------------------------------------------------------|
| CFS_CODE_GIT_USER | Git Username using basic auth (in combination with the PAT. see line below)                          |
| CFS_CODE_GIT_PAT  | personal access token using basic auth                                                               |

### Contract

See _./docs/contract-examples_ for examples of valid platform requests from idp-cfs:

- `new-contract` will request a new code and k8s infra
- `update-contract` will request to update the existing code in a new branch and validate that the remainder of the platform exists to ensure no corruption

### Run tests

```shell
go test -cover ./...
```

### Run idp-cfs

```shell
go run main.go --dryRunMode=true \
               --contractFile=./_docs/contract-examples/platform-order.yaml \
               --gpClonePath="/tmp/idp-cfs-code" \
               --codeClonePath="/tmp/idp-cfs-gp"
```
