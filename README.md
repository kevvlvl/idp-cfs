# Code Foundation Shop

First iteration of a minimal idp tool

## Setup

### Define authentication related vars

These environment variables will allow you to authenticate to various git systems for your `code` repository

| Env Var              | Purpose                                                                        |
|:---------------------|:-------------------------------------------------------------------------------|
| CFS_CODE_GITHUB_USER | Github Username using basic auth (in combination with the PAT. see line below) |
| CFS_CODE_GITHUB_PAT  | Github personal access token using basic auth                                  |
| CFS_CODE_GITLAB_USER | Gitlab Username using basic auth (in combination with the PAT. see line below) |
| CFS_CODE_GITLAB_PAT  | Gitlab personal access token using basic auth                                  |

These environment variables will allow you to authenticate to various git systems for your `golden path` repository

| Env Var              | Purpose                                                                        |
|:---------------------|:-------------------------------------------------------------------------------|
| CFS_GP_GITHUB_USER   | Github Username using basic auth (in combination with the PAT. see line below) |
| CFS_GP_GITHUB_PAT    | Github personal access token using basic auth                                  |
| CFS_GP_GITLAB_USER   | Gitlab Username using basic auth (in combination with the PAT. see line below) |
| CFS_GP_GITLAB_PAT    | Gitlab personal access token using basic auth                                  |

### Contract

See _./docs/contract-examples_ for examples of valid platform requests from idp-cfs:

- `new-code` will request a new code and k8s infra
- `update-contract` will request to update the existing code in a new branch and validate that the remainder of the platform exists to ensure no corruption

### Run tests

```shell
go test -cover ./...
```

### Run idp-cfs

```shell
go run main.go --dryRun=true \
               --contractFile=./_docs/contract-examples/platform-order-new-gh.yaml
```

### Gitlab integration

- local instance (container) @ `http://localhost:80` user root/superlab
- Golden Path: `http://localhost/idp-cfs/goldenpath`
- Code Repo: `http://localhost/idp-cfs/code`