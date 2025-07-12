package main

import (
	"log"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/lucchmielowski/oci-fs/ocifs"
)

// Sample code to test FS
func main() {
	mux := fsimpl.NewMux()
	mux.Add(ocifs.FS)

	// for example, a URL that points to a subdirectory at a specific tag in a
	// given git repo, hosted on GitHub and authenticated with SSH...
	fsys, err := mux.Lookup("oci://ghcr.io/nirmata/demo-image-compliance-policies:block-high-vulnerabilites")
	if err != nil {
		log.Fatal(err)
	}

	file, err := fsys.Open("sample.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, fi.Size())
	file.Read(data)

	//Better: data, err := fs.ReadFile(fsys, "sample.yaml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
