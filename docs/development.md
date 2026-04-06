<!-- Source: https://github.com/arduino/tooling-project-assets/blob/main/documentation-templates/contributor-guide/task/development.md -->

# Development Guide

## Prerequisites

The following development tools must be available in your local environment:

- [**Go**](https://go.dev/dl/) - programming language
  - The **Go** version in use is defined in the `go` directive of [`go.mod`](../go.mod).
  - [**gvm**](https://github.com/moovweb/gvm#installing) is recommended if you want to manage multiple installations of **Go** on your system.
- [**Node.js** / **npm**](https://nodejs.org/en/download/) - Node.js dependencies management tool
  - The **Node.js** version in use is defined by the `engines.node` key of [`package.json`](../package.json).
  - [**nvm**](https://github.com/nvm-sh/nvm#installing-and-updating) is recommended if you want to manage multiple installations of **Node.js** on your system.
- [**Python**](https://wiki.python.org/moin/BeginnersGuide/Download)
  - The **Python** version in use is defined by the `tool.poetry.dependencies.python` key of [`pyproject.toml`](../pyproject.toml).

## Building the Project

You can build the project by running this command from the root folder of the project:

```text
go tool task build
```

## Running Checks

Checks and tests are set up to ensure the project content is functional and compliant with the established standards.

You can run the checks by running this command from the root folder of the project:

```text
go tool task check
```

### Dependency License Metadata

Metadata about the license types of all dependencies is cached in the repository. To update this cache, run the following command from the repository root folder:

```text
go tool task general:cache-dep-licenses
```

The necessary **Licensed** tool can be installed by following [these instructions](https://github.com/licensee/licensed#installation).

An updated cache is also generated whenever the cache is found to be outdated by the "**Check Go Dependencies**" CI workflow and made available for download via the `dep-licenses-cache` [workflow artifact](https://docs.github.com/actions/how-tos/manage-workflow-runs/download-workflow-artifacts).

## Automatic Corrections

Tools are provided to automatically bring the project into compliance with some of the required checks.

You can make these automatic fixes by running this command from the root folder of the project:

```text
go tool task fix
```
