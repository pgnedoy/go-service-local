package errors

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"github.com/stretchr/testify/assert"

	"github.com/9count/go-services/core/log"
)

type LogLevel int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel LogLevel = 1
	// InfoLevel is the default logging priority.
	InfoLevel LogLevel = 2
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel LogLevel = 3
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel LogLevel = 4
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel LogLevel = 5
	// PanicLevel logs a message, then panics.
	PanicLevel LogLevel = 6
)

type ErrorOptions struct {
	Code             uint32
	Error            error
	HumanizedMessage string
	LogLevel         LogLevel
	LogMsg           string
	Message          string
}

type ErrorOption func(*ErrorOptions)

func newErrOpts(errOpts ...ErrorOption) *ErrorOptions {
	opts := &ErrorOptions{}

	for _, option := range errOpts {
		option(opts)
	}

	return opts
}

type Error interface {
	Error() string
	Code() uint32
	Message() string
	Cause() error
	Component() string
}

// Base error
type baseError struct {
	options          *ErrorOptions
	HumanizedMessage string
	Error            error
	component        string
}

func (b *baseError) Code() uint32 {
	if b.options != nil {
		return b.options.Code
	}

	return 0
}

func (b *baseError) Message() string {
	if b.options != nil && b.options.Message != "" {
		return b.options.Message
	}

	return b.HumanizedMessage
}

// Unwraps the error one level
func (b *baseError) Cause() error {
	if b.Error != nil {
		return b.Error
	}

	return nil
}

// Returns the component where the error originated
func (b *baseError) Component() string {
	return b.component
}

func ErrorCode(code uint32, message string) ErrorOption {
	return func(args *ErrorOptions) {
		args.Code = code
		args.Message = message
	}
}

func HumanizedMessage(msg string) ErrorOption {
	return func(args *ErrorOptions) {
		args.HumanizedMessage = msg
	}
}

func Cause(err error) ErrorOption {
	return func(args *ErrorOptions) {
		args.Error = err
	}
}

// Accepts a log level and then logs the error under that level with more useful information.
func Log(err error, level LogLevel) {
	msg := err.Error()

	switch level {
	case DebugLevel:
		log.Debug(context.Background(), msg, log.WithError(err))
	case InfoLevel:
		log.Info(context.Background(), msg, log.WithError(err))
	case WarnLevel:
		log.Warn(context.Background(), msg, log.WithError(err))
	case ErrorLevel:
		log.Error(context.Background(), msg, log.WithError(err))
	case DPanicLevel:
		log.DPanic(context.Background(), msg, log.WithError(err))
	case PanicLevel:
		log.Panic(context.Background(), msg, log.WithError(err))
	}
}

// Completed unwraps an error all the way down to its root error
func Unwrap(err error) error {
	e, ok := err.(Error)
	if !ok {
		return nil
	}

	cause := e.Cause()

	if cause != nil {
		return Unwrap(cause)
	}

	return err
}

