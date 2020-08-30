package main

import (
	"context"
	"fmt"

	"rkd/containers"
	"rkd/git"

	"github.com/google/go-github/v32/github"
)

var (
	url         = "https://kubernetes-charts.storage.googleapis.com"
	repoName    = "stable"
	chartName   = "mysql"
	releaseName = "mysql-dev"
)

func main() {

	// //Helm
	// Add helm repo
	// helm.RepoAdd(repoName, url)
	// helm.RepoAdd("rancher-stable", "https://releases.rancher.com/server-charts/stable")
	// helm.RepoUpdate()
	// helm.DownloadChart("rancher-stable", "rancher", "tmp")

	// Containers images
	client := github.NewClient(nil)

	// Get Rancher latest stable
	latestRelease, _, err := client.Repositories.GetLatestRelease(context.Background(), "rancher", "rancher")
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Download rancher-image.txt of from latest stable
	pathImageList, err := git.GetRancherImageList(latestRelease, "tmp")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	containers.DownloadImage(pathImageList, "tmp/img.tar")

	// app := cli.NewApp()
	// app.Name = "rkd"
	// app.Usage = "Rancher Kubernetes Downloader"

	// app.Commands = []cli.Command{
	// 	cmd.ListCommand(),
	// }

	// app.Run(os.Args)
}
