package service

import (
	"net/http"

	"github.com/Bek0sh/online-market-auth/internal/models"
	"github.com/Bek0sh/online-market-auth/internal/myerrors"
	"github.com/Bek0sh/online-market-auth/pkg/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func validateRegisterUser(req *models.RegisterUser) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.Username); err != nil {
		violations = append(violations, myerrors.FieldViolation("name", err))
	}
	if err := val.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		violations = append(violations, myerrors.FieldViolation("phone_number", err))
	}
	if err := val.ValidateUsername(req.Surname); err != nil {
		violations = append(violations, myerrors.FieldViolation("surname", err))
	}
	if err := val.ValidatePassword(req.Password); err != nil {
		violations = append(violations, myerrors.FieldViolation("password", err))
	}
	if err := val.ValidatePassword(req.ConfirmPassword); err != nil {
		violations = append(violations, myerrors.FieldViolation("confirm_password", err))
	}

	return violations
}

func validateSignInUser(req *models.SignInUser) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		violations = append(violations, myerrors.FieldViolation("phone_number", err))
	}
	if err := val.ValidatePassword(req.Password); err != nil {
		violations = append(violations, myerrors.FieldViolation("password", err))
	}

	return violations
}

func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badrequest := &errdetails.BadRequest{
		FieldViolations: violations,
	}
	statusInvaild := status.New(http.StatusBadRequest, "invalid parametres")

	statusDetails, err := statusInvaild.WithDetails(badrequest)
	if err != nil {
		return statusInvaild.Err()
	}

	return statusDetails.Err()
}
