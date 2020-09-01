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

// DownloadCommand
func DownloadCommand() cli.Command {
	DownloadFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "helm",
			Usage: "Download Rancher helm chart",
			// Value: "latest",
		},
		cli.StringFlag{
			Name:  "images",
			Usage: "Download Rancher iamges",
			// Value: "2.4.5",
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

// DownloadRancherRelease
func DownloadRancherRelease(c *cli.Context) {
	if c.String("helm") != "" {
		if c.String("helm") == "latest" {
			fmt.Printf("Getting latest Rancher helm chart\n")
			GetRancherHelmChart(c.String("helm"))
		} else {
			fmt.Printf("Getting Rancher helm chart %s\n", c.String("helm"))
		}
		return
	}

	if c.String("images") != "" {
		fmt.Printf("Getting Rancher %s images\n", c.String("images"))
		GetRancherImages(c.String("images"))
		return
	}

	// If no flag provided, download latest chart and images
	GetRancherHelmChart("latest")
	GetRancherImages("latest")
}

func GetRancherImages(version string) {
	fmt.Printf("%s\n", version)
	client := github.NewClient(nil)

	var release *github.RepositoryRelease
	var err error

	if version == "latest" {
		// Get latest release
		release, _, err = client.Repositories.GetLatestRelease(context.Background(), "rancher", "rancher")
		if err != nil {
			fmt.Printf("%s\n", err)
		}
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

func GetRancherHelmChart(version string) {
	helm.RepoAdd(rancherRepoName, rancherRepoURL)
	helm.RepoUpdate()
	helm.DownloadChart(rancherRepoName, rancherChartName, version)
}
