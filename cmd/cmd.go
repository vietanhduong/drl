package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vietanhduong/drl/pkg/config"
	"github.com/vietanhduong/drl/pkg/github"
)

func newCommand() *cobra.Command {
	var (
		version     bool
		githubToken string
		silent      bool
		output      string
	)
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s REPO_OWNER/REPO_NAME/[Tag]/FILE_NAME [--github.token PAT] [--output FILE]", os.Args[0]),
		Short: "Download a release asset from GitHub",
		Example: fmt.Sprintf(`# Download the latest release asset from the owner/repo repository
$ %s owner/repo/asset.tar.gz

# Download the latest release asset from the owner/repo repository using a GitHub token
$ GITHUB_TOKEN=ghp_token %s owner/repo/asset.tar.gz

# Download the latest release asset from the owner/repo repository to a file
$ %s owner/repo/asset.tar.gz --output asset.tar.gz

# Download a specific release asset from the owner/repo repository
$ %s owner/repo/v1.0.0/asset.tar.gz

# Decompress tarball file in the fly
$ %s owner/repo/asset.tar.gz | tar -xvz -C /tmp
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0]),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				config.PrintVersion()
				return nil
			}

			// Parse input args
			if len(args) == 0 {
				return fmt.Errorf("missing repository argument")
			}
			parts := strings.Split(args[0], "/")
			if len(parts) < 2 || len(parts) > 4 {
				return fmt.Errorf("invalid repository argument")
			}

			if len(parts) == 3 { // ["owner", "repo", "file"] => ["owner", "repo", "latest", "file"]
				parts = append(parts, parts[len(parts)-1])
				parts[2] = "latest"
			}

			f := os.Stdout
			if output != "" {
				var err error
				if f, err = os.Create(output); err != nil {
					return fmt.Errorf("failed to create file: %w", err)
				}
				defer f.Close()
			}

			gh := github.NewClient(github.WithToken(githubToken))
			return gh.DownloadRelease(parts[0], parts[1], parts[2], parts[3], f)
		},
	}

	cmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "Print the version")
	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "Silent or quiet mode. Do not show progress meter.")
	cmd.Flags().StringVar(&githubToken, "github.token", os.Getenv("GITHUB_TOKEN"), "The GitHub token is used to download releases from a private repository. You can use this by setting the environment variable `GITHUB_TOKEN`.")
	cmd.Flags().StringVarP(&output, "output", "o", "", "The output file to store the release asset. Default is stdout.")
	return cmd
}
