package docker

import (
	"bufio"
	"compress/gzip"
	"context"
	"io"
	"os"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

// PullImages Pull images listed in `imgListPath` in local registry
func PullImages(imgListPath string) (err error) {
	dockerCli, err := docker.NewEnvClient()
	if err != nil {
		return err
	}

	imgList, err := os.Open(imgListPath)
	if err != nil {
		return err
	}
	defer imgList.Close()

	scanner := bufio.NewScanner(imgList)
	for scanner.Scan() {
		rc, err := dockerCli.ImagePull(context.Background(), scanner.Text(), types.ImagePullOptions{})
		if err != nil {
			return err
		}

		_, err = io.Copy(os.Stdout, rc)
		if err != nil {
			return err
		}
	}

	return nil
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
	// helpers.CreateDestDir(dest)
	tarfile, err := os.Create(dest)
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

	return dest, nil
}
