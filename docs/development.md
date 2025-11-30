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
