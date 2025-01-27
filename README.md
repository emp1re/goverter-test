<p align="center">
    <img width="300" src=".github/logo.svg" />
</p>

<h1 align="center">goverter</h1>
<p align="center"><i>a "type-safe Go converter" generator</i></p>
<p align="center">
    <a href="https://github.com/emp1re/goverter-test/actions/workflows/build.yml">
        <img alt="Build Status" src="https://github.com/emp1re/goverter-test/actions/workflows/build.yml/badge.svg">
    </a>
     <a href="https://codecov.io/gh/jmattheis/goverter">
        <img alt="codecov" src="https://codecov.io/gh/jmattheis/goverter/branch/main/graph/badge.svg">
    </a>
    <a href="https://goreportcard.com/report/github.com/jmattheis/goverter">
        <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/jmattheis/goverter">
    </a>
    <a href="https://pkg.go.dev/github.com/jmattheis/goverter">
        <img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/jmattheis/goverter.svg">
    </a>
    <a href="https://github.com/emp1re/goverter-test/releases/latest">
        <img alt="latest release" src="https://img.shields.io/github/release/jmattheis/goverter.svg">
    </a>
</p>

goverter is a tool for creating type-safe converters. All you have to
do is create an interface and execute goverter. The project is meant as
alternative to [jinzhu/copier](https://github.com/jinzhu/copier) that doesn't
use reflection.

[Installation](https://goverter.jmattheis.de/guide/install) ᛫
[CLI](https://goverter.jmattheis.de/reference/cli) ᛫
[Config](https://goverter.jmattheis.de/reference/settings)

## Features

- **Fast execution**: No reflection is used at runtime
- Automatically converts builtin types: slices, maps, named types, primitive
  types, pointers, structs with same fields
- [Deep copies](https://en.wikipedia.org/wiki/Object_copying#Deep_copy) per
  default and supports [shallow
  copying](https://en.wikipedia.org/wiki/Object_copying#Shallow_copy)
- **Customizable**: [You can implement custom converter methods](https://goverter.jmattheis.de/reference/extend)
- [Clear errors when generating the conversion methods](https://goverter.jmattheis.de/guide/error-early) if
  - the target struct has unmapped fields
  - types cannot be converted without losing information
- Detailed [documentation](https://goverter.jmattheis.de/) with a lot examples
- Thoroughly tested, see [our test scenarios](./scenario)

## Getting Started

1. Ensure your `go version` is 1.16 or above

1. Create a go modules project if you haven't done so already

    ```bash
    $ go mod init module-name
    ```

1. Create your converter interface and mark it with a comment containing
   [`goverter:converter`](https://goverter.jmattheis.de/reference/converter)

    `input.go`

    ```go
    package example

    // goverter:converter
    type Converter interface {
      ConvertItems(source []Input) []Output

      // goverter:ignore Irrelevant
      // goverter:map Nested.AgeInYears Age
      Convert(source Input) Output
    }

    type Input struct {
      Name string
      Nested InputNested
    }
    type InputNested struct {
      AgeInYears int
    }
    type Output struct {
      Name string
      Age int
      Irrelevant bool
    }
    ```

    See [Settings](https://goverter.jmattheis.de/reference/settings) for more information.

1. Run `goverter`:

    ```bash
    $ go run github.com/emp1re/goverter-test/cmd/goverter@latest gen ./
    ```

    It's recommended to use an explicit version instead of `latest`. See
    [Installation](https://goverter.jmattheis.de/guide/install) and
    [CLI](https://goverter.jmattheis.de/reference/cli) for more information.

1. goverter created a file at `./generated/generated.go`, it may look like this:

    ```go
    package generated

    import example "goverter/example"

    type ConverterImpl struct{}

    func (c *ConverterImpl) Convert(source example.Input) example.Output {
        var exampleOutput example.Output
        exampleOutput.Name = source.Name
        exampleOutput.Age = source.Nested.AgeInYears
        return exampleOutput
    }
    func (c *ConverterImpl) ConvertItems(source []example.Input) []example.Output {
        var exampleOutputList []example.Output
        if source != nil {
            exampleOutputList = make([]example.Output, len(source))
            for i := 0; i < len(source); i++ {
                exampleOutputList[i] = c.Convert(source[i])
            }
        }
        return exampleOutputList
    }
    ```

    See [Generation](https://goverter.jmattheis.de/explanation/generation) for more information.

## Versioning

goverter uses [SemVer](http://semver.org/) for versioning the cli.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE)
file for details

_Logo by [MariaLetta](https://github.com/MariaLetta)_
