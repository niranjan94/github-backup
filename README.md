# Github Repository Backup

This Go utility let's you schedule backups of your GitHub repositories. It uses the GitHub API and saves the archive of each repository to your specified directory. It leverages goroutines for concurrent backups.

## Setup

You need to have [Go installed](https://golang.org/doc/install) and an environment compatible with Go modules.

## Usage

Firstly, clone the repository:

```
git clone https://github.com/niranjan94/github-backup.git
```

Next, enter the directory:

```
cd github-backup
```

Build the binary:

```
go build -v -o github_backup cmd/github-backup/main.go
```

Run the binary with environment variables:

```
GITHUB_BACKUP_OWNERS="<github_owners>" GITHUB_BACKUP_TOKEN="<your_token>" GITHUB_BACKUP_DIRECTORY="<backup_directory>" ./github-backup
```

## Environment Variables

You need to pass in a few environment variables:

- `GITHUB_BACKUP_OWNERS`: A comma-separated list of GitHub usernames whose repositories you want to backup. Mandatory.
- `GITHUB_BACKUP_TOKEN`: Your GitHub Personal Access Token (PAT). To generate a PAT, follow [this guide](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token). The token should have the `repo` scope. Mandatory.
- `GITHUB_BACKUP_DIRECTORY`: The directory where the backups will be stored. The directory will be created if it doesn't exist. Mandatory.
- `GITHUB_BACKUP_CONCURRENCY`: The number of concurrent backups running. This is limited by your system's maximum number of goroutines. Defaults to 10.
- `GITHUB_BACKUP_PERIOD_SECONDS`: The backup frequency, in seconds. Defaults to `86400`, which is 24 hours.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE) file for details.

## Contributing

Your contributions are always welcome! Please raise an issue or open a pull request.
