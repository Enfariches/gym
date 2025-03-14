package utils

import "errors"

const (
	Ok                      string = "ok"
	InternalError                  = "internal_error"
	InvalidBody                    = "invalid_body"
	InvalidRequest                 = "invalid_request"
	MissingToken                   = "missing_token"
	InvalidToken                   = "invalid_token"
	WrongEmail                     = "not_allowed_email"
	LoginErrorStatus               = "login_error"
	RegistrationErrorStatus        = "registration_error"
	AlreadyExists                  = "already_exists"
	DoesntExist                    = "not_exist"
	WrongPassword                  = "wrong_password"
	FieldRequired                  = "field_required"
	UploadError                    = "upload_error"
	PermissionsErrorStatus         = "permission_denied"
	EmailConfirmNotFound           = "email_confirm_not_found"
	WrongConfirmationCode          = "wrong_code"
)

var (
	EmailRequired                  = errors.New("email address is required")
	EmailExists                    = errors.New("email address already in use")
	PasswordRequired               = errors.New("password of at least 6 characters is required")
	ConnectionError                = errors.New("connection error")
	EmailConfirmationRequiredError = errors.New("email is not confirmed")
)
