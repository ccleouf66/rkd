package cmd

import (
	"github.com/urfave/cli"
)

// UploadCommand return downlaod cli command
func UploadCommand() cli.Command {
	UploadFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "dest",
			Usage: "set the destination registry",
		},
		cli.BoolFlag{
			Name:  "signature",
			Usage: "boolean to enable the download of the Cosign signature, optional, false by default",
		},
	}

	return cli.Command{
		Name:   "upload",
		Usage:  "upload the datapack",
		Action: UploadDataPack,
		Flags:  UploadFlags,
	}
}

// DownloadDataPack downlaod Rancher chart and images
func UploadDataPack(c *cli.Context) error {
	// V1 - Upload container image
	// test presence of tarball
	// Fetch the manifest json from the tarball
	// for each repoTags in the tarball, upload images with dest
	// Check if the signature dir exist -> set bool
	// for each signature, read the manifest and upload the signature with dest

	// V2 - Upload helm chart

	return nil
}
