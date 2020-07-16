package errs

import (
	"context"
	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CtxCancelled checks whether a given context has been cancelled
func CtxCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
	}
	return false
}

// CtxError wraps context error to a gRPC error
func CtxError(ctx context.Context, operation string) error {
	if _, ok := ctx.Err().(interface{ Timeout() bool }); ok {
		// Should retry the request
		return status.Errorf(codes.DeadlineExceeded, "couldn't complete %s operation: %v", operation, ctx.Err())
	}
	return status.Errorf(codes.Canceled, "couldn't complete %s operation: %v", operation, ctx.Err())
}

// LogWarn logs a warn message to std logger
func LogWarn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// LogInfo logs infos to std logger
func LogInfo(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// LogError logs an error to std error
func LogError(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// FromJSONMarshal wraps error returned from json.Marshal to a status error
func FromJSONMarshal(err error, obj string) error {
	return status.Errorf(codes.Internal, "failed to json marshal %s: %v", obj, err)
}

// FromJSONUnMarshal wraps error returned from json.Unmarshal to a status error
func FromJSONUnMarshal(err error, obj string) error {
	return status.Errorf(codes.InvalidArgument, "failed to json unmarshal %s: %v", obj, err)
}

// SQLQueryFailed wraps sql error to a status error
func SQLQueryFailed(err error, queryType string) error {
	return status.Errorf(codes.Internal, "failed to execute %s query: %v", queryType, err)
}

// DuplicateField wraps a duplicate field error to a status error
func DuplicateField(field, value string) error {
	return status.Errorf(codes.AlreadyExists, "resource with %s %s already exists", field, value)
}

// SQLQueryNoRows wraps sql no rows found error to a status error
func SQLQueryNoRows(err error) error {
	return status.Errorf(codes.NotFound, "no rows found for query: %v", err)
}

// MissingField creates a status error caused by missing credentials
func MissingField(field string) error {
	return status.Errorf(codes.InvalidArgument, "missing field: %s", field)
}

// CheckingCreds wraps error returned while checking credentials to a status error
func CheckingCreds(err error) error {
	return status.Errorf(codes.Internal, "failed while checking credentials: %v", err)
}

// PermissionDenied result from performing non-priviledged operations
func PermissionDenied(op string) error {
	return status.Errorf(codes.PermissionDenied, "not authorised to perform %s operation", op)
}

// NotFound returns status error that indicates resource does not exist
func NotFound(resource, id string) error {
	return status.Errorf(codes.NotFound, "%s with id: %v not found", resource, id)
}

// AccountNotFound returns status error the account is not exist
func AccountNotFound(username string) error {
	return status.Errorf(codes.NotFound, "no account found for %s", username)
}

// PathogenExists returns a status error indicating that a pathogen exists
func PathogenExists() error {
	return status.Errorf(codes.ResourceExhausted, "pathogen with id or name exists")
}

// PathogenNotSet returns a status error that indicates pathogen is not set in cache
func PathogenNotSet() error {
	return status.Error(codes.Unavailable, "pathogen not set in cache")
}

// RedisCmdFailed wraps an error returned from redis to a status error
func RedisCmdFailed(err error, queryType string) error {
	return status.Errorf(codes.Internal, "failed to execute %s command: %v", queryType, err)
}

// SearchFailed wraps an error returned while searching to a status error
func SearchFailed(err error) error {
	return status.Errorf(codes.Internal, "search failed: %v", err)
}

// ConvertingType wraps error from type assertion to status error
func ConvertingType(err error, from, to string) error {
	return status.Errorf(codes.InvalidArgument, "couldn't convert from %s to %s: %v", from, to, err)
}

// NilObject is error resulting from using nil objects
func NilObject(obj string) error {
	return status.Errorf(codes.InvalidArgument, "nil object not allowed: %s", obj)
}

// FailedToPerformOperation wraps error from caller returning status error to a proper status error
func FailedToPerformOperation(err error, operation string) error {
	return status.Errorf(status.Code(err), "failed to perform operation %s: %v", operation, err)
}

// TokenCredentialNotMatching creates a status error caused by mismatch in token credential
func TokenCredentialNotMatching(cred string) error {
	return status.Errorf(codes.PermissionDenied, "token credential %v do not match", cred)
}

// WrapMessage is a wraps message provided to status error
func WrapMessage(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

// WrapErrWithMessage is a wraps message provided to status error
func WrapErrWithMessage(code codes.Code, err error, msg string) error {
	return status.Errorf(code, "%s: %v", msg, err)
}

// FailedToGenToken wraps error caused while generating jwt token to a status error
func FailedToGenToken(err error) error {
	return status.Errorf(codes.Internal, "failed to generate jwt token: %v", err)
}

// FailedToGenHashedPass wraps error while generating password hash to a grpc error
func FailedToGenHashedPass(err error) error {
	return status.Errorf(codes.Internal, "failed to generate hashed password: %v", err)
}
