[![Release BadgerCLI](https://github.com/duds-fw/dchangelog/actions/workflows/release.yml/badge.svg)](https://github.com/duds-fw/dchangelog/actions/workflows/release.yml)

# Changelog

Changelog is CLI for generating document for changes code betwen 2 git branch. This CLI used for tracking changes code and TSD.

## Features

- **Generate**: Base on config file, creating a pdf file for changes code.
- **Merge**: Merge all pdf file in 1 folder to 1 pdf file.

## Installation

To install Badger ORM, use the following command:

```bash
go install github.com/duds-fw/dchangelog
```

## Usage

### Generate

```bash
dchangelog generate --config=config.json --dest=parent-branch --src=child-branch
```

### Merge

```bash
dchangelog merge --folder=./tsd
```

# Contributing

Contributions are welcome! Please open an issue or submit a pull request.

# License

This library is under **MIT** License.
