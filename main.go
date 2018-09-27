package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("usage: ./download-etcd-snapshot-b2 [FILEPATH]")
	}
	filePath := os.Args[1]

	spec, err := ParseENV()
	if err != nil {
		log.Fatalf("%v", err)
	}

	client := http.Client{
		Timeout: time.Second * 10,
	}

	dl := B2Downloader{
		ApplicationID:         spec.B2ApplicationID,
		ApplicationKey:        spec.B2ApplicationKey,
		BucketName:            spec.B2BucketName,
		Object:                spec.B2Object,
		ObjectID:              spec.B2ObjectID,
		DownloadRetryInterval: spec.B2DownloadRetryInterval,
		HTTPClient:            &client,
	}

	err = dl.Download(filePath)

	if err != nil && err.Error() == "File does not exist on B2" {
		log.Printf("File does not exist on B2")
		return
	}

	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%s downloaded", filePath)

}
