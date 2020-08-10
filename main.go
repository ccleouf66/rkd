package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
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

	// Download rancher-image.txt
	path, err := GetRancherImageList(latestRelease, "tmp")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Println(path)

	releases, err := GetRepoStablRelease("rancher", "rancher")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Num. Name - TagName\n")
	for index, release := range releases {
		fmt.Printf("%d. %s - %s\n", index, release.GetName(), release.GetTagName())
	}

	// Pull images listed in rancher-images.txt
	// path, err = PullImages(path)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// }

	// Create tar.gz with images listed in rancher-images.txt
	// path, err = SaveImageFromFile(path)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// }

}

// CreateDestDir check and create directory if not exist
func CreateDestDir(name string) error {
	destInfo, err := os.Stat(name)
	if os.IsNotExist(err) {
		errCreateDir := os.MkdirAll(name, 0755)
		if errCreateDir != nil {
			return err
		}
		return nil
	} else if !destInfo.IsDir() {
		return errors.New("Destination is file but should be a folder")
	} else {
		return err
	}
}

// GetRepoStablRelease list stable release for given repo
func GetRepoStablRelease(user string, repo string) ([]*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), "rancher", "rancher", nil)
	if err != nil {
		return nil, err
	}

	var stableReleases []*github.RepositoryRelease
	for _, release := range releases {
		if !release.GetPrerelease() {
			stableReleases = append(stableReleases, release)
		}
	}

	return stableReleases, nil
}

// GetRancherImageList download rancher-images.txt file from rancher release
func GetRancherImageList(release *github.RepositoryRelease, dest string) (path string, err error) {
	client := github.NewClient(nil)

	releaseAssets, _, err := client.Repositories.ListReleaseAssets(context.Background(), "rancher", "rancher", release.GetID(), nil)
	if err != nil {
		return "", err
	}

	for _, asset := range releaseAssets {
		if asset.GetName() == "rancher-images.txt" {

			// Create destination file
			CreateDestDir(dest)
			pathFile := fmt.Sprintf("%s/%s", dest, asset.GetName())
			out, err := os.Create(pathFile)
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
				return pathFile, nil
			}

			_, err = io.Copy(out, rc)
			if err != nil {
				return "", err
			}
			return pathFile, nil
		}
	}

	return "", nil
}

// PullImages Pull images listed in `imgListPath` in local registry
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

// SaveImageFromFile save image from local registry listed in `imgListPath` to tarball
func SaveImageFromFile(imgListPath string, dest string) (tarPath string, err error) {

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

	// Create destination dir and file
	CreateDestDir(dest)
	pathFile := fmt.Sprintf("%s/rancher-images.tar.gz", dest)
	tarfile, err := os.Create(pathFile)
	if err != nil {
		return "", err
	}
	defer tarfile.Close()

	// gzip writer
	gw := gzip.NewWriter(tarfile)
	defer gw.Close()

	// Save images
	rc, err := dockerCli.ImageSave(context.Background(), imgs)
	if err != nil {
		return "", err
	}

	// Write it in dest file
	if _, err := io.Copy(gw, rc); err != nil {
		return "", err
	}

	return pathFile, nil
}
