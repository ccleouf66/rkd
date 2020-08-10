package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/github"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

func main() {
	client := github.NewClient(nil)

	latestRelease, _, err := client.Repositories.GetLatestRelease(context.Background(), "rancher", "rancher")
	if err != nil {
		fmt.Printf("%s", err)
	}

	path, err := GetRancherImageList(latestRelease)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	fmt.Println(path)

	// path, err = PullImages(path)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// }

	path, err = SaveImageFromFile(path)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	// releases, _, err := client.Repositories.ListReleases(context.Background(), "rancher", "rancher", nil)
	// if err != nil {
	// 	fmt.Printf("%s", err)
	// }

	// fmt.Printf("%s\n", releases[0].GetTagName())
	// fmt.Printf("%s - %s\n", latestRelease.GetName(), latestRelease.GetTagName())
}

// GetRancherImageList download rancher-images.txt file from rancher release
func GetRancherImageList(release *github.RepositoryRelease) (path string, err error) {
	client := github.NewClient(nil)

	releaseAssets, _, err := client.Repositories.ListReleaseAssets(context.Background(), "rancher", "rancher", release.GetID(), nil)
	if err != nil {
		return "", err
	}

	for _, asset := range releaseAssets {
		if asset.GetName() == "rancher-images.txt" {

			// Create destination file
			out, err := os.Create(asset.GetName())
			if err != nil {
				return "", err
			}
			defer out.Close()

			// Try to download file
			rc, redirect, err := client.Repositories.DownloadReleaseAsset(context.Background(), "rancher", "rancher", asset.GetID(), nil)
			if err != nil {
				return "", err
			}
			// If redirect, follow and download file at new source
			if rc == nil {
				resp, err := http.Get(redirect)
				if err != nil {
					return "", err
				}

				_, err = io.Copy(out, resp.Body)
				if err != nil {
					return "", err
				}
				return asset.GetName(), nil
			}

			_, err = io.Copy(out, rc)
			if err != nil {
				return "", err
			}
			return asset.GetName(), nil
		}
	}

	return "", nil
}

// PullImages get an image list and pull images to tar.gz file
func PullImages(imgListPath string) (archivePath string, err error) {
	dockerCli, err := docker.NewEnvClient()
	if err != nil {
		return "", err
	}

	imgList, err := os.Open(imgListPath)
	if err != nil {
		return "", err
	}
	defer imgList.Close()

	scanner := bufio.NewScanner(imgList)
	for scanner.Scan() {
		rc, err := dockerCli.ImagePull(context.Background(), scanner.Text(), types.ImagePullOptions{})
		if err != nil {
			return "", err
		}

		_, err = io.Copy(os.Stdout, rc)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

// SaveImageFromFile save image from local registry listed in imgListPath to tarball
func SaveImageFromFile(imgListPath string) (tarPath string, err error) {
	dockerCli, err := docker.NewEnvClient()
	if err != nil {
		return "", err
	}

	imgList, err := os.Open(imgListPath)
	if err != nil {
		return "", err
	}
	defer imgList.Close()

	var imgs []string

	scanner := bufio.NewScanner(imgList)
	for scanner.Scan() {
		imgs = append(imgs, scanner.Text())
	}

	tarfile, err := os.Create("rancher-images.tar.gz")
	if err != nil {
		return "", err
	}
	defer tarfile.Close()

	gw := gzip.NewWriter(tarfile)
	defer gw.Close()

	rc, err := dockerCli.ImageSave(context.Background(), imgs)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(gw, rc); err != nil {
		return "", err
	}

	return "rancher-images.tar.gz", nil
}
