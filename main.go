package main

import (
	"context"
	"fmt"
	"os"

	"rkd/cmd"
	"rkd/docker"
	"rkd/git"
	"rkd/helm"

	"github.com/google/go-github/v32/github"
	"github.com/urfave/cli"
)

var (
	url         = "https://kubernetes-charts.storage.googleapis.com"
	repoName    = "stable"
	chartName   = "mysql"
	releaseName = "mysql-dev"
)

func main() {

	// Add helm repo
	helm.RepoAdd(repoName, url)

	client := github.NewClient(nil)

	latestRelease, _, err := client.Repositories.GetLatestRelease(context.Background(), "rancher", "rancher")
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Download rancher-image.txt
	pathImageList, err := git.GetRancherImageList(latestRelease, "tmp")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Println(pathImageList)

	// Pull images listed in rancher-images.txt
	err = docker.PullImages(pathImageList)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	// Create tar.gz with images listed in rancher-images.txt
	path, err := docker.SaveImageFromFile(pathImageList, "tmp/rancher-images.tar.gz")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s creates", path)

	app := cli.NewApp()
	app.Name = "rkd"
	app.Usage = "Rancher Kubernetes Downloader"

	app.Commands = []cli.Command{
		cmd.ListCommand(),
	}

	app.Run(os.Args)

}
