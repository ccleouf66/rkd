package main

import (
	"context"
	"fmt"
	"os"
	"rkd/cmd"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

var (
	rancherStableRepoURL   = "https://releases.rancher.com/server-charts/stable"
	rancherLatestRepoURL = "https://releases.rancher.com/server-charts/latest"
	rancherStableRepoName  = "rancher-stable"
	rancherLatestRepoName  = "rancher-latest"
	rancherChartName = "rancher"
)

func main() {

	///test

	// get latest from github
	client := github.NewClient(nil)
	release, _, err = client.Repositories.GetLatestRelease(context.Background(), "rancher", "rancher")
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	// check if version exist in stable repo
	repo.FindChartInRepoURL(rancherStableRepoURL, rancherChartName, release.GetTagName(), certFile, keyFile, caFile string, getters getter.Providers) (string, error)


	///

	app := cli.NewApp()
	app.Name = "rkd"
	app.Usage = "Rancher Kubernetes Downloader"

	app.Commands = []cli.Command{
		cmd.ListCommand(),
		cmd.DownloadCommand(),
	}

	app.Run(os.Args)
}
