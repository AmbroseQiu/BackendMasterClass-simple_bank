package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FieldViolation(filed string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       filed,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

	statusDetails, _ := statusInvalid.WithDetails(badRequest)

	return statusDetails.Err()
}

func unauthenticationError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthentication err %s", err)
}

func permissionDeniedError(err error) error {
	return status.Errorf(codes.PermissionDenied, "request user info not allowed err %s", err)
}
