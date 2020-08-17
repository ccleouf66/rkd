package git

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"rkd/helpers"

	"github.com/google/go-github/v32/github"
)

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
			helpers.CreateDestDir(dest)
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
