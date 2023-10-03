package service

import (
	"github.com/Bek0sh/online-market-auth/internal/myerrors"
	"github.com/Bek0sh/online-market-auth/pkg/proto"
	"github.com/Bek0sh/online-market-auth/pkg/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateRegisterUser(req *proto.RegisterUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetName()); err != nil {
		violations = append(violations, myerrors.FieldViolation("name", err))
	}
	if err := val.ValidatePhoneNumber(req.GetPhoneNumber()); err != nil {
		violations = append(violations, myerrors.FieldViolation("phone_number", err))
	}
	if err := val.ValidateUsername(req.GetSurname()); err != nil {
		violations = append(violations, myerrors.FieldViolation("surname", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, myerrors.FieldViolation("name", err))
	}

	return violations
}

func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badrequest := &errdetails.BadRequest{
		FieldViolations: violations,
	}
	statusInvaild := status.New(codes.InvalidArgument, "invalid parametres")

	statusDetails, err := statusInvaild.WithDetails(badrequest)
	if err != nil {
		return statusInvaild.Err()
	}

	return statusDetails.Err()
}
