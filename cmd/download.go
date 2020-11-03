package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"rkd/containers"
	"rkd/git"
	"rkd/helm"
	"rkd/helpers"

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
			Usage: "Download a helm chart list.\n rkd download --helm",
		},
		cli.StringSliceFlag{
			Name:  "image",
			Usage: "Download an image list.\n rkd download --image busybox --image alpine",
		},
		cli.StringFlag{
			Name:  "rancher",
			Usage: "Download Rancher images, rke and helm chart.\n rkd download --rancher v2.5.0",
		},
		cli.StringFlag{
			Name:  "dest",
			Usage: "Set the name of the destination archive , optional a generic name will be set\n rkd download --rancher v2.5.0 --dest myArchive",
		},
	}

	return cli.Command{
		Name:    "download",
		Aliases: []string{"d"},
		Usage:   "Download helm chart list, container images and Rancher release",
		Action:  DownloadDataPack,
		Flags:   DownloadFlags,
	}
}

// DownloadRancherRelease downlaod Rancher chart and images
func DownloadDataPack(c *cli.Context) {

	dest := helpers.GenFileName("datapack")
	if c.String("dest") != "" {
		dest = c.String("dest")
	}
	helpers.CreateDestDir(dest)

	if len(c.StringSlice("image")) > 0 {
		destImg := fmt.Sprintf("%s/images.tar", dest)
		containers.DownloadImage(c.StringSlice("image"), destImg)
	}

	if c.String("rancher") != "" {
		fmt.Printf("Getting Rancher %s\n", c.String("rancher"))
		GetRancherHelmChart(c.String("rancher"), dest)
		GetRancherImages(c.String("rancher"), dest)
	}

	// if c.String("helm") != "" {
	// 	fmt.Printf("Getting Rancher chart %s\n", c.String("helm"))
	// 	GetRancherHelmChart(c.String("helm"))
	// }

	// If no flag provided, download latest chart and images
	if c.String("rancher") == "" && c.String("image") == "" {
		GetRancherHelmChart("latest", dest)
		GetRancherImages("latest", dest)
	}
}

// GetRancherImages downalod rancher images
func GetRancherImages(version string, dest string) {
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
	pathImageList, err := git.GetRancherImageList(release, dest)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// Create imageList
	imgListFile, err := os.Open(pathImageList)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer imgListFile.Close()

	scanner := bufio.NewScanner(imgListFile)
	var imgList []string
	for scanner.Scan() {
		imgList = append(imgList, scanner.Text())
	}

	// Downlaod container images
	destImg := fmt.Sprintf("%s/rancher-images-%s.tar", dest, release.GetTagName())
	containers.DownloadImage(imgList, destImg)
}

// GetRancherHelmChart downalod rancher chart
func GetRancherHelmChart(version string, dest string) {
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
	helm.DownloadChart(rancherRepoName, rancherChartName, version, dest)
}
