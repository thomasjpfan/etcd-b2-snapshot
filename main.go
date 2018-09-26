package main

import (
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("usage: ./download-etcd-snapshot-b2 [FILEPATH]")
	}
	filePath := os.Args[1]

	// spec, err := ParseENV()
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// dc, err := NewDownloadClientFromSpec(*spec)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// err = dc.Get(spec.EtcdSnapshotBucket, spec.EtcdObjectName, filePath)

	// if err != nil {
	// 	log.Printf("Unable to download snapshot: %v", err)
	// 	return
	// }

	log.Printf("%s downloaded", filePath)

}
