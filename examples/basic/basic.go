package main

import (
	"fmt"

	"github.com/romdo/go-validate"
)

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
