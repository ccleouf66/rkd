package gcs

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func InitCli(ctx context.Context, bktName string) (*storage.BucketHandle, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}

	bkt := client.Bucket(bktName)

	return bkt, nil
}

func UploadObject(ctx context.Context, bucket *storage.BucketHandle, path string, fName string) (err error) {

	// build file name
	if fName == "" {
		fName = fmt.Sprintf("%s-%s", path, time.Now().Format("2006-01-02"))
		strSlice := strings.Split(fName, "/")
		if len(strSlice) > 1 {
			fName = strSlice[len(strSlice)-1]
		}
	}

	// src file
	src, err := os.Open(path)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	r := bufio.NewReader(src)

	// upload file
	fmt.Printf("Uploading %s to %s\n", path, fName)
	w := bucket.Object(fName).NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		fmt.Printf("%s", err)
		return err
	}

	// close connexion
	if err := w.Close(); err != nil {
		fmt.Printf("%s", err)
		return err
	}
	return nil
}

func ListObject(ctx context.Context, bucket *storage.BucketHandle) (l []string, err error) {
	query := &storage.Query{Prefix: ""}

	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		l = append(l, attrs.Name)
	}
	return l, nil
}
