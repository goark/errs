# [errs] -- Error handling for Golang

[![check vulns](https://github.com/goark/errs/workflows/vulns/badge.svg)](https://github.com/goark/errs/actions)
[![lint status](https://github.com/goark/errs/workflows/lint/badge.svg)](https://github.com/goark/errs/actions)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/goark/errs/master/LICENSE)
[![GitHub release](http://img.shields.io/github/release/goark/errs.svg)](https://github.com/goark/errs/releases/latest)

Package [errs] implements functions to manipulate error instances.
This package is required Go 1.20 or later.

**Migrated repository to [github.com/goark/errs][errs]**

## Usage

### Create new error instance with cause

```go
package main

import (
    "fmt"
    "os"

    "github.com/goark/errs"
)

func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return errs.New(
            "file open error",
            errs.WithCause(err),
            errs.WithContext("path", path),
        )
    }
    defer file.Close()

    return nil
}

func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        fmt.Printf("%v\n", err)  // file open error: open not-exist.txt: no such file or directory
        fmt.Printf("%#v\n", err) // *errs.Error{Err:&errors.errorString{s:"file open error"}, Cause:&fs.PathError{Op:"open", Path:"not-exist.txt", Err:0x2}, Context:map[string]interface {}{"function":"main.checkFileOpen", "path":"not-exist.txt"}}
        fmt.Printf("%+v\n", err) // {"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"file open error"},"Context":{"function":"main.checkFileOpen","path":"not-exist.txt"},"Cause":{"Type":"*fs.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}}}
    }
}
```

### Wrapping error instance

```go
package main

import (
    "fmt"
    "os"

    "github.com/goark/errs"
)

func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return errs.Wrap(
            err,
            errs.WithContext("path", path),
        )
    }
    defer file.Close()

    return nil
}

func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        fmt.Printf("%v\n", err)  // open not-exist.txt: no such file or directory
        fmt.Printf("%#v\n", err) // *errs.Error{Err:&fs.PathError{Op:"open", Path:"not-exist.txt", Err:0x2}, Cause:<nil>, Context:map[string]interface {}{"function":"main.checkFileOpen", "path":"not-exist.txt"}}
        fmt.Printf("%+v\n", err) // {"Type":"*errs.Error","Err":{"Type":"*fs.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}},"Context":{"function":"main.checkFileOpen","path":"not-exist.txt"}}
    }
}
```

### Wrapping error instance with cause

```go
package main

import (
    "errors"
    "fmt"
    "os"

    "github.com/goark/errs"
)

func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return errs.Wrap(
            errors.New("file open error"),
            errs.WithCause(err),
            errs.WithContext("path", path),
        )
    }
    defer file.Close()

    return nil
}

func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        fmt.Printf("%v\n", err)  // file open error: open not-exist.txt: no such file or directory
        fmt.Printf("%#v\n", err) // *errs.Error{Err:&errors.errorString{s:"file open error"}, Cause:&fs.PathError{Op:"open", Path:"not-exist.txt", Err:0x2}, Context:map[string]interface {}{"function":"main.checkFileOpen", "path":"not-exist.txt"}}
        fmt.Printf("%+v\n", err) // {"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"file open error"},"Context":{"function":"main.checkFileOpen","path":"not-exist.txt"},"Cause":{"Type":"*fs.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}}}
    }
}
```

### Create new error instance with multiple causes

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"

    "github.com/goark/errs"
)

func generateMultiError() error {
    return errs.New("error with multiple causes", errs.WithCause(errors.Join(os.ErrInvalid, io.EOF)))
}

func main() {
    err := generateMultiError()
    fmt.Printf("%+v\n", err)            // {"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"error with multiple causes"},"Context":{"function":"main.generateMultiError"},"Cause":{"Type":"*errors.joinError","Msg":"invalid argument\nEOF","Cause":[{"Type":"*errors.errorString","Msg":"invalid argument"},{"Type":"*errors.errorString","Msg":"EOF"}]}}
    fmt.Println(errors.Is(err, io.EOF)) // true
}
```

### Handling multiple errors

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"

    "github.com/goark/errs"
)

func generateMultiError() error {
    return errs.Join(os.ErrInvalid, io.EOF)
}

func main() {
    err := generateMultiError()
    fmt.Printf("%+v\n", err)            // {"Type":"*errs.Errors","Errs":[{"Type":"*errors.errorString","Msg":"invalid argument"},{"Type":"*errors.errorString","Msg":"EOF"}]}
    fmt.Println(errors.Is(err, io.EOF)) // true
}
```

```go
package main

import (
    "fmt"
    "sync"

    "github.com/goark/errs"
)

func generateMultiError() error {
    errlist := &errs.Errors{}
    var wg sync.WaitGroup
    for i := 1; i <= 2; i++ {
        i := i
        wg.Add(1)
        go func() {
            defer wg.Done()
            errlist.Add(fmt.Errorf("error %d", i))
        }()
    }
    wg.Wait()
    return errlist.ErrorOrNil()
}

func main() {
    err := generateMultiError()
    fmt.Printf("%+v\n", err) // {"Type":"*errs.Errors","Errs":[{"Type":"*errors.errorString","Msg":"error 2"},{"Type":"*errors.errorString","Msg":"error 1"}]}
}
```

[errs]: https://github.com/goark/errs "goark/errs: Error handling for Golang"
