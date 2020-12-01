package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
		cli.StringSliceFlag{
			Name:  "helm",
			Usage: "Download a helm chart list.\nrkd download --helm https://REPO_URL/REPO_NAME/CHART_NAME \nrkd download --helm https://charts.bitnami.com/bitnami/postgresql-ha",
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
			Usage: "Set the name of the destination folder containing all downloaded componets (images, rancher images, rancher chart, ...), optional, a generic name will be set: datapack-YYY_MM_JJThh_mm_ss\n rkd download --rancher v2.5.0 --dest destFolder",
		},
		cli.StringFlag{
			Name:  "imgarchname",
			Usage: "Set the name of the generated archive containing additionals images provided by --image flag, optional, default archive name is images.tar\n rkd download --rancher v2.5.0 --image busybox --imgarchname busybox.tar",
		},
		cli.StringFlag{
			Name:  "helmuser",
			Usage: "Set the username for private helm repo with --helm flag, optional\nrkd download --helm https://charts.bitnami.com/bitnami/postgresql-ha --helmuser user --helmpwd pwd",
		},
		cli.StringFlag{
			Name:  "helmpwd",
			Usage: "Set the password for private helm repo with --helm flag, optional\nrkd download --helm https://charts.bitnami.com/bitnami/postgresql-ha --helmuser user --helmpwd pwd",
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

// DownloadDataPack downlaod Rancher chart and images
func DownloadDataPack(c *cli.Context) {

	var dest string

	if c.String("dest") != "" {
		dest = c.String("dest")
	} else {
		dest = helpers.GenFileName("datapack")
	}
	helpers.CreateDestDir(dest)

	// Images
	if len(c.StringSlice("image")) > 0 {
		var destImg string

		if c.String("imgarchname") != "" {
			destImg = fmt.Sprintf("%s/%s", dest, c.String("imgarchname"))
		} else {
			destImg = fmt.Sprintf("%s/images.tar", dest)
		}
		containers.DownloadImage(c.StringSlice("image"), destImg)
	}

	// Helm
	if len(c.StringSlice("helm")) > 0 {
		for _, el := range c.StringSlice("helm") {
			s := strings.Split(el, "/")
			chartName := s[len(s)-1]
			repoName := s[len(s)-2]
			repoURL := strings.Join(s[:len(s)-1], "/")
			chartDest := fmt.Sprintf("%s/%s", dest, chartName)

			fmt.Printf("Getting %s chart\n", chartName)

			err := helm.RepoAdd(repoName, repoURL)
			if err != nil {
				log.Println(err)
			}
			helm.RepoUpdate()
			chartPath := helm.DownloadChart(repoName, chartName, "", chartDest)

			// Get chart image list
			imgList, err := helm.GetChartImages(chartPath)
			if err != nil {
				log.Printf("Err getting image list of chart %s.\n%s\n", chartName, err)
				continue
			}

			// Downlaod container images
			destImg := fmt.Sprintf("%s/%s-images.tar", chartDest, chartName)
			err = containers.DownloadImage(imgList, destImg)
			if err != nil {
				log.Printf("%s\n", err)
			}
		}
	}

	// Rancher
	if c.String("rancher") != "" {
		fmt.Printf("Getting Rancher %s\n", c.String("rancher"))
		GetRancherHelmChart(c.String("rancher"), dest)
		GetRancherImages(c.String("rancher"), dest)
	}

	// If no flag provided, download latest chart and images
	if c.String("rancher") == "" && len(c.StringSlice("image")) == 0 && len(c.StringSlice("helm")) == 0 {
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
	err := helm.RepoAdd(rancherRepoName, rancherRepoURL)
	if err != nil {
		log.Println(err)
	}
	helm.RepoUpdate()
	helm.DownloadChart(rancherRepoName, rancherChartName, version, dest)
}
