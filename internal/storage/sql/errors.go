package sql

import "errors"

var (
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrCustomerInvalid       = errors.New("customer invalid")
	ErrCustomerInternal      = errors.New("customer internal error")

	ErrFeedbackNotFound      = errors.New("feedback not found")
	ErrFeedbackInternal      = errors.New("feedback internal error")
	ErrFeedbackInvalid       = errors.New("feedback invalid")
	ErrFeedbackAlreadyExists = errors.New("feedback already exists")
)
