package xerrors

import (
	"errors"
	"testing"
)

func TestAppend(t *testing.T) {
	var errs Errors
	errs = append(errs, errors.New("a"))
	errs = append(errs, errors.New("b"))
	if len(errs) != 2 {
		t.Fatalf("invalid errors count - got: %d, want: 2", len(errs))
	}
}

func TestError(t *testing.T) {
	var errs Errors
	if errs.Error() != "" {
		t.Fatalf("invalid error message - got: %q, want: %q", errs.Error(), "")
	}

	errs = append(errs, errors.New("a"))
	errs = append(errs, errors.New("b"))
	if errs.Error() != "a\nb" {
		t.Fatalf("invalid error message - got: %q, want: %q", errs.Error(), "a\nb")
	}
}
