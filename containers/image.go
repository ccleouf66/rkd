package containers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker/archive"
	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
)

// DownloadImage download docker images from src and create docker-archive
func DownloadImage(src string, dest string) {

	// Contexts
	defaultPolicy, err := signature.NewPolicyFromFile("policy.json")
	if err != nil {
		fmt.Printf("default policy err: %s\n", err)
	}
	policyContext, err := signature.NewPolicyContext(defaultPolicy)
	if err != nil {
		fmt.Printf("Policy context err: %s\n", err)
	}
	defer policyContext.Destroy()

	// Create systemContext to force linux os
	// sysCtx := &types.SystemContext{
	// 	ArchitectureChoice: "arm",
	// 	OSChoice:           "linux",
	// }

	// Create new dest archive
	aw, err := archive.NewWriter(nil, dest)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	defer aw.Close()

	// Read image lits
	imgList, err := os.Open(src)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer imgList.Close()
	scanner := bufio.NewScanner(imgList)

	for scanner.Scan() {

		// Ref
		imgRef := fmt.Sprintf("%s%s", "docker://", scanner.Text())
		srcRef, err := alltransports.ParseImageName(imgRef)
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		////////// Dest
		// Get image name
		imgNamed, err := reference.ParseDockerRef(scanner.Text())
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		// Get tag from image ref
		strSlice := strings.Split(scanner.Text(), ":")
		tag := "latest"
		if len(strSlice) > 1 {
			tag = strSlice[len(strSlice)-1]
		}
		imgNameTagged, err := reference.WithTag(imgNamed, tag)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		fmt.Printf("NamedTagged ========> %s\n", imgNameTagged)

		// Create dest ref
		destRef, err := aw.NewReference(imgNameTagged)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		//////////

		// Download and create tar
		fmt.Printf("Copy %s to %s\n", imgRef, dest)
		_, err = copy.Image(context.Background(), policyContext, destRef, srcRef, &copy.Options{
			ReportWriter: os.Stdout,
		})
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