// Completed unwraps an error all the way down to its root error
func WrapWithType(err error, message string) Error {
	cause := Cause(err)

	switch err.(type) {
	case *InvalidEntityError:
		return &InvalidEntityError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	case *InvalidParameterError:
		return &InvalidParameterError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	case *MissingParameterError:
		return &MissingParameterError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	case *InvalidContextError:
		return &InvalidContextError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	case *NotImplementedError:
		return &NotImplementedError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	case *ResourceExhausted:
		return &ResourceExhausted{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	case *UnauthenticatedError:
		return &UnauthenticatedError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	case *UnauthorizedError:
		return &UnauthorizedError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	case *NotFoundError:
		return &NotFoundError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	default:
		return &UnknownError{
			baseError:        newBaseError(cause),
			FormattedMessage: message,
		}
	}
}

// Creates a trace from each level of errors component
func Trace(err error) []string {
	stack := make([]string, 0)

	e, ok := err.(Error)

	for ok {
		stack = append(stack, e.Component())

		e, ok = e.Cause().(Error)
	}

	return stack
}

func newBaseError(opts ...ErrorOption) baseError {
	pc := make([]uintptr, 15)
	// Warning: This is sort of inferred that each constructor calls this base error function
	// So because this skips 2 calls deep, if some higher wrapping happens, this frame will be wrong.
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	component := frame.Function

	errOpts := newErrOpts(opts...)

	return baseError{
		options:          errOpts,
		HumanizedMessage: errOpts.HumanizedMessage,
		Error:            errOpts.Error,
		component:        component,
	}
}

type InvalidContextError struct {
	baseError
	FormattedMessage string
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func (e *InvalidContextError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewInvalidContextError(ctx context.Context, opts ...ErrorOption) *InvalidContextError {
	return &InvalidContextError{
		baseError:        newBaseError(opts...),
		FormattedMessage: fmt.Sprintf("Invalid: %v is not a valid context", getType(ctx)),
	}
}

// Invalid Parameter - User exceptions that should not be logged
type InvalidParameterError struct {
	baseError
	FormattedMessage string
}

func (e *InvalidParameterError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewInvalidParameterError(param string, value interface{}, opts ...ErrorOption) *InvalidParameterError {
	return &InvalidParameterError{
		baseError:        newBaseError(opts...),
		FormattedMessage: fmt.Sprintf("%v is not a valid %s", value, param),
	}
}

// Missing parameter error indicates that a required parameter was not passed to the caller,
// or that an empty value was passed where it should have been supplied
type MissingParameterError struct {
	baseError
	FormattedMessage string
}

func (e *MissingParameterError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewMissingParameterError(param string, opts ...ErrorOption) *MissingParameterError {
	return &MissingParameterError{
		baseError:        newBaseError(opts...),
		FormattedMessage: fmt.Sprintf("Missing parameter: %s", param),
	}
}

// Not Found error indicates that a resource was requested under some specified parameters, but that
// resource does not exist where it was expected to
type NotFoundError struct {
	baseError
	FormattedMessage string
}

func (e *NotFoundError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewNotFoundError(opts ...ErrorOption) *NotFoundError {
	return &NotFoundError{
		baseError:        newBaseError(opts...),
		FormattedMessage: "The record could not be found",
	}
}

// Unauthenticated error indicates that valid identity credentials could not be determined from
// the requester
type UnauthenticatedError struct {
	baseError
	FormattedMessage string
}

func (e *UnauthenticatedError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewUnauthenticatedError(opts ...ErrorOption) *UnauthenticatedError {
	return &UnauthenticatedError{
		baseError:        newBaseError(opts...),
		FormattedMessage: "The user is not authenticated",
	}
}

// Unauthorized error indicates that the requester's identity has been authenticated but is
// not allowed to perform the operation
type UnauthorizedError struct {
	baseError
	FormattedMessage string
}

func (e *UnauthorizedError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewUnauthorizedError(opts ...ErrorOption) *UnauthorizedError {
	return &UnauthorizedError{
		baseError:        newBaseError(opts...),
		FormattedMessage: "The user is not authorized for this resource",
	}
}

// Network error indicates that there was an issue at the network layer between external resources
// that prevents the transfer of that data
type NetworkError struct {
	baseError
	FormattedMessage string
}

func (e *NetworkError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewNetworkError(opts ...ErrorOption) *NetworkError {
	return &NetworkError{
		baseError:        newBaseError(opts...),
		FormattedMessage: "Cannot complete the request",
	}
}

// Not implemented error indicates that the requested action is not currently implemented by the program
type NotImplementedError struct {
	baseError
	FormattedMessage string
}

func (e *NotImplementedError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewNotImplementedError(opts ...ErrorOption) *NotImplementedError {
	return &NotImplementedError{
		baseError:        newBaseError(opts...),
		FormattedMessage: "This feature is not implemented",
	}
}

// Marshalling error indicates that the
type MarshallingError struct {
	baseError
	FormattedMessage string
}

func (e *MarshallingError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewMarshallingError(opts ...ErrorOption) *MarshallingError {
	return &MarshallingError{
		baseError:        newBaseError(opts...),
		FormattedMessage: "Could not parse",
	}
}

// Invalid Entity - User exceptions that should not be logged
type InvalidEntityError struct {
	baseError
	FormattedMessage string
}

func (e *InvalidEntityError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewInvalidEntityError(msg string, opts ...ErrorOption) *InvalidEntityError {
	return &InvalidEntityError{
		baseError:        newBaseError(opts...),
		FormattedMessage: msg,
	}
}

// Resource Exhausted
type ResourceExhausted struct {
	baseError
	FormattedMessage string
}

func (e *ResourceExhausted) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewResourceExhaustedError(opts ...ErrorOption) *ResourceExhausted {
	return &ResourceExhausted{
		baseError:        newBaseError(opts...),
		FormattedMessage: "Resource exhausted",
	}
}

// Unknown Error
type UnknownError struct {
	baseError
	FormattedMessage string
}

func (e *UnknownError) Error() string {
	if e.HumanizedMessage != "" {
		return e.HumanizedMessage
	}

	return e.FormattedMessage
}

func NewUnknownError(err error, opts ...ErrorOption) *UnknownError {
	return &UnknownError{
		baseError:        newBaseError(opts...),
		FormattedMessage: err.Error(),
	}
}

// AssertEqual Compares two errors
func AssertEqual(t assert.TestingT, expected, actual error) bool {
	return assert.Equal(t, expected.Error(), actual.Error())
}
