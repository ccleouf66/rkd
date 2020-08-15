package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ccleouf66/rkd/cmd"

	"gopkg.in/yaml.v2"

	"github.com/gofrs/flock"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"

	helm_cli "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

var settings *helm_cli.EnvSettings

var (
	url         = "https://kubernetes-charts.storage.googleapis.com"
	repoName    = "stable"
	chartName   = "mysql"
	releaseName = "mysql-dev"
	namespace   = "mysql-test"
	args        = map[string]string{
		// comma seperated values to set
		"set": "mysqlRootPassword=admin@123,persistence.enabled=false,imagePullPolicy=Always",
	}
)

func main() {

	settings = helm_cli.New()
	// Add helm repo
	RepoAdd(repoName, url)

	app := cli.NewApp()
	app.Name = "rkd"
	app.Usage = "Rancher Kubernetes Downloader"

	app.Commands = []cli.Command{
		cmd.ListCommand()
	}

	// app.Run(os.Args)

	// client := github.NewClient(nil)

	// latestRelease, _, err := client.Repositories.GetLatestRelease(context.Background(), "rancher", "rancher")
	// if err != nil {
	// 	fmt.Printf("%s", err)
	// }

	// // Download rancher-image.txt
	// path, err := GetRancherImageList(latestRelease, "tmp")
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// 	return
	// }
	// fmt.Println(path)

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

// RepoAdd adds repo with given name and url
func RepoAdd(name, url string) {
	repoFile := settings.RepositoryConfig
	fmt.Println(repoFile)

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Fatal(err)
	}

	if f.Has(name) {
		fmt.Printf("repository name (%s) already exists\n", name)
		return
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		log.Fatal(err)
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q has been added to your repositories\n", name)
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
			rc, redirect, err := client.Repositories.DownloadReleaseAsset(context.Background(), "rancher", "rancher", asset.GetID())
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
