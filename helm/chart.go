package helm

import (
	"fmt"
	"log"
	"os"
	"rkd/helpers"
	"strings"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}

// DownloadChart download a chart from public repo to local folder
func DownloadChart(repo string, chart string, version string, dest string) (chartPath string, err error) {

	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		return "", err
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
		return "", err
	}

	fmt.Printf("Chart downloaded to %s\n", path)

	return path, nil
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
