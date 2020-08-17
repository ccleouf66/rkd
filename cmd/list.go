package cmd

import (
	"fmt"

	"rkd/git"

	"github.com/urfave/cli"
)

// ListCommand .
func ListCommand() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"c"},
		Usage:   "List Rancher stable release",
		Action:  listRepoStableRelease,
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
