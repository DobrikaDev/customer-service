package delivery

import (
	"DobrikaDev/customer-service/internal/domain"
	customerpb "DobrikaDev/customer-service/internal/generated/proto/customer"
	"DobrikaDev/customer-service/internal/service/customer"
	"context"

	"github.com/dr3dnought/gospadi"
	"go.uber.org/zap"
)

func (s *Server) CreateCustomer(ctx context.Context, req *customerpb.CreateCustomerRequest) (*customerpb.CreateCustomerResponse, error) {
	if req.Customer == nil {
		return &customerpb.CreateCustomerResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "customer is required",
			},
		}, nil
	}
	if req.Customer.MaxId == "" {
		return &customerpb.CreateCustomerResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "max id is required",
			},
		}, nil
	}
	if req.Customer.Name == "" {
		return &customerpb.CreateCustomerResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "name is required",
			},
		}, nil
	}

	customer := &domain.Customer{
		MaxID: req.Customer.MaxId,
		Name:  req.Customer.Name,
		About: req.Customer.About,
		Type:  converCustomerTypeToDomain(req.Customer.Type),
	}
	customer, err := s.customerService.CreateCustomer(ctx, customer)
	if err != nil {
		return &customerpb.CreateCustomerResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Info("customer created", zap.Any("customer", customer))

	return &customerpb.CreateCustomerResponse{
		Customer: convertCustomerToProto(customer),
	}, nil
}

func (s *Server) GetCustomers(ctx context.Context, req *customerpb.GetCustomersRequest) (*customerpb.GetCustomersResponse, error) {
	if req.MaxId == "" {
		return &customerpb.GetCustomersResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "max id is required",
			},
		}, nil
	}
	customers, count, err := s.customerService.GetCustomers(ctx, req.MaxId, int(req.Limit), int(req.Offset))
	if err != nil {
		return &customerpb.GetCustomersResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	s.logger.Info("customers fetched", zap.Any("customers", customers), zap.Int("count", count))
	return &customerpb.GetCustomersResponse{
		Customers: gospadi.Map(customers, convertCustomerToProto),
		Total:     int32(count),
	}, nil
}

func (s *Server) GetCustomerByMaxID(ctx context.Context, req *customerpb.GetCustomerByMaxIDRequest) (*customerpb.GetCustomerByMaxIDResponse, error) {
	if req.MaxId == "" {
		return &customerpb.GetCustomerByMaxIDResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "max id is required",
			},
		}, nil
	}
	customer, err := s.customerService.GetCustomerByMaxID(ctx, req.MaxId)
	if err != nil {
		return &customerpb.GetCustomerByMaxIDResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	s.logger.Info("customer fetched", zap.Any("customer", customer))
	return &customerpb.GetCustomerByMaxIDResponse{
		Customer: convertCustomerToProto(customer),
	}, nil
}

func (s *Server) UpdateCustomer(ctx context.Context, req *customerpb.UpdateCustomerRequest) (*customerpb.UpdateCustomerResponse, error) {
	if req.Customer == nil {
		return &customerpb.UpdateCustomerResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "customer is required",
			},
		}, nil
	}
	if req.Customer.MaxId == "" {
		return &customerpb.UpdateCustomerResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "max id is required",
			},
		}, nil
	}

	customer := &domain.Customer{
		MaxID: req.Customer.MaxId,
		Name:  req.Customer.Name,
		About: req.Customer.About,
		Type:  converCustomerTypeToDomain(req.Customer.Type),
	}
	customer, err := s.customerService.UpdateCustomer(ctx, customer)
	if err != nil {
		return &customerpb.UpdateCustomerResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	s.logger.Info("customer updated", zap.Any("customer", customer))

	return &customerpb.UpdateCustomerResponse{
		Customer: convertCustomerToProto(customer),
	}, nil
}

func (s *Server) DeleteCustomer(ctx context.Context, req *customerpb.DeleteCustomerRequest) (*customerpb.DeleteCustomerResponse, error) {
	if req.MaxId == "" {
		return &customerpb.DeleteCustomerResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "max id is required",
			},
		}, nil
	}
	err := s.customerService.DeleteCustomer(ctx, req.MaxId)
	if err != nil {
		return &customerpb.DeleteCustomerResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	s.logger.Info("customer deleted", zap.String("max_id", req.MaxId))
	return &customerpb.DeleteCustomerResponse{
		MaxId: req.MaxId,
	}, nil
}

func convertCustomerToProto(customer *domain.Customer) *customerpb.Customer {
	return &customerpb.Customer{
		MaxId:     customer.MaxID,
		Name:      customer.Name,
		About:     customer.About,
		Type:      convertCustomerTypeToProto(customer.Type),
		CreatedAt: int32(customer.CreatedAt.Unix()),
		UpdatedAt: int32(customer.UpdatedAt.Unix()),
	}
}
func converCustomerTypeToDomain(customerType customerpb.CustomerType) domain.CustomerType {
	switch customerType {
	case customerpb.CustomerType_CUSTOMER_TYPE_INDIVIDUAL:
		return domain.CustomerTypeIndividual
	case customerpb.CustomerType_CUSTOMER_TYPE_BUSINESS:
		return domain.CustomerTypeCompany
	}
	return domain.CustomerTypeIndividual
}
func convertCustomerTypeToProto(customerType domain.CustomerType) customerpb.CustomerType {
	switch customerType {
	case domain.CustomerTypeIndividual:
		return customerpb.CustomerType_CUSTOMER_TYPE_INDIVIDUAL
	case domain.CustomerTypeCompany:
		return customerpb.CustomerType_CUSTOMER_TYPE_BUSINESS
	}
	return customerpb.CustomerType_CUSTOMER_TYPE_UNSPECIFIED
}

func convertErrorToProto(err error) *customerpb.Error {
	switch err {
	case customer.ErrCustomerNotFound:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_NOT_FOUND,
			Message: err.Error(),
		}
	case customer.ErrCustomerAlreadyExists:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_ALREADY_EXISTS,
			Message: err.Error(),
		}
	case customer.ErrCustomerInvalid:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
			Message: err.Error(),
		}
	case customer.ErrCustomerInternal:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_INTERNAL,
			Message: err.Error(),
		}
	case customer.ErrFeedbackNotFound:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_NOT_FOUND,
			Message: err.Error(),
		}
	case customer.ErrFeedbackAlreadyExists:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_ALREADY_EXISTS,
			Message: err.Error(),
		}
	case customer.ErrFeedbackInvalid:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
			Message: err.Error(),
		}
	case customer.ErrFeedbackInternal:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_INTERNAL,
			Message: err.Error(),
		}
	default:
		return &customerpb.Error{
			Code:    customerpb.ErrorCode_ERROR_CODE_UNSPECIFIED,
			Message: err.Error(),
		}
	}
}
