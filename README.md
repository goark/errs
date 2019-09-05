# [errs] -- Error handling for Golang

[![Build Status](https://travis-ci.org/spiegel-im-spiegel/errs.svg?branch=master)](https://travis-ci.org/spiegel-im-spiegel/errs)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/spiegel-im-spiegel/errs/master/LICENSE)
[![GitHub release](http://img.shields.io/github/release/spiegel-im-spiegel/errs.svg)](https://github.com/spiegel-im-spiegel/errs/releases/latest)

## Usage

```go
package main

import (
    "fmt"
    "os"

    "github.com/spiegel-im-spiegel/errs"
)

func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return errs.Wrap(
            err,
            "file open error",
            errs.WithParam("path", path),
        )
    }
    defer file.Close()

    return nil
}

func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        fmt.Printf("%v\n", err)             //file open error: open not-exist.txt: no such file or directory
        fmt.Printf("%+v\n", err)            //{"Type":"*errs.Error","Msg":"file open error","Params":{"function":"main.checkFileOpen","path":"not-exist.txt"},"Cause":{"Type":"*os.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}}}
        fmt.Printf("%v\n", errs.Cause(err)) //no such file or directory
    }
}
```

[errs]: https://github.com/spiegel-im-spiegel/errs "spiegel-im-spiegel/errs: Error handling for Golang"
