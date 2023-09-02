# ghpm

A simple tool in golang to manage Github repositories.

## Token Authentication

Export the token as an environment variable.

```bash
$ export GHPM_TOKEN=TOKEN
```

## Usage

```
Usage:
  ghpm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  milestones  Manage milestones
  version     Display version information

Flags:
  -c, --config string          config file (Default; config.yml in the current directory) (default "./config.yml")
  -h, --help                   help for ghpm
      --repositories strings   The GitHub repositories
      --token string           The GitHub token
      --username string        The GitHub user or organization name

Use "ghpm [command] --help" for more information about a command.
```

## License

The contents of this repository are packaged under the terms of the [Apache 2.0 License](./LICENSE).
