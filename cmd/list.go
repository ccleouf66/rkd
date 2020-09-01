package cmd

import (
	"fmt"

	"rkd/git"

	"github.com/urfave/cli"
)

// ListCommand .
func ListCommand() cli.Command {
	listFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "helm",
			Usage: "List",
			Value: "latest",
		},
	}

	return cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List Rancher stable release",
		Action:  listRepoStableRelease,
		Flags:   listFlags,
	}
}

// ListRepoStableRelease .
func listRepoStableRelease(c *cli.Context) {
	releases, err := git.GetRepoStablRelease("rancher", "rancher")
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	fmt.Printf("Num. Name - TagName\n")
	for index, release := range releases {
		fmt.Printf("%d. %s - %s\n", index, release.GetName(), release.GetTagName())
	}
}
