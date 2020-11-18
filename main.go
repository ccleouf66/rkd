package main

import (
	"os"
	"rkd/cmd"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "rkd"
	app.Usage = "Rancher Kubernetes Downloader"

	app.Commands = []cli.Command{
		cmd.ListCommand(),
		cmd.DownloadCommand(),
	}

	app.Run(os.Args)
}
