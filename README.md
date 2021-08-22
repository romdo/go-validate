<h1 align="center">
  romdo/go-validate
</h1>

<p align="center">
  <strong>
    Yet another Go struct/object validation package, with a focus on simplicity,
    flexibility, and full control over validation logic.
  </strong>
</p>

---

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
