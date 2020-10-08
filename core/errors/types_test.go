package errors

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func newContextWithUserAgent() context.Context {
	ctx := context.Background()

	md := metadata.Pairs(
		"user-agent",
		"OneSellStaging/1.0.0 (iPhone; iOS 12.3.1; Scale/3.00)",
		"testkey",
		"testval",
	)

	ctx = metadata.NewIncomingContext(ctx, md)

	return ctx
}

func TestNewInvalidContextError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &InvalidContextError{})
	})

	t.Run("when given an context", func(t *testing.T) {
		myCtx := newContextWithUserAgent()
		err := NewInvalidContextError(myCtx)
		msg := err.Error()

		t.Run("returns a formatted error message", func(t *testing.T) {
			assert.Equal(t, msg, "Invalid: *valueCtx is not a valid context")
			t.Log()
		})
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewInvalidContextError(context.Background(), HumanizedMessage("missing metadata"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "missing metadata")
		})
	})
}

func TestInvalidParameterError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &InvalidParameterError{})
	})

	t.Run("when given not an argument name and value", func(t *testing.T) {
		err := NewInvalidParameterError("id", 44455)

		t.Run("returns a formatted error message", func(t *testing.T) {
			msg := err.Error()

			assert.Equal(t, msg, "44455 is not a valid id")
		})
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewInvalidParameterError("id", 44455, HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestMissingParameterError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &MissingParameterError{})
	})

	t.Run("when given not an argument name", func(t *testing.T) {
		err := NewMissingParameterError("id")

		t.Run("returns a formatted error message", func(t *testing.T) {
			msg := err.Error()

			assert.Equal(t, msg, "Missing parameter: id")
		})
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewMissingParameterError("id", HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNotFoundError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &NotFoundError{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		msg := NewNotFoundError().Error()

		assert.Equal(t, msg, "The record could not be found")
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewNotFoundError(HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNewUnauthenticatedError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &UnauthenticatedError{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		msg := NewUnauthenticatedError().Error()

		assert.Equal(t, msg, "The user is not authenticated")
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewUnauthenticatedError(HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNewUnauthorizedError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &UnauthorizedError{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		msg := NewUnauthorizedError().Error()

		assert.Equal(t, msg, "The user is not authorized for this resource")
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewUnauthorizedError(HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNewNetworkError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &NetworkError{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		msg := NewNetworkError().Error()

		assert.Equal(t, msg, "Cannot complete the request")
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewNetworkError(HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNewNotImplementedError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &NotImplementedError{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		msg := NewNotImplementedError().Error()

		assert.Equal(t, msg, "This feature is not implemented")
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewNotImplementedError(HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNewMarshallingError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &MarshallingError{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		msg := NewMarshallingError().Error()

		assert.Equal(t, msg, "Could not parse")
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewMarshallingError(HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNewInvalidEntityError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &InvalidEntityError{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		subject := "This is an error message"

		msg := NewInvalidEntityError(subject).Error()

		assert.Equal(t, msg, subject)
	})

	t.Run("when a humanized option", func(t *testing.T) {
		subject := "This is an error message"
		err := NewInvalidEntityError(subject, HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNewResourceExhaustedError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &ResourceExhausted{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		msg := NewResourceExhaustedError().Error()

		assert.Equal(t, msg, "Resource exhausted")
	})

	t.Run("when a humanized option", func(t *testing.T) {
		err := NewResourceExhaustedError(HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestNewUnknownError(t *testing.T) {
	t.Run("is compliant with Error interface", func(t *testing.T) {
		assert.Implements(t, (*Error)(nil), &UnknownError{})
	})

	t.Run("returns a formatted error message", func(t *testing.T) {
		err := errors.New("this is an error")

		msg := NewUnknownError(err).Error()

		assert.Equal(t, msg, "this is an error")
	})

	t.Run("when a humanized option", func(t *testing.T) {
		newErr := errors.New("this is an error")
		err := NewUnknownError(newErr, HumanizedMessage("this is an error"))
		msg := err.Error()

		t.Run("returns the humanized option", func(t *testing.T) {
			assert.Equal(t, msg, "this is an error")
		})
	})
}

func TestLog(t *testing.T) {
	t.Run("when given a debug log level", func(t *testing.T) {
		t.Run("doesn't panic", func(t *testing.T) {
			assert.NotPanics(t, func() { Log(assert.AnError, DebugLevel) })
		})
	})
	t.Run("when given a info log level", func(t *testing.T) {
		t.Run("doesn't panic", func(t *testing.T) {
			assert.NotPanics(t, func() { Log(assert.AnError, InfoLevel) })
		})
	})
	t.Run("when given a warn log level", func(t *testing.T) {
		t.Run("doesn't panic", func(t *testing.T) {
			assert.NotPanics(t, func() { Log(assert.AnError, WarnLevel) })
		})
	})
	t.Run("when given a warn log level", func(t *testing.T) {
		t.Run("doesn't panic", func(t *testing.T) {
			assert.NotPanics(t, func() { Log(assert.AnError, ErrorLevel) })
		})
	})
	t.Run("when given a error log level", func(t *testing.T) {
		t.Run("it panics", func(t *testing.T) {
			assert.NotPanics(t, func() { Log(assert.AnError, DPanicLevel) })
		})
	})

	t.Run("when given a error log level", func(t *testing.T) {
		t.Run("it panics", func(t *testing.T) {
			assert.Panics(t, func() { Log(assert.AnError, PanicLevel) })
		})
	})
}

func TestMessageCode(t *testing.T) {
	t.Run("returns default values with simple object creation", func(t *testing.T) {
		baseErr := &baseError{}

		assert.EqualValues(t, 0, baseErr.Code())
		assert.EqualValues(t, "", baseErr.Message())
	})

	t.Run("returns default values when using the constructor", func(t *testing.T) {
		baseErr := newBaseError()

		assert.EqualValues(t, 0, baseErr.Code())
		assert.EqualValues(t, "", baseErr.Message())
	})

	t.Run("returns the custom code and message from the options", func(t *testing.T) {
		var errCode uint32 = 42
		errMsg := "life, the universe, and everything"

		baseErr := newBaseError(ErrorCode(errCode, errMsg))

		assert.EqualValues(t, errCode, baseErr.Code())
		assert.EqualValues(t, errMsg, baseErr.Message())
	})
}

func TestCause(t *testing.T) {
	t.Run("returns nil when not given a cause", func(t *testing.T) {
		err := NewMissingParameterError("err")

		assert.Nil(t, err.Cause())
	})

	t.Run("returns an error when given a cause", func(t *testing.T) {
		oldErr := NewMissingParameterError("err1")
		err := NewMissingParameterError("err2", Cause(oldErr))

		assert.Equal(t, oldErr, err.Cause())
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("returns nil when not given a cause", func(t *testing.T) {
		err := NewMissingParameterError("err")

		assert.Equal(t, err, Unwrap(err))
	})

	t.Run("returns an error when given a cause", func(t *testing.T) {
		olderErr := NewMissingParameterError("err0")
		oldErr := NewMissingParameterError("err1", Cause(olderErr))
		err := NewMissingParameterError("err2", Cause(oldErr))

		assert.Equal(t, olderErr, Unwrap(err))
	})
}

func TestBaseError_Component(t *testing.T) {
	t.Run("returns the text component", func(t *testing.T) {
		err := NewMissingParameterError("err")

		assert.Equal(t, "github.com/9count/go-services/core/errors.TestBaseError_Component.func1", err.Component())
	})
}

func TestTrace(t *testing.T) {
	t.Run("returns the text component", func(t *testing.T) {
		var olderErr error
		var oldErr error
		var err error

		func() {
			olderErr = NewMissingParameterError("err0")
		}()

		func() {
			oldErr = NewMissingParameterError("err1", Cause(olderErr))
		}()

		func() {
			err = NewMissingParameterError("err0", Cause(oldErr))
		}()

		stack := Trace(err)

		assert.Equal(t, []string{
			"github.com/9count/go-services/core/errors.TestTrace.func1.3",
			"github.com/9count/go-services/core/errors.TestTrace.func1.2",
			"github.com/9count/go-services/core/errors.TestTrace.func1.1",
		}, stack)
	})
}

func TestAssertEqual(t *testing.T) {
	t.Run("returns true when two errors are the same", func(t *testing.T) {
		foo1 := struct {
			NewError func() error
		}{
			NewError: func() error { return NewMissingParameterError("foo") },
		}

		foo2 := struct {
			NewError func() error
		}{
			NewError: func() error { return NewMissingParameterError("foo") },
		}

		expErr := foo1.NewError()
		actErr := foo2.NewError()

		// this should fail because they have different components
		assert.False(t, assert.ObjectsAreEqualValues(expErr, actErr))

		// this should pass because their .Error() messages only being compared
		AssertEqual(t, expErr, actErr)
	})
}

func TestWrapWithType(t *testing.T) {
	t.Run("returns the same type of error with a new message when called with a new message", func(t *testing.T) {
		err := NewNotFoundError()
		newMsg := "This is a new message"
		newErr := WrapWithType(err, newMsg)

		assert.Equal(t, newMsg, newErr.Error())
		assert.IsType(t, &NotFoundError{}, newErr)
		assert.Equal(t, err, newErr.Cause())
	})
}
