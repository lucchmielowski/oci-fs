# OCI-FS

This repo POCs usage of OCI manifests as an FS module for `go-fsimpl`

## Usage

```go
package main

import
    "github.com/hairyhenderson/go-fsimpl"
    "github.com/lucchmielowski/oci-fs/ocifs"
)

// Sample code to test FS
func main() {
    mux := fsimpl.NewMux()
    mux.Add(ocifs.FS)
	
    // Get manifest from OCI registry
    fsys, err := mux.Lookup("oci://ghcr.io/nirmata/demo-image-compliance-policies:block-high-vulnerabilites")
    if err != nil {
       log.Fatal(err)
    }

    // Read `sample.yaml` file from the extracted layers
    data, err := fs.ReadFile(fsys, "sample.yaml")
    if err != nil {
       log.Fatal(err)
    }
	
    // Should print a kvyerno policy
    log.Println(string(data))
}
```

## Features

- [x] Implement fs.FS
- [ ] Implement fs.ReadDirFS
- [ ] Add perfomance tests