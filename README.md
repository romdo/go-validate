<h1 align="center">
  go-validate
</h1>

<p align="center">
  <strong>
    Yet another Go struct/object validation package, with a focus on simplicity,
    flexibility, and full control over validation logic.
  </strong>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/romdo/go-validate">
    <img src="https://img.shields.io/badge/%E2%80%8B-reference-387b97.svg?logo=go&logoColor=white"
  alt="Go Reference">
  </a>
  <a href="https://github.com/romdo/go-validate/releases">
    <img src="https://img.shields.io/github/v/tag/romdo/go-validate?label=release" alt="GitHub tag (latest SemVer)">
  </a>
  <a href="https://github.com/romdo/go-validate/actions">
    <img src="https://img.shields.io/github/workflow/status/romdo/go-validate/CI.svg?logo=github" alt="Actions Status">
  </a>
  <a href="https://codeclimate.com/github/romdo/go-validate">
    <img src="https://img.shields.io/codeclimate/coverage/romdo/go-validate.svg?logo=code%20climate" alt="Coverage">
  </a>
  <a href="https://github.com/romdo/go-validate/issues">
    <img src="https://img.shields.io/github/issues-raw/romdo/go-validate.svg?style=flat&logo=github&logoColor=white"
alt="GitHub issues">
  </a>
  <a href="https://github.com/romdo/go-validate/pulls">
    <img src="https://img.shields.io/github/issues-pr-raw/romdo/go-validate.svg?style=flat&logo=github&logoColor=white" alt="GitHub pull requests">
  </a>
  <a href="https://github.com/romdo/go-validate/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/romdo/go-validate.svg?style=flat" alt="License Status">
  </a>
</p>

Add validation to any type, by simply implementing the `Validatable` interface:

```go
type Validatable interface {
	Validate() error
}
```

## Import

```go
import "github.com/romdo/go-validate"
```

## Example

```go
type Order struct {
	Books []*Book `json:"books"`
}

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (s *Book) Validate() error {
	var errs error
	if s.Title == "" {
		errs = validate.Append(errs, &validate.Error{
			Field: "Title", Msg: "is required",
		})
	}

	// Helper to perform the same kind of check as above for Title.
	errs = validate.Append(errs, validate.RequireField("Author", s.Author))

	return errs
}

func main() {
	errs := validate.Validate(&Order{Books: []*Book{{Title: ""}}})

	for _, err := range validate.Errors(errs) {
		fmt.Println(err.Error())
	}
}
```

Above example produces the following output:

```
books.0.title: is required
books.0.author: is required
```

## Documentation

Please see the
[Go Reference](https://pkg.go.dev/github.com/romdo/go-validate#section-documentation)
for documentation and examples.

## LICENSE

[MIT](https://github.com/romdo/go-conventionalcommit/blob/main/LICENSE)
