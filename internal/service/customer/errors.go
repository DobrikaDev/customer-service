package customer

import "errors"

var ErrCustomerNotFound = errors.New("customer not found")
var ErrCustomerAlreadyExists = errors.New("customer already exists")
var ErrCustomerInternal = errors.New("customer internal error")
var ErrCustomerInvalid = errors.New("customer invalid")

var ErrFeedbackNotFound = errors.New("feedback not found")
var ErrFeedbackInternal = errors.New("feedback internal error")
var ErrFeedbackInvalid = errors.New("feedback invalid")
var ErrFeedbackAlreadyExists = errors.New("feedback already exists")