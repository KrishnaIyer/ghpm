# ghpm

A simple tool in golang to manage Github repositories.

## Usage

1. Export the token as an environment variable.

```bash
$ export GHPM_TOKEN=TOKEN
```

2. Create a YAML config file with the list of repositories. See the example [config.yml](./config.yml).

3. Run the required commands. Ex: List all milestones across all the repositories in the `config.yml` file.

```bash
$ ghpm milestones get
##################################################
		 Milestones
##################################################
Repository: ghpm
1. {"title":"Second Test Milestone","description":"This is another test milestone"}
2. {"title":"Test Milestone","due_on":"08 Sep 23"}
--------------------------------------------------
```

## Options

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
