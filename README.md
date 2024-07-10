# Download GitHub Release

Download GitHub releases from the command line. This also supports private repositories by providing a GitHub Personal Access Token (PAT) without requiring any additional permissions.

```console
$ drl --help
Download a release asset from GitHub
Notes:
  * The ASSET URL should be in 3 formats:
    - REPO_OWNER/REPO_NAME/[Tag]/FILE_NAME => The [Tag] is optional, if not specified will be treated as "latest"
    - github.com/REPO_OWNER/REPO_NAME/releases/download/[Tag]/FILE_NAME
    - github.com/REPO_OWNER/REPO_NAME/releases/latest/download/FILE_NAME

  * Currently, this tool only supports hosting on GitHub. Enterprise GitHub is not supported yet.

Usage:
  drl REPO_OWNER/REPO_NAME/[Tag]/FILE_NAME [--github.token PAT] [--output FILE] [flags]

Examples:
# Download the latest release asset from the owner/repo repository
$ drl owner/repo/asset.tar.gz

# Download the latest release asset from the owner/repo repository using a GitHub token
$ GITHUB_TOKEN=ghp_token drl owner/repo/asset.tar.gz

# Download the latest release asset from the owner/repo repository to a file
$ drl owner/repo/asset.tar.gz --output asset.tar.gz

# Download a specific release asset from the owner/repo repository
$ drl owner/repo/v1.0.0/asset.tar.gz

# Decompress tarball file in the fly
$ drl owner/repo/asset.tar.gz | tar -xvz -C /tmp


Flags:
  --github.token GITHUB_TOKEN   The GitHub token is used to download releases from a private repository. You can use this by setting the environment variable GITHUB_TOKEN.
  -h, --help                        help for drl
  -o, --output string               The output file to store the release asset. Default is stdout.
  -s, --silent                      Silent or quiet mode. Do not show progress meter.
  -v, --version                     Print the version
```
