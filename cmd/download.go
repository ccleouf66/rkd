package cmd

import (
	"context"
	"fmt"

	"rkd/containers"
	"rkd/git"
	"rkd/helm"

	"github.com/google/go-github/v32/github"
	"github.com/urfave/cli"
)

var (
	rancherRepoURL   = "https://releases.rancher.com/server-charts/stable"
	rancherRepoName  = "rancher-stable"
	rancherChartName = "rancher"
)

// DownloadCommand return downlaod cli command
func DownloadCommand() cli.Command {
	DownloadFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "helm",
			Usage: "Download Rancher helm chart",
		},
		cli.StringFlag{
			Name:  "images",
			Usage: "Download Rancher iamges",
		},
	}

	return cli.Command{
		Name:    "download",
		Aliases: []string{"d"},
		Usage:   "Downlaod Rancher stable release",
		Action:  DownloadRancherRelease,
		Flags:   DownloadFlags,
	}
}

// DownloadRancherRelease downlaod Rancher chart and images
func DownloadRancherRelease(c *cli.Context) {
	if c.String("helm") != "" {
		fmt.Printf("Getting Rancher chart %s\n", c.String("helm"))
		GetRancherHelmChart(c.String("helm"))
	}

	if c.String("images") != "" {
		fmt.Printf("Getting images for Rancher %s\n", c.String("images"))
		GetRancherImages(c.String("images"))
	}

	// If no flag provided, download latest chart and images
	if c.String("helm") == "" && c.String("images") == "" {
		GetRancherHelmChart("latest")
		GetRancherImages("latest")
	}
}

// GetRancherImages downalod rancher images
func GetRancherImages(version string) {
	client := github.NewClient(nil)

	var release *github.RepositoryRelease
	var err error

	if version == "latest" {
		// Get latest release
		release, _, err = client.Repositories.GetLatestRelease(context.Background(), "rancher", "rancher")
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		version = release.GetTagName()
	} else {
		// Get release by tag
		release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "rancher", "rancher", version)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}

	// Download rancher-image.txt of from latest stable
	pathImageList, err := git.GetRancherImageList(release, version)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// Downlaod container images
	dest := fmt.Sprintf("%s/rancher-%s.tar", version, release.GetTagName())
	containers.DownloadImage(pathImageList, dest)
}

// GetRancherHelmChart downalod rancher chart
func GetRancherHelmChart(version string) {
	if version == "latest" {
		client := github.NewClient(nil)
		release, _, err := client.Repositories.GetLatestRelease(context.Background(), "rancher", "rancher")
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		version = release.GetTagName()
	}
	helm.RepoAdd(rancherRepoName, rancherRepoURL)
	helm.RepoUpdate()
	helm.DownloadChart(rancherRepoName, rancherChartName, version, version)
}
