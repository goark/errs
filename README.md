# [errs] -- Error handling for Golang

[![ci status](https://github.com/goark/errs/workflows/ci/badge.svg)](https://github.com/goark/errs/actions)
[![codeql status](https://github.com/goark/errs/workflows/CodeQL/badge.svg)](https://github.com/goark/errs/actions)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/goark/errs/master/LICENSE)
[![GitHub release](http://img.shields.io/github/release/goark/errs.svg)](https://github.com/goark/errs/releases/latest)

Package [errs] implements functions to manipulate error instances.
This package is required Go 1.20 or later.

**Migrated repository to [github.com/goark/errs][errs]**

## Design goals

- Wrap any `error` and collect context at the point of failure
- Add arbitrary key/value context with `WithContext`
- Include caller function name in context by default
- Print structured error data with `%+v` (JSON-like output)
- Handle multi-errors in a concurrency-safe way via `errs.Errors`

## Development

### Requirements

- Go 1.20 or later
- [Task](https://taskfile.dev/) command

### Local validation

```text
task test
task govulncheck
```

Run all maintenance tasks:

```text
task
```

## CI Workflows

- `ci`: lint (`golangci-lint` with `gosec`), tests, and `govulncheck`
- `CodeQL`: scheduled and push/PR static analysis

## Usage

### Sample programs

All sample files under `sample/` use the `run` build tag.

```text
go run -tags run ./sample/sample1.go
```

- New error with cause and context: [sample/sample1.go](sample/sample1.go)
- Wrap existing error with context: [sample/sample2.go](sample/sample2.go)
- Wrap sentinel error with cause: [sample/sample3.go](sample/sample3.go)
- Multiple causes with `errors.Join`: [sample/sample4.go](sample/sample4.go)
- Multi-error with `errs.Join`: [sample/sample5.go](sample/sample5.go)
- Concurrency-safe multi-error accumulation: [sample/sample6.go](sample/sample6.go)
- Zap structured logging object example: [zapobject/example_test.go](zapobject/example_test.go)

### Print formats

- `%v`: human readable error message
- `%#v`: Go-syntax-like internal structure
- `%+v`: structured JSON-like representation

### Helper functions compatible with stdlib

`errs.Is`, `errs.As`, and `errs.Unwrap` are thin wrappers around `errors.Is`, `errors.As`, and `errors.Unwrap`.

`errs.Unwraps` returns `[]error` and works for both single-cause and multi-cause errors.
For a single cause, it returns a one-element slice.
For multiple causes, it returns all causes as a slice.

`errs.EncodeJSON` serializes generic `error` values by traversing unwrap chains when possible.

### Edge-case behavior

- `errs.New("")` returns `nil`
- `errs.Wrap(nil)` returns `nil`
- If `WithCause` is given multiple times, the last cause is used
- `errs.Join(...)` ignores `nil` arguments and returns `nil` if all arguments are `nil`

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

### Structured logging with Zap

Use the submodule `github.com/goark/errs/zapobject` to log errors as structured objects.

```go
logger.Error("failed", zap.Object("error", zapobject.New(err)))
```

Without `zapobject`, `zap.Error(err)` writes string fields only.

## Background article (Japanese)

- [Go 言語用エラーハンドリング・パッケージ](https://text.baldanders.info/release/errs-package-for-golang/)

[errs]: https://github.com/goark/errs "goark/errs: Error handling for Golang"
