package main

import (
	"fmt"
	"github.com/hairyhenderson/go-fsimpl"
	"github.com/lucchmielowski/oci-fs/ocifs"
	"log"
)

// Sample code to test FS
func main() {
	mux := fsimpl.NewMux()
	mux.Add(ocifs.FS)

	// for example, a URL that points to a subdirectory at a specific tag in a
	// given git repo, hosted on GitHub and authenticated with SSH...
	fsys, err := mux.Lookup("oci://ghcr.io/lucchmielowski/ivpol:multilayered")
	if err != nil {
		log.Fatal(err)
	}

	//file, err := fsys.Open("sample.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()
	//
	//fi, err := file.Stat()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//data := make([]byte, fi.Size())
	//_, err = file.Read(data)
	//
	////Better: data, err := fs.ReadFile(fsys, "sample.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(string(data))

	//files, err := fs.ReadDir(fsys, ".")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//for _, f := range files {
	//	fmt.Println(f.Name())
	//	data, err := fs.ReadFile(fsys, path.Join("dir", f.Name()))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(string(data))
	//}

	file, err := fsys.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fi.Name(), fi.IsDir())
}
