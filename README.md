# Golang Hivex bindings

Golang bindings for [hivex](https://github.com/libguestfs/hivex).

## Basic usage

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    hivex "github.com/gabriel-samfira/go-hivex"
)

func main() {
    gopath := os.Getenv("GOPATH")
    // There are a few test hives inside the package
    // Feel free to use your own
    hivePath := filepath.Join(
        gopath,
        "src/github.com/gabriel-samfira/go-hivex",
        "testdata",
        "rlenvalue_test_hive")

    // If you plan to write to the hive, replace
    hive, err := hivex.NewHivex(hivePath, hivex.READONLY)
    if err != nil {
        log.Fatal(err)
    }

    root, err := hive.Root()
    if err != nil {
        log.Fatal(err)
    }
    // Get a child node called ModerateValueParent
    child, err := hive.NodeGetChild(root, "ModerateValueParent")
    if err != nil {
        log.Fatal(err)
    }
    // child will hold an int64 representing the offset of the node
    fmt.Println(child)

    // fetch the name of the Node
    childName, err := hive.NodeName(child)
    if err != nil {
        log.Fatal(err)
    }
    // print out the name (should be "ModerateValueParent")
    fmt.Println(childName)
}
```