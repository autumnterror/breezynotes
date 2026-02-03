package test

import (
	"errors"
	"testing"

	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrors(t *testing.T) {
	err := format.Error("op", status.Error(codes.FailedPrecondition, "test"))

	s, ok := status.FromError(err)
	if assert.True(t, ok) {
		log.Println(s.Code())
	}
}

func TestErrorsUnw(t *testing.T) {
	err := format.Error("op", status.Error(codes.FailedPrecondition, "test"))
	var newErr error
	if unwErr := errors.Unwrap(err); unwErr != nil {
		newErr = unwErr
	} else {
		newErr = err
	}
	s, ok := status.FromError(newErr)
	if assert.True(t, ok) {
		log.Green(s.Code())
		log.Green(s.Message())
	}
	// err := format.Error("test", errors.New("testerr"))
	// log.Println(err)
	// log.Println(errors.Unwrap(err))

}
