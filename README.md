# ino-platform-discovery

[![Check Exclusions File status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-exclusions.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-exclusions.yml)
[![Check Go status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-go-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-go-task.yml)
[![Check Go Dependencies status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-go-dependencies-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-go-dependencies-task.yml)
[![Check Markdown status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-markdown-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-markdown-task.yml)
[![Check npm status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-npm-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-npm-task.yml)
[![Check Poetry status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-poetry-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-poetry-task.yml)
[![Check Prettier Formatting status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-prettier-formatting-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-prettier-formatting-task.yml)
[![Check ToC status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-toc-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-toc-task.yml)
[![Check YAML status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-yaml-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/check-yaml-task.yml)
[![Spell Check status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/spell-check-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/spell-check-task.yml)
[![Sync Labels status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/sync-labels-npm.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/sync-labels-npm.yml)
[![Test Go status](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/test-go-task.yml/badge.svg)](https://github.com/per1234-org/ino-platform-discovery/actions/workflows/test-go-task.yml)
[![Codecov](https://codecov.io/gh/per1234-org/ino-platform-discovery/branch/main/graph/badge.svg)](https://codecov.io/gh/per1234-org/ino-platform-discovery)

A tool to discover Arduino [package indexes](https://arduino.github.io/arduino-cli/latest/package_index_json-specification/) and [boards platforms](https://arduino.github.io/arduino-cli/latest/platform-specification/).

## Table of Contents

<!-- toc -->

- [Usage](#usage)
  - [Prerequisites](#prerequisites)
  - [A. Obtain inoplatforms Catalog File](#a-obtain-inoplatforms-catalog-file)
  - [B. Set Up GitHub Access Token](#b-set-up-github-access-token)
  - [C. Run the Tool](#c-run-the-tool)
  - [D. Review Discoveries](#d-review-discoveries)
  - [Notes](#notes)
- [Configuration](#configuration)
  - [Exclusions](#exclusions)
    - [Keys](#keys)
- [Limitations](#limitations)
- [Contributing](#contributing)

<!-- tocstop -->

## Usage

### Prerequisites

The following development tools must be available in your local environment:

- [**Go**](https://go.dev/dl/) - programming language
  - The **Go** version in use is defined in the `go` directive of [`go.mod`](go.mod).
  - [**gvm**](https://github.com/moovweb/gvm#installing) is recommended if you want to manage multiple installations of **Go** on your system.

### A. Obtain inoplatforms Catalog File

[**inoplatforms**](https://github.com/per1234/inoplatforms) is a catalog of all known Arduino boards platforms. **ino-platform-discovery** is intended to be used to discover previously unknown platforms. For this reason, this tool compares the candidate platforms it finds against the content of the **inoplatforms** catalog, and excludes those platforms that are already cataloged.

For this reason, it is necessary to download a copy of the **inoplatforms** catalog file for **ino-platform-discovery** to access.

1. Click the following link to open the GitHub page for the catalog file in your web browser:<br />
   https://github.com/per1234/inoplatforms/blob/main/ino-hardware-package-list.tsv
1. Click the button that looks like an arrow pointing downwards into a tray ("Download raw file"), which you will see on the toolbar on that page.

The file will be downloaded to your hard drive.

### B. Set Up GitHub Access Token

The tool makes requests to the GitHub API. These requests must be authenticated with a [GitHub access token](https://docs.github.com/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens), for the following reasons:

- GitHub [requires authentication of all requests to the `/search/code` endpoint](https://docs.github.com/rest/search/search#search-code:~:text=This%20endpoint%20requires%20you%20to%20authenticate), even though we are only interested in public code.
- To avoid rate limiting while making requests to the `/repos/{owner}/{repo}` endpoint.

1. Click the following link to open the token creation page in your web browser:<br />
   https://github.com/settings/personal-access-tokens/new
1. Type a meaningful name into the "**Token name**" field.
1. Select an appropriate expiration from the "**Expiration**" menu.
1. Select the **Repository access > Public repositories** radio button.
1. Leave the "**Permissions**" section empty.
1. Click the "**Generate token**" button.<br />
   The "**New personal access token**" dialog will open.
1. Click the "**Generate token**" button.<br />
   The token will be generated, and its value displayed.
1. Save the displayed token value to a safe place.

### C. Run the Tool

1. Open a terminal in the project folder.
1. Type the following command in the terminal:
   ```text
   GITHUB_TOKEN="<token>" go run main.go --catalog "<catalog path>"
   ```
1. Replace the `<token>` placeholder with the value of the GitHub access token you created for use by the script.
1. Replace the `<catalog path>` placeholder with the path of the [**inoplatforms** catalog file](#a-obtain-inoplatforms-catalog-file) on your hard drive.
1. Press the <kbd>**Enter**</kbd> key.

The tool run will take some time to complete.

### D. Review Discoveries

A successful run of **ino-platform-discovery** produces a spreadsheet of discoveries. The tool attempts to filter out repositories that do not represent novel Arduino boards platforms. However, the spreadsheet is still likely to contain items that are not of value to the user. For this reason, the discoveries must be manually reviewed.

A discovery might fall into one of the following classifications:

- **Original:** An independent creation.
- **Hard Fork:** A derivative of an existing platform containing significant modifications.
- **Staging Fork:** A repository with the sole purpose of staging work for contribution to the parent platform.
  - **ⓘ** In the case where a proposal of significant modifications is not accepted by the maintainer of the parent, the creator of what was originally intended to be a staging fork may decide to maintain it as a **hard fork**.
- **Trivial Fork:** A fork that contains modifications, but the modifications are insignificant. These may be created in the case where the owner performed some experimentation, but did not produce something relevant to platform users.
- **Duplicate:** A copy of a platform. The copy may have been made at any point in the development history of the parent project, so these can have different content from the latest revision of the parent, but only in the absence of recent changes in the parent.
  - **ⓘ** We would expect a copy to be marked as a fork by GitHub (in which case it would have been filtered out by **ino-platform-discovery**). However, copies may have be created in a manner that does not produce that linkage.

Add any discoveries that are determined to be invalid, trivial forks, or duplicates to the [exclusion](#exclusions) so you can avoid the need to review them again for future runs.

### Notes

The content of the "**Boards Manager URL**" column in the discoveries spreadsheet is the URL of the "raw" package index source file. This will generally work as the URL to use in the Arduino IDE "**Additional Boards Manager URL**" preference. However, the platform maintainer may specify a different URL in the installation instructions. In this case, that canonical URL should generally be given preference over the URL generated by this tool.

## Configuration

### Exclusions

**ino-platform-discovery** can be configured to exclude items that would otherwise be included in the discoveries. This is done via a data file in [YAML](https://www.yaml.info/learn/index.html) format. Pass the path to the file as an argument to the `--exclusions` flag in the `ino-platform-discovery` invocation.

An exclusions file is maintained in this project: [**here**](./exclusions.yml).

#### Keys

The exclusions file is a [sequence](https://yaml.org/spec/1.2.2/ext/glossary/#sequence) (i.e., array) of [mappings](https://yaml.org/spec/1.2.2/ext/glossary/#mapping) (i.e., objects/dictionaries). Each mapping may contain the following keys:

- **`host`:** (required) The Git host.
- **`owner`:** (required) The name of the repository owner.
- **`name`:** (optional) The name of the repository.
  - Default: `.*` (exclude all repositories from the given owner)
- **`path`:** (optional) The path of the discovery. This is the path of the package index file or the folder containing a platform.
  - Default: `.*` (exclude all discoveries from the given repository)

The values of the keys are regular expressions. The regular expression syntax is that of the Go [`regexp` package](https://pkg.go.dev/regexp).

## Limitations

Discoveries are made via data provided by the [GitHub code search API](https://docs.github.com/rest/search/search#search-code), and so is subject to the [limitations imposes on that endpoint](https://docs.github.com/search-github/searching-on-github/searching-code#considerations-for-code-search). This means the discoveries are limited to projects that meet the following conditions:

- Hosted on GitHub.
- Present in the [default branch](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-branches#about-the-default-branch) of the repository.
- Repository is not [archived](https://docs.github.com/repositories/archiving-a-github-repository/archiving-repositories).
- Repository has had activity, or been returned in search results, within the last year.
- (For package index discovery) Index file is smaller than 384 kB.
- (For platform discovery) [`boards.txt`](https://arduino.github.io/arduino-cli/latest/platform-specification/#boardstxt) file is smaller than 384 kB.
- Repository contains less than 500000 files.

## Contributing

See [the **Contributor Guide**](docs/CONTRIBUTING.md).
