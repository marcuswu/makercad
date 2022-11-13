sysfont
=======
[![Build Status](https://github.com/adrg/sysfont/workflows/CI/badge.svg)](https://github.com/adrg/sysfont/actions?query=workflow%3ACI)
[![pkg.go.dev documentation](https://pkg.go.dev/badge/github.com/adrg/sysfont)](https://pkg.go.dev/github.com/adrg/sysfont)
[![MIT license](https://img.shields.io/badge/license-MIT-red.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Go report card](https://goreportcard.com/badge/github.com/adrg/sysfont)](https://goreportcard.com/report/github.com/adrg/sysfont)
[![GitHub issues](https://img.shields.io/github/issues/adrg/sysfont)](https://github.com/adrg/sysfont/issues)
[![Buy me a coffee](https://img.shields.io/static/v1.svg?label=%20&message=Buy%20me%20a%20coffee&color=579fbf&logo=buy%20me%20a%20coffee&logoColor=white)](https://ko-fi.com/T6T72WATK)

sysfont is a small package that makes it easy to identify installed fonts. It
is useful for listing installed fonts or for matching fonts based on user
queries. The matching process also suggests viable font alternatives.

The package uses a collection of standard fonts compiled from the
[os-font-list](https://github.com/adrg/os-font-list) project along with string
processing and similarity metrics for scoring font matches, in order to account
for partial or inexact input queries.

Full documentation can be found at: https://pkg.go.dev/github.com/adrg/sysfont.

## Installation

```
go get github.com/adrg/sysfont
```

## Usage

#### List fonts

```go
finder := sysfont.NewFinder(nil)

for _, font := range finder.List() {
    fmt.Println(font.Family, font.Name, font.Filename)
}
```

#### Match fonts

The matching process has three steps. Identification of the best matching
installed font, based on the specified query, is attempted first. If no close
match is found, alternative fonts are searched. If no alternative font is
found, a suitable default font is returned.

```go
finder := sysfont.NewFinder(nil)

terms := []string{
    "AmericanTypewriter",
    "AmericanTypewriter-Bold",
    "Andale",
    "Arial",
    "Arial Bold",
    "Arial-BoldItalicMT",
    "ArialMT",
    "Baskerville",
    "Candara",
    "Corbel",
    "Gill Sans",
    "Hoefler Text Bold",
    "Impact",
    "Palatino",
    "Symbol",
    "Tahoma",
    "Times",
    "Times Bold",
    "Times BoldItalic",
    "Times Italic Bold",
    "Times Roman",
    "Verdana",
    "Verdana-Italic",
    "Webdings",
    "ZapfDingbats",
}

for _, term := range terms {
    font := finder.Match(term)
    fmt.Printf("%-30s -> %-30s (%s)\n", term, font.Name, font.Filename)
}
```

Output:
![sysfont test output minimal](https://raw.githubusercontent.com/adrg/adrg.github.io/master/assets/projects/sysfont/output_minimal.png)

A more comprehensive test made on Ubuntu:
![sysfont test output full](https://raw.githubusercontent.com/adrg/adrg.github.io/master/assets/projects/sysfont/output-full.png)

## References

For more information see:
- [os-font-list](https://github.com/adrg/os-font-list)
- [strutil](https://github.com/adrg/strutil)
- [xdg](https://github.com/adrg/xdg)

## Contributing

Contributions in the form of pull requests, issues or just general feedback,
are always welcome.  
See [CONTRIBUTING.MD](https://github.com/adrg/sysfont/blob/master/CONTRIBUTING.md).

## License
Copyright (c) 2019 Adrian-George Bostan.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](https://github.com/adrg/sysfont/blob/master/LICENSE) for more details.
