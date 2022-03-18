# plutus

A query tool for HashiCorp Vault.

Allows users to query: 
1. Vault Alias and view
    - All Vault Groups they are a part of
    - All Vault Roles they have
    - All Vault Policies they have attached to them (and how they are attached)
    - All Vault Paths they can access along with the capabilities they have for that path (also which policy allows them access to that path)
2. Vault Path and view
    - All the Vault Aliases that have access to the path along with the capabiliites and which policies grant them access.

## Configuration

#### Environment variables

You need to set the following environment variables

```bash
export VAULT_ADDR="a vault addr"
export VAULT_TOKEN="a vault token"
export VAULT_NAMESPACE="a vault namespace"  # You only need to use this if you use the CLI to query vault directly

export GITHUB_ACCESS_TOKEN="a github personal access token"
export REDIS_ADDR="redis address"
export REST_ADDR="REST API address"
```

#### YAML file

You need to add a `config.yaml` file in a config folder. So the file path from the root directory is `config/config.yaml`. A sample config is provided below

```yaml
namespaces:                                                                 # Namespaces you want plutus to cover
  - <example-namespace>
uiAddress: "localhost:4200?baseURL=localhost:8000"                          # UI redirect that can be used to redirect to the proper UI address
githubEnterpise:
  baseURL: "https://github.com/api/v3"                                  # Github Enterprise Reader API base URL
  groupsRepoPath: "/path/to/repo"                                       # Github Enterprise Repo that has the groups information
```

Make sure that the files in github.groupsRepoPath folder are of type `<group-name>-groups.yaml` and look like the following

```yaml
name: group-name
description: ""
spec:
  type: "Security"
  reason: "Access"
  attributes: []
  owners:
  - person-a
  ...
  members:
  - person-b
  ...
  ```
  As of now, only the Entterprise Github Groups Reader is supported but more can be added easily. Look at the group-reader package [README.md](https://github.com/cisco-open/plutus/blob/main/groups-reader/README.md)
## How to run

#### Running locally (Docker)

To run the REST API:
1. Run `docker build . -t plutus:dev` to build the image locally.
2. Run `docker-compose up`

To run the UI:
1. Clone the PlutusUI repo(unpublished)
2. Run `ng serve` in the root directory for the PlutusUI repo

#### Running locally (Binary)
1. Run `go build -o plutus` to generate the executable binary
2. Run `./plutus s` to start the REST API server

## Development

All packages have READMEs in them that can be read to learn more about them. 

## Trivia

Plutus is the Greek god of wealth and so will know exactly where to look in a vault!

## Contributors

[Pranav Bansal](https://github.com/prnvbn)
