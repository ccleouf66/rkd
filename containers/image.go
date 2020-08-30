package containers

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/containers/image/docker/archive"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
)

// DownloadImage download docker images from src and create docker-archive
func DownloadImage(src string, dest string) {

	// Contexts
	defaultPolicy, err := signature.DefaultPolicy(nil)
	if err != nil {
		fmt.Printf("default policy err: %s\n", err)
	}
	policyContext, err := signature.NewPolicyContext(defaultPolicy)
	if err != nil {
		fmt.Printf("Policy context err: %s\n", err)
	}
	defer policyContext.Destroy()

	// // Ref and Dest
	// imgRef := fmt.Sprintf("%s%s", "docker://", src)
	// srcRef, err := alltransports.ParseImageName(imgRef)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// }
	// archDest := fmt.Sprintf("%s:%s:%s", "docker-archive", dest, src)
	// destRef, err := alltransports.ParseImageName(archDest)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// }

	// data, err := copy.Image(context.Background(), policyContext, destRef, srcRef, &copy.Options{
	// 	ReportWriter: os.Stdout,
	// })
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// }
	// fmt.Printf("%s\n", data)

	// Read image lits
	imgList, err := os.Open(src)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer imgList.Close()

	scanner := bufio.NewScanner(imgList)

	wip := archive.NewReference
	for scanner.Scan() {

		// Ref and Dest
		imgRef := fmt.Sprintf("%s%s", "docker://", src)
		srcRef, err := alltransports.ParseImageName(imgRef)
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		destRef := wip.NewReference(scanner.Text())
		copy.Image(context.Background(), policyContext, srcRef, destRef, &copy.Options{
			ReportWriter: os.Stdout,
		})
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
	wip.Finish()
}
