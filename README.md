# [errs] -- Error handling for Golang

[![Build Status](https://travis-ci.org/spiegel-im-spiegel/errs.svg?branch=master)](https://travis-ci.org/spiegel-im-spiegel/errs)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/spiegel-im-spiegel/errs/master/LICENSE)
[![GitHub release](http://img.shields.io/github/release/spiegel-im-spiegel/errs.svg)](https://github.com/spiegel-im-spiegel/errs/releases/latest)

## Usage

```
err := errs.Wrap(os.ErrInvalid, "wrapped message")
fmt.Println(err)
fmt.Printf("errs.Cause(err): %v\n", errs.Cause(err))
// Output:
// wrapped message: invalid argument
// errs.Cause(err): invalid argument
```

[errs]: https://github.com/spiegel-im-spiegel/errs "spiegel-im-spiegel/errs: Error handling for Golang"
