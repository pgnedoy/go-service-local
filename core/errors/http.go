package errors

import (
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/status"
)

// Transforms an application error into a http error by mapping the error type to a response code and writing to the
func ToHttp(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *InvalidParameterError, *MissingParameterError, *InvalidEntityError:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	case *NotImplementedError:
		http.Error(w, err.Error(), http.StatusNotImplemented)
	case *InvalidContextError:
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
	case *ResourceExhausted:
		http.Error(w, err.Error(), http.StatusTooManyRequests)
	case *UnauthenticatedError:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case *UnauthorizedError:
		http.Error(w, err.Error(), http.StatusForbidden)
	case *NotFoundError:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "Oops, something went wrong", http.StatusBadRequest)
	}
}

// Transforms an application error into an http error with a body by mapping a grpc error to an http error
func GrpcToHTTP(w http.ResponseWriter, err error) error {
	st, ok := status.FromError(err)

	if !ok {
		return NewMarshallingError()
	}

	jsonResp, err := json.Marshal(st.Message())

	if err != nil {
		return err
	}

	w.WriteHeader(runtime.HTTPStatusFromCode(st.Code()))

	// TODO: transform this into a header
	_, writeErr := w.Write(jsonResp)

	return writeErr
}
