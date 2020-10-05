package containers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker/archive"
	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

// DownloadImage download docker images from src and create docker-archive
func DownloadImage(imgList []string, dest string) {

	// Check regs auth

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

	// Create new dest archive
	aw, err := archive.NewWriter(nil, dest)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	defer aw.Close()

	//for scanner.Scan() {
	for _, img := range imgList {

		// Ref
		imgRef := fmt.Sprintf("%s%s", "docker://", img)
		srcRef, err := alltransports.ParseImageName(imgRef)
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		////////// Dest
		// Get image name
		imgNamed, err := reference.ParseDockerRef(img)
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		// Get tag from image ref
		strSlice := strings.Split(img, ":")
		tag := "latest"
		if len(strSlice) > 1 {
			tag = strSlice[len(strSlice)-1]
		}
		imgNameTagged, err := reference.WithTag(imgNamed, tag)
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		// Create dest ref
		destRef, err := aw.NewReference(imgNameTagged)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		//////////

		// Create systemContext to select os and arch
		sysCtx := &types.SystemContext{
			ArchitectureChoice: "amd64",
			OSChoice:           "linux",
		}

		// Download and create tar
		fmt.Printf("Copy %s to %s\n", imgRef, dest)
		_, err = copy.Image(context.Background(), policyContext, destRef, srcRef, &copy.Options{
			ReportWriter: os.Stdout,
			SourceCtx:    sysCtx,
		})
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
