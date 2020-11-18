package helm

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"rkd/helpers"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

type imageValue struct {
	Registry   string
	Name       string
	Repository string
	Tag        string
}

// RepoAdd adds repo with given name and url
func RepoAdd(name, url string) {

	settings := cli.New()

	repoFile := settings.RepositoryConfig

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
}

// RepoUpdate updates charts for all helm repos
func RepoUpdate() {
	settings := cli.New()

	repoFile := settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		log.Fatal(errors.New("No repositories found. You must add one before updating"))
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, r)
	}

	fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Printf("Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Printf("Chart updated. ⎈ Happy Helming! ⎈\n")
}

// DownloadChart download a chart from public repo to local folder
func DownloadChart(repo string, chart string, version string, dest string) (chartPath string) {

	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(err)
	}

	client := action.NewInstall(actionConfig)

	p := getter.All(settings)

	chartDownloader := &downloader.ChartDownloader{
		Out:              os.Stdout,
		Verify:           downloader.VerifyIfPossible,
		Keyring:          client.ChartPathOptions.Keyring,
		Getters:          p,
		Options:          []getter.Option{},
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
	}

	chartRef := fmt.Sprintf("%s/%s", repo, chart)
	helpers.CreateDestDir(dest)

	path, _, err := chartDownloader.DownloadTo(chartRef, version, dest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Chart downloaded to %s\n", path)

	return path
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}

// GetChartImages take a chart in tgz format and return an image list
func GetChartImages(chartPath string) ([]string, error) {

	chart, err := loader.Load(chartPath)
	if err != nil {
		fmt.Printf("Error when loading chart.\n")
		return nil, err
	}

	cvals, err := chartutil.CoalesceValues(chart, nil)
	if err != nil {
		fmt.Printf("ERR CoalesceValues.\n")
		return nil, err
	}

	imageList, err := GetChartImagesFromValues(cvals)
	if err != nil {
		return nil, err
	}
	return imageList, nil
}

// GetChartImagesFromValues take a map of string corresponding to a values.yaml an return an image list
func GetChartImagesFromValues(values map[string]interface{}) ([]string, error) {
	var imgList []string
	for k, v := range values {
		if k == "image" {
			if _, ok := v.(map[string]interface{}); ok {
				yml, err := yaml.Marshal(v)
				if err != nil {
					return nil, err
				}

				var img imageValue
				err = yaml.Unmarshal(yml, &img)
				if err != nil {
					return nil, err
				}

				imgStr := fmt.Sprintf("%s:%s\n", img.Repository, img.Tag)
				if img.Repository == "" {
					imgStr = fmt.Sprintf("%s:%s\n", img.Name, img.Tag)
				}

				if img.Registry != "" {
					imgStr = fmt.Sprintf("%s/%s", img.Registry, imgStr)
				}
				// remove \n to prevent of reference format error
				imgStr = strings.TrimSuffix(imgStr, "\n")
				imgList = append(imgList, imgStr)
			}
		} else {
			if _, ok := v.(map[string]interface{}); ok {
				il, err := GetChartImagesFromValues(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				imgList = append(imgList, il...)
			}
		}
	}
	return imgList, nil
}
