package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs     []string
}

func (e *MultiError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}

	template := "%d errors occured:\n\t* %s\n"

	errs := strings.Join(e.errs, "\t* ")
	return fmt.Sprintf(template, len(e.errs), errs)
}

func Append(err error, errs ...error) *MultiError {
	multiError, ok := err.(*MultiError)
	if !ok {
		multiError = new(MultiError)
		if err != nil {
			multiError.errs = []string{err.Error()}
		}
	}

	for _, err := range errs {
		multiError.errs = append(multiError.errs, err.Error())
	}

	return multiError
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
